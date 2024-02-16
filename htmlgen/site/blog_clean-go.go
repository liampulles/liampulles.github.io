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
		cleanGo_whatIs,
		cleanGo_howDoIKnow,
		cleanGo_howDoIStructure,
		cleanGo_theWirePackage,
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
	figure("The Layers of The Clean Architecture.", image(
		"clean-architecture-diagram.png",
		731, 731,
		"Diagram of the clean architecture",
	)),
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
	figure("Any notion of a specific DB implementation or router framework certainly belongs in the Drivers and Frameworks layer.",
		image("postgresql-icon.png", 110, 110, "PostgreSQL logo"),
		image("gorilla-icon.jpeg", 110, 110, "Gorilla mux logo"),
	),
	figure("A helpful way to think through your layers is to imagine radically changing the DB and web frameworks.",
		image("mongodb-icon.png", 110, 110, "MongoDB logo"),
		image("grpc-icon.png", 110, 110, "gRPC logo"),
	),
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
	figure(`Package structure for <a href="https://github.com/liampulles/matchstick-video">Matchstick Video!</a>`,
		image(
			"matchstick-video-package-structure.png",
			232, 414,
			"Screenshot of a folder tree view for the app",
		),
	),
)

var cleanGo_theWirePackage = section(
	`How do I structure my packages?`,
	markdown(`
The *wire* layer sits outside of the other layers, and it deals with dependency
injection. Basically, this layer contains a function which first creates the
service instances which have no dependencies, then uses those to construct the
services which rely on those services, and so-on until you've created the
server, which you can then execute.

`+codeFigureMarkdown("go",
		`// CreateServerFactory injects all the dependencies
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
    //...`)+`

`),
)
