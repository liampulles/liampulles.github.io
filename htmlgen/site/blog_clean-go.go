package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"clean-go",
		`Notes on applying "The Clean Architecture" in Go`,
		"This article introduces Clean Architecture in Go. It discusses the layers involved, provides guidance on code placement, and suggests package structuring.",
		civil.Date{Year: 2020, Month: time.September, Day: 29},
		markdown(cleanGo_opening),
		"clean-architecture-diagram.png",
		cleanGo_whatIs,
		cleanGo_howDoIKnow,
		cleanGo_howDoIStructure,
		cleanGo_theWirePackage,
		cleanGo_jsonStruct,
		cleanGo_entities,
		cleanGo_conclusion,
	))
}

const cleanGo_opening = `
Having had an appetite for experimenting with a REST API in Go, my research
indicated that Robert "Uncle Bob" Martin's Clean Architecture is something like
the "architecture of choice" among Gophers. However, I have found that there is
a lack of good resources out there for applying it in Go. So here is my attempt
to fill that gap...

***Note:** I am not an Architect, and I've only applied this pattern on one
application - so please take everything I'm about to say with a big pinch of
salt. You can find the example REST API I implemented
[here](https://github.com/liampulles/matchstick-video)*`

var cleanGo_whatIs = section(
	`So, what is this "Clean Architecture" thing anyway?`,
	markdown(`
The best resource on this is definitely Uncle Bob's
[blog post](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
on the subject, but in essence the architecture is concerned with abstracting
implementation details away from your business logic. It does this by separating
code into a series of layers:

* **The Entity Layer**: This is where your business logic sits.
* **The Use Case Layer**: This is where your application logic sits.
* **The Adapter Layer**: This is where the glue between the use case and driver
  layers sits.
* **The Frameworks and Drivers Layer**: This is where code interacting with
  external libraries sits.`),
	withAsideFigure(figure("The Layers of The Clean Architecture.", image(
		"clean-architecture-diagram.png",
		731, 731,
		"Diagram of the clean architecture",
	))),
)

var cleanGo_howDoIKnow = section(
	`How do I know where to put my Code?`,
	markdown(`
This is a question I wrestled with quite a bit in this project. Firstly,
remember that you are learning, so its okay to put down the code in a layer with
a TODO and then come back later to refactor it. However, I have come up with a
set of questions you can ask yourself in order to come to a considered decision:

1. **Can this code be copied and pasted into another application,** and be
   useful without modification? If so, then it probably belongs in the entity
   layer. Examples include: string validation functions (e.g. isBlank) and your
   core business logic.
1. **Does this code deal with orchestrating the logic of a transaction,**
   e.g. finding all users in the database? If so, it probably belongs in the use
   case layer. Examples include: the core orchestration logic (as mentioned) as
   well as factories for constructing entity types, and interfaces which code in
   the adapter layer must implement.
1. **Does the code call any external code or use any external types?** If so, it
   probably belongs in the frameworks and drivers layer. Examples include:
   PostgreSQL driver configurations, Mux handlers and server setup, etc.
1. **If your code does not fit with one of the above questions,** and/or deals
   with bridging calls to and from the use case and driver and frameworks
   layers, then it probably belongs in the adapter layer. Examples include: Your
   controllers, JSON transformers, SQL code, and configuration utilities.

You may be thinking that a lot of the things I've put into the adapter layer may
belong in the use case layer. What I would recommend is: imagine what layers
would need to change if you switched your API and DB implementations to a
completely different philosophy.

For example: what if I used GRPC instead of a REST API, and a NoSQL DB instead
of an SQL DB? The use case logic should not be (hugely) affected by this
switch - thus we (e.g.) put the Repository interface in the use case layer, but
we put the SQL (and potentially NoSQL) implementations of that repository in the
adapter layer.
`),
	withAsideFigure(optionalFigure("Any notion of a specific DB implementation or router framework certainly belongs in the Drivers and Frameworks layer.",
		image("postgresql-icon.png", 110, 110, "PostgreSQL logo"),
		image("gorilla-icon.jpeg", 110, 110, "Gorilla mux logo"),
	)),
	withAsideFigure(optionalFigure("A helpful way to think through your layers is to imagine radically changing the DB and web frameworks.",
		image("mongodb-icon.png", 110, 110, "MongoDB logo"),
		image("grpc-icon.png", 110, 110, "gRPC logo"),
	)),
)

var cleanGo_howDoIStructure = section(
	`How do I structure my packages?`,
	markdown(`
Firstly, we're going to rename some of our layers for practical reasons. I've
renamed the entity layer to the domain layer, because entities have their own
meaning in a DDD sense which I wish to maintain. Secondly, I'm going to simplify
"drivers and frameworks" to just "drivers".

Remember that the names are just a suggestion - as are the number of layers. You
can (and should) adjust them to what makes sense for your team. The important
thing is that you separate logic so as to try and minimize the time it takes to
figure out where you need to make changes, and to minimize the number of lines
which need to be touched for a change. We are trying to enable ongoing change
to the greatest degree possible.
	
Anyway, on the side you'll see my package structure. The important ones here are
the *adapter*, *domain*, *driver*, *usecase*, and *wire* packages. Which brings
us to...`),
	withAsideFigure(figure(`Package structure for <a href="https://github.com/liampulles/matchstick-video">Matchstick Video!</a>`,
		image(
			"matchstick-video-package-structure.png",
			232, 414,
			"Screenshot of a folder tree view for the app",
		),
	)),
)

var cleanGo_theWirePackage = section(
	`How do I structure my packages?`,
	markdown(`
The *wire* layer sits outside of the other layers, and it deals with dependency
injection. Basically, this layer contains a function which first creates the
service instances which have no dependencies, then uses those to construct the
services which rely on those services, and so-on until you've created the
server, which you can then execute.

~~~go
// CreateServerFactory injects all the dependencies
// needed to create http.ServerFactory
func CreateServerFactory(
    source goConfig.Source,
) (http.ServerFactory, error) {
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
~~~

This is fairly mundane, boiler-platery code - but it is pretty easy to
understand and update, and I haven't found good enough cause to use an external
package to do it (though if you want to use an external package,
[Google's wire tool](https://github.com/google/wire) seems to be a good choice).

What IS important is that you make a separate package (and layer) for wiring, as
it is going to be importing code from all over your project, and you want to
make sure that there aren't circular dependencies.`),
)

var cleanGo_jsonStruct = section(
	`JSON struct tags`,
	markdown(`
One might be compelled to put JSON struct tags (as well as ORM stuct tags, etc.)
on your entities in the domain package, but this of course would be a violation
of our segregation rules: application communication does not form part of the
business rules. If we go back to our thought experiment to reinforce this point:
what if we wanted to use GRPC instead? This should not require us touching the
domain package, so clearly we cannot put any JSON tags on the entity to begin
with.

This does not mean that we cannot customize how our objects are serialized - it
just means that we need to make use of an "intermediary" struct in order to do
this. For example:

~~~go
// FromInventoryItemView converts a view to JSON
func (e *EncoderServiceImpl) FromInventoryItemView(
    view *inventory.ViewVO,
) ([]byte, error) {
    intermediary := mapViewIntermediary(view)

    bytes, err := json.Marshal(intermediary)
    if err != nil {
        return nil, fmt.Errorf("could not convert inventory"+
        " item view to json - marshal error: %w", err)
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
    ID        entity.ID \'json:"id"\'
    Name      string    \'json:"name"\'
    Location  string    \'json:"location"\'
    Available bool      \'json:"available"\'
}
~~~

Here we first map our use case view (which contains the elements of the entity
we want to expose) to an intermediary struct, and then marshal the struct. By
doing this, we've decoupled the entity from concerns over how it is viewed
externally, and we've decoupled that view from its encoding.`),
)

var cleanGo_entities = section(
	`Entities`,
	markdown(`
I like to think of Entities as bratty Beverly Hill teenagers: they have an
entitled view of the world which may not map to reality.

This abstracted view includes everything ranging from:
* Method parameters
* Responses
* Errors
* Services that Entities might need to use
* etc.

Here is a good example from the Inventory Item Entity:

~~~go
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
        return fmt.Errorf("cannot check out inventory"+
        " item - it is unavailable")
    }
    i.available = false
    return nil
}

// CheckIn will mark the inventory item as available.
// If the inventory item is available, then an
// error is returned.
func (i *InventoryItemImpl) CheckIn() error {
    if i.available {
        return fmt.Errorf("cannot check in inventory"+
        " item - it is already checked in")
    }
    i.available = true
    return nil
}
~~~

Here we have not exposed the available field on the struct. Instead, we
encapsulate access via methods, some of which may throw errors. This is done to
protect the entity - it is the use case layer's job to deal with these errors.
`),
	withAsideFigure(figure("The Entity is King, and treats itself as such.", image(
		"my-super-sweet-16.jpg",
		375, 500,
		"A boy dressed as a king on a throne in a super sweet 16 reality show",
	))),
)

var cleanGo_conclusion = section(
	`Conclusion`,
	markdown(`
I found The Clean Architecture to work very well for a REST API, and for Go. It
takes a fair bit of time to set up, but what you are left with is a very modular
and easy to change structure. I will definitely use it for my web apps GOing
forward. ;)`),
)
