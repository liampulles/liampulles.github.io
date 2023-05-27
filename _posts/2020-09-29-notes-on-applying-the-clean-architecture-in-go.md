---
layout: post
title: Notes on applying "The Clean Architecture" in Go
shareable: true
---
<section>
    <p>Having had an appetite for experimenting with a REST API in Go, my research indicated that Robert
    "Uncle Bob" Martin's <i>Clean Architecture</i> is something like the "architecture of choice" among Gophers.
    However, I have found that there is a lack of good resources out there for applying it in Go. So here is my
    attempt to fill that gap...</p>
    <p><i><b>Note:</b> I am not an Architect, and I've only applied this pattern on one application - so please take everything
    I'm about to say with a big pinch of salt. You can find the example REST API I implemented
    <b><a href="https://github.com/liampulles/matchstick-video">here</a></b></i></p>
</section>

<section>
    <h2>So, what is this "Clean Architecture" thing anyway?</h2>
    <aside>
        <figure>
            <img src="/images/clean-architecture-diagram.png" height="731">
            <figcaption><i>The Layers of The Clean Architecture.</i></figcaption>
        </figure>
    </aside>
    <p>The best resource on this is definitely Uncle Bob's
    <a href="https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html">blog post</a>
    on the subject, but in essence the architecture is concerned with abstracting implementation details
    away from your business logic. It does this by separating code into a series of layers:
    <ul>
        <li><b>The Entity Layer</b>: This is where your business logic sits.</li>
        <li><b>The Use Case Layer</b>: This is where your application logic sits.</li>
        <li><b>The Adapter Layer</b>: This is where the glue between the use case and driver layers sits.</li>
        <li><b>The Frameworks and Drivers Layer</b>: This is where code interacting with external libraries sits.</li>
    </ul>
    </p>
</section>

<section>
    <h2>How do I know where to put my Code?</h2>
    <aside class="optional">
        <figure>
            <img src="/images/rules-of-the-game-vhs-front.jpg" max-width="180" height="2216">
            <figcaption><i>Enterprise business rules belong in the Entity layer.</i></figcaption>
        </figure>
        <figure>
            <img src="/images/postgresql-icon.png" width="50%" height="256">
            <img src="/images/gorilla-icon.jpeg" width="50%" height="225">
            <figcaption><i>Any notion of a specific DB implementation or router framework certainly belongs in the Drivers and Frameworks layer.</i></figcaption>
        </figure>
        <figure>
            <img src="/images/mongodb-icon.png" width="50%" height="512">
            <img src="/images/grpc-icon.png" width="50%" height="1500">
            <figcaption><i>A helpful way to think through your layers is to imagine radically changing the DB and web frameworks.</i></figcaption>
        </figure>
    </aside>
    <p>This is a question I wrestled with quite a bit in this project. Firstly, remember that you are learning,
    so its okay to put down the code in a layer with a TODO and then come back later to refactor it.
    However, I have come up with a set of questions you can ask yourself in order to come to a considered
    decision:
    <ol>
        <li><b>Can this code be copied and pasted into another application,</b> and be useful without modification? If so,
        then it probably belongs in the entity layer. Examples include: string validation functions (e.g. <i>isBlank</i>) and
        your core business logic.</li>
        <li><b>Does this code deal with orchestrating the logic of a transaction,</b> e.g. finding all users in the database?
        If so, it probably belongs in the use case layer. Examples include: the core orchestration logic (as mentioned) as
        well as factories for constructing entity types, and interfaces which code in the adapter layer must implement.</li>
        <li><b>Does the code call any external code or use any external types?</b> If so, it probably belongs in the
        frameworks and drivers layer. Examples include: PostgreSQL driver configurations, Mux handlers and server setup,
        etc.</li>
        <li><b>If your code does not fit with one of the above questions,</b> and/or deals with bridging calls to and from
        the use case and driver and frameworks layers, then it probably belongs in the adapter layer. Examples include:
        Your controllers, JSON transformers, SQL code, and configuration utilities.</li>
    </ol>
    </p>
    <p>You may be thinking that a lot of the things I've put into the adapter layer may belong in the use case layer.
    What I would recommend is: imagine what layers would need to change if you switched your API and DB implementations to a
    completely different philosophy.</p>
    <p>For example: what if I used GRPC instead of a REST API, and a NoSQL DB instead of an SQL
    DB? The use case logic should not be (hugely) affected by this switch - thus we (e.g.) put the Repository interface in
    the use case layer, but we put the SQL (and potentially NoSQL) implementations of that repository in the adapter layer.</p>
</section>

<section>
    <h2>How do I structure my packages?</h2>
    <aside>
        <figure>
            <img src="/images/matchstick-video-package-structure.png" height="414">
            <figcaption><i>Package structure for</i> <a href="https://github.com/liampulles/matchstick-video">Matchstick Video!</a></figcaption>
        </figure>
    </aside>
    <p>Firstly, we're going to rename some of our layers for practical reasons. I've renamed the entity layer to the domain
    layer, because entities have their own meaning in a DDD sense which I wish to maintain. Secondly, I'm going to simplify
    "drivers and frameworks" to just "drivers".</p>
    <p>Remember that the
    names are just a suggestion - as are the number of layers. You can (and should) adjust them to what makes sense for your
    team. The important thing is that you separate logic so as to try and minimize the time it takes to figure out where you
    need to make changes, and to minimize the number of lines which need to be touched for a change. We are trying to enable
    ongoing change to the greatest degree possible.</p>
    <p>Anyway, on the side you'll see my package structure. The important ones here are the <i>adapter</i>, <i>domain</i>,
    <i>driver</i>, <i>usecase</i>, and <i>wire</i> packages. Which brings us to...</p>
</section>

<section>
    <h2>The "Wire" Package/Layer</h2>
    <p>The <i>wire</i> layer sits outside of the other layers, and it deals with dependency injection. Basically, this
    layer contains a function which first creates the service instances which have no dependencies, then uses those to construct
    the services which rely on those services, and so-on until you've created the server, which you can then execute.</p>
{% highlight go %}
// CreateServerFactory injects all the dependencies needed to create
// http.ServerFactory
func CreateServerFactory(source goConfig.Source) (http.ServerFactory, error) {
    // Each "tap" below indicates a level of dependency
    configStore, err := config.NewStoreImpl(
        source,
    )
    if err != nil {
        return nil, err
    }
    errorParser := adapterDb.NewErrorParserImpl()

    // --- NEXT TAP ---
    helperService := sql.NewHelperServiceImpl(errorParser)
    databaseService, err := db.NewDatabaseServiceImpl(
        configStore,
    )
    if err != nil {
        return nil, err
    }
    inventoryItemConstructor := entity.NewInventoryItemConstructorImpl()
    muxWrapper := mux.NewWrapperImpl()
    //...
{% endhighlight %}
    <p>This is fairly mundane, boiler-platery code - but it is pretty easy to understand and update, and I haven't found good
    enough cause to use an external package to do it (though if you want to use an external package,
    <a href="https://github.com/google/wire"> Google's wire tool</a> seems to be a good choice).</p>
    <p>What IS important is that you make a separate package (and layer) for wiring, as it is going to be importing code from all
    over your project, and you want to make sure that there aren't circular dependencies.</p>
</section>

<section>
    <h2>JSON struct tags</h2>
    <p>One might be compelled to put JSON struct tags (as well as ORM stuct tags, etc.) on your entities in the domain package,
    but this of course would be a violation of our segregation rules: application communication does not form part of the
    business rules. If we go back to our thought experiment to reinforce this point: what if we wanted to use GRPC instead? This should not require us
    touching the domain package, so clearly we cannot put any JSON tags on the entity to begin with.</p>
    <p>This does not mean that we cannot customize how our objects are serialized - it just means that we need to make use of an "intermediary" struct
    in order to do this. For example:</p>
{% highlight go %}
// FromInventoryItemView converts a view to JSON
func (e *EncoderServiceImpl) FromInventoryItemView(view *inventory.ViewVO) ([]byte, error) {
    intermediary := mapViewIntermediary(view)

    bytes, err := json.Marshal(intermediary)
    if err != nil {
        return nil, fmt.Errorf("could not convert inventory item view to json - marshal error: %w", err)
    }
    return bytes, nil
}

func mapViewIntermediary(view *inventory.ViewVO) *jsonViewVO {
    return &jsonViewVO{
        ID:        view.ID,
        Name:      view.Name,
        Location:  view.Location,
        Available: view.Available,
    }
}

type jsonViewVO struct {
    ID        entity.ID `json:"id"`
    Name      string    `json:"name"`
    Location  string    `json:"location"`
    Available bool      `json:"available"`
}
{% endhighlight %}
    <p>Here we first map our use case view (which contains the elements of the entity we want to expose) to an
    intermediary struct, and then marshal the struct. By doing this, we've decoupled the entity from concerns over how
    it is viewed externally, and we've decoupled that view from its encoding.</p>
</section>

<section>
    <h2>Entities</h2>
    <aside>
        <figure>
            <img src="/images/my-super-sweet-16.jpg" height="375">
            <figcaption><i>The Entity is King, and treats itself as such.</i></figcaption>
        </figure>
    </aside>
    <p>I like to think of Entities as bratty Beverly Hill teenagers: they have an entitled view of the world which may not map to reality.</p>
    This abstracted view includes everything ranging from:
    <ul>
        <li>Method parameters</li>
        <li>Responses</li>
        <li>Errors</li>
        <li>Services that Entities might need to use</li>
        <li>etc.</li>
    </ul>
    <p>Here is a good example from the Inventory Item Entity:</p>
{% highlight go %}
// IsAvailable will return true if the inventory item may
// be checked out - false otherwise.
func (i *InventoryItemImpl) IsAvailable() bool {
    return i.available
}

// Checkout will mark the inventory item as unavilable.
// If the inventory item is not available,
// then an error is returned.
func (i *InventoryItemImpl) Checkout() error {
    if !i.available {
        return fmt.Errorf("cannot check out inventory item - it is unavailable")
    }
    i.available = false
    return nil
}

// CheckIn will mark the inventory item as available.
// If the inventory item is available, then an
// error is returned.
func (i *InventoryItemImpl) CheckIn() error {
    if i.available {
        return fmt.Errorf("cannot check in inventory item - it is already checked in")
    }
    i.available = true
    return nil
}
{% endhighlight %}
    <p>Here we have not exposed the <i>available</i> field on the struct. Instead, we encapsulate access via methods,
    some of which may throw errors. This is done to protect the entity - it is the use case layer's job to deal with
    these errors.</p>
</section>

<section>
    <h2>Conclusion</h2>
    <p>I found The Clean Architecture to work very well for a REST API, and for Go. It takes
    a fair bit of time to set up, but what you are left with is a very modular and easy to change structure. I will
    definitely use it for my web apps GOing forward. ;)</p>
</section>