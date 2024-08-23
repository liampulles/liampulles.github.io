package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"moving-blog",
		`My journey to create a static site generator`,
		"Recounts my reasoning and experiences in moving away from Jekyll towards my own Go based static site generator",
		civil.Date{Year: 2024, Month: time.February, Day: 20},
		markdown(moving_opening),
		"markdown-nightmare.jpg",
		moving_jekyllLimits,
		moving_acceptCriteria,
		moving_firstStab,
		moving_widgetDSL,
		moving_conclusion,
	))
}

const moving_opening = `
My site has been generated using Jekyll since its inception. Jekyll has
served me well, but I decided recently to
write my own static site generator.

Below you'll read about: why I did this,
my initial false start, and my current
widget based approach written in Go.`

var moving_jekyllLimits = section("Jekyll's Limitations", markdown(`
I'm kind of anal retentive when it comes to the HTML of this site. It really
irks me if I have to use classed \'div\'s instead of more appropriate tags.

In particular, I really wanted to make use of HTML5 \'articles\',
\'sections\', and \'asides\'. This is tricky in Jekyll: Jekyll works by
converting markdown files into HTML, and applying them to the appropriate HTML
page template for the file category.

Here's a super basic example for a Jekyll blog post template:

~~~html
---
layout: default
---
<h1>{{ page.title }}</h1>

<article>
    {{ content }}
</article>
~~~

Okay so we've got the \'article\' at least. But the \'{{ content }}\'
part is not
something we can directly template, its just whatever HTML corresponds to
the markdown. So I couldn't have my \'sections\' and \'asides\' in there.

Well not directly anyway - strictly speaking, markdown is in some ways a
superset of HTML - so you can write actual HTML in the markdown files and get
that out as the \'{{ content }}\'. And believe it or not, that is exactly what
I did - such is the degree to which I dislike using classed divs.

But I really want to write markdown (or some kind of
human-oriented markup language)
instead. Also, if I ever changed my mind about what kind of structure to use - I
would've had to have gone and re-structured all my posts by hand (yuck!).

A smaller issue I had with Jekyll was that it uses Ruby to manage its
versioning and installation. I've nothing against Ruby - I'm just not familiar
with it. So if I ever left the site alone
for a few months and came back to it, I had to get back up to speed with Ruby
package management before I could really get going again.
`),
	withAsideFigure(optionalFigure("The Jekyll logo", image(
		"jekyll-icon.png",
		562, 1024,
		"The Jekyll logo: a test tube with red bubbling liquid.",
	))),
	withAsideFigure(figure(`Generated image: "A programmer stares out the window. He sees div HTML tags flying around in a dystopian nightmare. Surreal."`, image(
		"dystopian-window.jpg",
		512, 512,
		`Generated image: "A programmer stares out the window. He sees div HTML tags flying around in a dystopian nightmare. Surreal."`,
	))),
)

var moving_acceptCriteria = section("Acceptance criteria", markdown(`
For my own solution to really work *long term*, it needs to do the
following:
1. I need to be able to write 95% of the actual body of blog posts
   in a **human-friendly**
   markup language (like Markdown, or reStructuredText, etc.).
1. The HTML required to structure the pages needs to be entirely
   under **my control**.
1. BUT I also need to be able to create **rules** for how to 
   structure blog posts,
   etc. for consistency and to guide my writing.
1. I like my site's current design, and it should not change.
1. It must not be so complex that I'll need to spend a lot of time getting back
   upto speed with it after a few months away. I.e. **keep it simple**.
`),
	withAsideFigure(optionalFigure(`Generated image: "A man barking orders at a bunch of computer programs sitting at the other side of a boardroom"`, image(
		"man-barking-orders.jpg",
		512, 512,
		`Generated image: "A man barking orders at a bunch of computer programs sitting at the other side of a boardroom"`,
	))),
)

var moving_firstStab = superSection("Failed first attempt",
	markdown(`
My initial thought process was: Could I add my own syntax
elements to markdown (or
reStructuredText) such that I could parse and generate the
HTML structure of the site according to my rules?

So I looked at what good libraries there were in Go for parsing 
different markup languages.
Probably the most popular is [goldmark](https://github.com/yuin/goldmark), and
fortunately goldmark does support extensions, but it
seemed complex to write these extensions to do what I wanted.

And then I thought: how hard could it be to
write my own markdown parser?`),

	section("Markdown, what art thou?", markdown(`
It turns out that there is no such thing as **a** markdown spec. There are **many**
markdown specs. Markdown is one of those formats which has several
actively used variants, which are not fully inter-compatible.

The closest thing to a commonly agreed spec is (appropriately named)
[CommonMark](https://spec.commonmark.org/current/). My issue with CommonMark
is that it is a very tolerant spec which, while making things easy for
humans, makes parsing quite difficult (for me at least).

In principle, markdown consists of \'block elements\' and \'inline elements\'.
Block elements include things like paragraphs, lists, and fenced code blocks.
Inline elements are for things like italics or links.
Now block elements should ideally be separated by blank lines, but
CommonMark tolerates
certain exceptions - for example you need not separate a paragraph and a list by
a blank line.

This is great for writers, but makes segregating different block
parts somewhat tricky for the parser.

Even ignoring such tolerations, there were a lot more complexities then I had
initially accounted for. For example: a proper markdown parser needs to allow
for things like fenced code blocks in a list, which means you have to
maintain a tree-like
structure to be able to
parse individual lines.

Suffice to say that I did not get far in writing
my own parser, and I decided to
abandon the custom markup approach.
`),
		withAsideFigure(optionalFigure("The CommonMark logo", image(
			"commonmark-icon.png",
			324, 206,
			"The CommonMark logo",
		))),
		withAsideFigure(figure(`Generated image: "A programmer is hallucinating. He sees the logos of various markdown implementations floating above him. He screams."`, image(
			"markdown-nightmare.jpg",
			512, 512,
			`Generated image: "A programmer is hallucinating. He sees the logos of various markdown implementations floating above him. He screams."`,
		))),
	),
)

var moving_widgetDSL = section("Widget DSL", markdown(`
What I've settled on is a kind of widget building / markdown mix. I say widget
because the way the layouts are built is similar to the way it is done in
Flutter and Android Native, i.e. composing nested UI elements.

Here is a truncated example of a digital restoration post:

~~~go
func init() {
	DigitalRestorations = append(DigitalRestorations, digitalRestoration(
		"some-short-url",
		"Some title",
		"Some meta SEO description",
		// Date written
		civil.Date{Year: 2020, Month: time.October, Day: 1},
		image("some-image.png", 7026, 9933, "Some alt text"),
		some_desc,
	))
}

var some_desc = markdown(\'
This bit is just standard markdown. e.g.
* Here's
* a
* list
\')
~~~

Notice how we are appending to a globally defined \'DigitalRestorations\' slice
using the output of the \'digitalRestoration\' function. This function
constrains the task of making a digital restoration post to providing a few key
parameters (which is exactly what I want), while still allowing me
to write the actual "meat"
of the post in markdown. If the definition is getting too big, I can also split
off the nested elements into a separate var, 
(as I have done with \'some_desc\').

This is what the \'digitalRestoration\' func is doing:

~~~go
func digitalRestoration(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
	image Image,
	description template.HTML,
) DatedPost {
	t := date.In(time.Local)
	page := page(rootTmpl, short,
		root(title, seoDesc,
			article(title, mul(
				withHeaderContent(markdown(fmt.Sprintf("*Written %s*",
					t.Format("2 January 2006")))),
				withRawContent(restorationPage(image, description)),
			)),
			withCommentsFooter(short),
		))
	return DatedPost{
		Page: page,
		Date: t,
	}
}
~~~

So it is making use of more general helper funcs to build a \'DatedPost\'
struct. We can see that a digital restoration on my site is its own \'page\',
consisting of an \'article\' body with a header, some raw \'restorationPage\'
content
(which is ultimately defining a HTML5 \'section\') as well as a comments footer.

Using these more general helper funcs means it is quite easy for me to define
a new type of element on my site, while keeping these elements well defined.

Finally, this is an example of what the HTML templating
looks like for an article
(other definitions can get more boiler-platey, so I'll just show you this):

~~~html
{{define "article"}}
<article>
    <header>
        <h1>{{.Header}}</h1>
        {{.HeaderContent}}
    </header>
    {{.RawContent}}
    {{range .Sections}}
        {{template "section" .}}
    {{end}}
</article>
{{end}}
~~~

The template has a name (\'"article"\') and is implanting fields from a
data structure (e.g. \'{{.Header}}\').
We can also do things
like loops or if conditions, etc.

The template is delegating out the templating of a \'"section"\'
to another definition. This allows for reuse and keeps the templates
somewhat manageable.

The widget approach also gives a very nice congruency, in that
each template definition has its own corresponding struct and
widget building func. Indeed there is an \'Article\' struct and \'article\'
func in the code, and similar for \'Section\'/\'section\'.
`))

var moving_conclusion = section("Conclusion", markdown(`
I'm quite pleased with the final result. It is both easy to write a
blog post in the exacting structure I want, but also very easy for me to
add radically new concepts to the site and to change things.

More than that, it has also served as a great learning experience for
Go templating and web design. Suddenly my site is fun to play with again!
`),
	withAsideFigure(figure(`Generated image: "The code for a website has come to life. It is a wonderous and gleeful creature composed of text and images and wires. A man is riding this creature."`, image(
		"man-riding-code.jpg",
		512, 512,
		`Generated image: "The code for a website has come to life. It is a wonderous and gleeful creature composed of text and images and wires. A man is riding this creature."`,
	))),
)
