package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"proverbs",
		`Coding Proverbs`,
		"Contains proverbs and advice around programming and developer life.",
		civil.Date{Year: 2022, Month: time.July, Day: 9},
		markdown(proverb_opening),
		superSection("Code Design", "",
			proverb_applyingDry,
			proverb_designByContract,
		),
		superSection("Golang", "",
			proverb_implChecks,
			proverb_dontNeedDep,
		),
		superSection("Makefiles", "",
			proverb_toolCheck,
		),
	))
}

const proverb_opening = `
Over the years I discover small pieces of cogent advice about programming and
being a developer that I find useful.

This page tries to collect those thoughts - for my benefit and perhaps others.
Email me if you have something good, I'll consider adding it. ;)`

var proverb_applyingDry = section(
	"Applying DRY: intrinsically vs. coincidentally similar code",
	markdown(`
In Uni I learned that you should NEVER have duplicated code. But in the
workplace, I've learned from wise old devs and from firsthand experience that
this is not always the case.

My philosophy is as follows: **Refactor intrinsically similar code - not
coincidentally similar.**

Coincidentally similar code is likely to change and differentiate. If you
preemptively pull both sets of logic out into a common function, you'll later
have to create parameters for differences, and you'll dirty your previously
clean code.

An example of this might be a plain text template function and an HTML template
function. While the initial implementation of these two functions is the same,
the HTML function will need to be able to escape certain characters in the
future - if you've pulled out a common function, you'll need now need to add a
boolean argument to switch between HTML or plain text (which is ugly).

Intrinsically similar code NEEDS to be put in a common function, however. If you
don't do this, you risk a business rule being inconsistently implemented.

An example of this is formatting dates as strings. If you identify two separate
pieces (or even one piece) of code which use an anonymous string like
"DD/MM/YYYY" to format a date, that should be pulled out into a global constant
or date utility function. If you don't do this, you risk someone changing one of
the format strings (or adding another anonymous one somewhere) to "MM/DD/YYYY" -
and you'll end up confusing/misinforming your users and potentially making the
client liable.

Determining whether similarities are coincidental or intrinsic is a matter of
interrogating your domain, which is an essential part of the job of a developer
anyway.
`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/Applying DRY: intrinsically vs. coincidentally similar code.webp",
		266, 639,
		"Japanese watercolor of a bird on a branch",
	))),
)

var proverb_designByContract = section(
	"Design by contract",
	markdown(`
*Source: The Pragmatic Programmer*

Check that the input and output of important functions are sensible. This can be
as simple as a few assert statements, e.g.

~~~go
public sendSms(String cell) {
    assert PhoneValidationUtil.ValidNumber(cell);
   // ...
}
~~~

Do this instead of defensive programming, and apply selectively - put it on your
core business logic for specific variables; resources.

By doing this, you'll help minimize unexpected costly mistakes where bad data
gets through to your business logic.

`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/Design by contract.webp",
		266, 400,
		"Japanese watercolor of an outset to a forest",
	))),
)

var proverb_implChecks = section(
	"Implementation checks",
	markdown(`
Go has *structural typing*, so a type implements an interface if it defines the
same method signatures.

While this can be useful for abstracting external types, it is often nice to be
able to check at compile time that a type implements an interface.
	
A nice pattern I've seen for this is to assign an unnamed global variable of the
interface type to an implementation instance, e.g.

~~~go
type Hashable interface {
    Hash() []byte
}

type BasicThingy string
type StructThingy struct{}

// Implementation checks:
var _ Hashable = (BasicThingy)("")
var _ Hashable = &StructThingy{}
~~~

Here we check that \'BasicThingy\' and \'*StructThingy\'
implement \'Hashable\'. If they don't, the compiler will generate a nice error
message pointing at the assignment. For this reason, it is nice to put this
assignment next to the implementation type declaration.

`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/Implementation checks.webp",
		266, 639,
		"",
	))),
)

var proverb_dontNeedDep = section(
	"You don't need a dependency injection library",
	markdown(`
If you follow the Service pattern from DDD, then you'll have a number of
services which transitively depend on each other. Typically, you'll use some
dependency injection container or wiring library to have these services be
connected upon startup.

But there's no reason you can't do this wiring yourself:

~~~go
func main() {
    server := Wire()
    server.Run()
}
  
func Wire() *http.Server {
    cfgProvider := config.NewProvider()
    svc := usecase.NewService(cfgProvider)
    return http.NewServer(cfgProvider, svc)
}
~~~

Advantages:
* You don't need to depend on an external library
* If it compiles then you know everything is wired
* You can have custom Wire functions which e.g. substitute external client
  services with mock versions for integration tests.

It is very convenient for small projects.
`))

var proverb_toolCheck = section(
	"Tool check and install",
	markdown(`
Often as part of your Makefile, you'll have directives which run dev tools to do
things like lint your code, generate mocks, etc. What you may not be aware of is
that - provided your dev tool can be installed via the shell - you can automate
the installation of these tools when the directives are run for the first time.
This makes it trivial for new starters on the repo to get up and going.

Here's an example for a Go tool:

~~~make
GOBIN := $(shell go env GOPATH)/bin

inspect: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run
  
$(GOBIN)/golangci-lint:
	$(MAKE) install-golangci-lint
install-golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.46.2
	rm -rf ./v1.46.2
~~~

If you run inspect, Make will check if a file at \'$(GOBIN)/golangci-lint\'
exists. If it does not, the associated directive will be executed, which then
leads to the \'install\' directive being executed, and finally whatever
shell commands are needed to install the program are run.

Another benefit of this abstraction is that if you update the tool, you can just
advise the team to run \'make install-golangci-lint\' manually.
`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/Tool check and install.webp",
		261, 754,
		"",
	))),
)
