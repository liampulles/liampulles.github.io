package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"htmx-snippets",
		"A rant, a static site, and htmx popups",
		"This article is a tutorial on how to use htmx to create popups on static websites. There is some HTML, CSS, and Javascript",
		civil.Date{Year: 2024, Month: time.February, Day: 18},
		markdown(htmxSnip_opening),
		htmxSnip_background,
	))
}

const htmxSnip_opening = `
Let me give some background on the why, and then I will get into the how.`

var htmxSnip_background = superSection("The Rant", `
I (like many others) use a static site generator for my blog. There are many
advantages to this, bonified enough that I won't go
through them here. One
of the big drawbacks though is the lack of dynamic content - i.e.
changing the current page on the fly.`,

	section("SPAs", markdown(`
Well that's not true actually. I used to think this, because I started my career
by using [Angular](https://angular.io/) to build dynamic websites, and I kind of assumed
that that was the "correct, modern, goodboy" way to do it. This was strange,
given that every time I've worked with Angular, it seemed to me to always be a 
distributed-monolith, injection-hell nightmare of a system to work with.

Okay maybe Angular was an unlucky one to begin with. I've heard others talk about
[React](https://react.dev/), [Vue](https://vuejs.org/), and
[Svelte](https://svelte.dev) but I've never used
them and so cannot comment. However I've come
to think (rightly or wrongly) that SPAs are really only a good
idea if, due to other
stronger forces, you must serve JSON APIs from your backend.`)),

	section("Phoenix LiveView", markdown(`
I've found in my own experimentation that
[Phoenix LiveView](https://www.phoenixframework.org/)
is a good all-encompassing; frontend & backend solution. What is really remarkable
about it is how it solves the problems of a
user client interacting with an
eventually consistent system.

Instead of waiting for a response, ignoring the response, or polling, LiveView
allows you to just
submit commands and carry on to leave a later WebSocket-sent event to update
the view.

This approach is both efficient and tolerant of distributed system concerns such
as node failures, rolling releases, asynchronous message-driven service
interaction, etc.`)),
)
