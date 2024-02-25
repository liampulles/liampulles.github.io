package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	post := blogPost(
		"proverbs",
		`Proverbs`,
		"Contains proverbs and advice around programming and developer life.",
		civil.Date{Year: 2024, Month: time.February, Day: 21},
		markdown(proverb_opening),
		superSection("Code Design", "",
			proverb_applyingDry,
			proverb_minimizeShit,
			proverb_designByContract,
			proverb_testRefactor,
			proverb_commentsAreGood,
		),
		superSection("Peopleware", "",
			proverb_comraderyOverAll,
			proverb_reviews,
		),
		superSection("Golang", "",
			proverb_implChecks,
			proverb_dontNeedDep,
		),
		superSection("Makefiles", "",
			proverb_toolCheck,
		),
	)
	// Its unlisted, in the nav.
	post.Unlisted = true
	BlogPosts = append(BlogPosts, post)
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

var proverb_minimizeShit = section(
	"Aim to minimize shit, not maximize perfection.",
	markdown(`
I think as a junior developer especially, ones desire with code design is to try
and find the absolute perfect abstractions and algorithms, and to over-engineer
solutions.

At least that's what I did as a junior. I would spend too much time trying to
perfect rather unimportant aspects of the system, and be disheartened with
the continuing existence of bad code in our codebase.

And then I decided to flip my perspective.
Rather than trying to make good things
better, I would try to make bad things less terrible.

So as an example: rather than trying to fine-tune a usecase that is already
working okay, I will try to identify important
parts of a
system which are not working well or are of poor quality, and make
a small improvement.

This is an acceptance of the fact that software is the product of flawed
corporations and people. And once you accept that, you have achieved zen,
and you can move on to
being productive.
`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/roman-sewer.jpg",
		768, 1024,
		"Exposed top-view of an ancient Roman sewer",
	))),
)

var proverb_designByContract = section(
	"Design by contract",
	markdown(`
*Source: The Pragmatic Programmer*

Check that the input (and sometimes output) of important
functions are sensible. This can be
as simple as a few assert statements, e.g.

~~~java
public sendSms(String cell) {
    assert PhoneValidationUtil.ValidNumber(cell);
   // ...
}
~~~

Do this instead of defensive programming, and apply selectively: put it on your
core business logic for specific variables.

By doing this, you'll help minimize unexpected costly mistakes where bad data
gets through to your business logic.

`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/Design by contract.webp",
		266, 400,
		"Japanese watercolor of an outset to a forest",
	))),
)

var proverb_testRefactor = section(
	"Test first, refactor second",
	markdown(`
Good tests allows you to delete code with confidence. Good tests are piercing
acceptance or integrations tests which test the main happy paths of the program
as well as known edge cases. Refactoring generally means deleting
and rewriting code.

The lesson: write acceptance tests before doing serious refactoring.
`),
)

var proverb_commentsAreGood = section(
	"Curt comments are good",
	markdown(`
Having the code describe itself is a wonderful, but terribly hubristic idea.
Once you've written and reviewed a piece of code that nonetheless had a bug in
it, you should abandon the silly notion that you can understand code just by
looking at it. Good code has comments in it.

BUT you still shouldn't comment for the sake of it. Here I think are good
opportunities for commenting:
* The steps of a large usecase function. Moving those steps to helper functions
  can help too.
* Any function or piece of code which is working around something; e.g. a badly
  designed external API.
* Assumptions around the form of input which cannot be expressed clearly
  with types (though see [Design by Contract](#Design%20by%20contract)).
* Any code which is logically related but where the relationship 
  cannot be enforced by the code (e.g. needing a struct definition to stay in
  sync between a client and repository package).
* An overview for a complex or relatively alien algorithm.

`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/epictetus.jpg",
		339, 600,
		"18th-century portrait of Epictetus, including his crutch",
	))),
)

var proverb_comraderyOverAll = section(
	"Prioritize people over architecture",
	markdown(`
Comradery is the most vital resource on a software project. It is more important
than having a good architecture or a clean, well designed codebase.

Consider two hypotheticals. In the first case
you have a well-architected, clear, and performant system - but with a team that
does
not trust itself, struggles to make decisions, and which has endless bad
arguments about all-manner of minutia.

Then consider a second case. The codebase is full of spaghetti code, takes
hours to deploy, and is relatively unclear - but the team embraces its
constituent members, people jump in and help each other at a moments notice,
decisions are made dispassionately and quickly, and individuals
are largely empowered to decide how to code their own tickets as they see fit.

My bet is that the first team is going to be relatively unproductive in real
terms (that is, being slow to deliver the *right* business functionality)
and the second team will conversely be strangely effective. And the second team will
be able to pragmatically refactor its system over time, whereas the first team
will lack the openness to discuss and address real arising problems with their
system.
`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/sailors.jpg",
		479, 600,
		"Propaganda artwork depicting two sailors sharing a drink",
	))),
)

var proverb_reviews = section(
	"Corollary #1: Major review comments are for broken code only",
	markdown(`
Don't make major, blocking code review comments if it is a difference of
opinion, even if you are only doing so to force a discussion. I would
even say you shouldn't make a major comment if something  is violating the
architecture, but still works.

Major comments should only be used if there is a bug that will cause business
requirements to fail.

This is based on the ideas of comradery
being key, the architecture of a system being a living thing,
and that people should be
empowered to do their work as they see fit (provided their solution works).

You can use minor comments though. Make it clear to the team that your nits are
purely optional and they need not even be read.

I'm quite a rigorous reviewer, and this balance means I generally give 1 or
two major comments and a dozen (sometimes several dozen) minors.
I find that 90% of the time the
reviewee will happily correct all the majors and most of the minors without
any argument. Thats fine with me.`),
	withAsideFigure(optionalFigure("", image(
		"proverbs/scholar-reading.jpg",
		513, 600,
		"Dutch painting of a scholar reading in study",
	))),
)

var proverb_implChecks = section(
	"Implementation checks",
	markdown(`
Go has *structural typing*, so a type implements an interface merely
by defining the same methods on itself.

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
