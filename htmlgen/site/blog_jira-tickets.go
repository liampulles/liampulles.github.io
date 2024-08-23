package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"jira-tickets",
		"How I use JIRA tickets",
		"This article speaks to how I use JIRA tickets or kanban cards: as a blueprint, as a second brain, and as a form of asynchronous communication.",
		civil.Date{Year: 2023, Month: time.June, Day: 1},
		markdown(jira_opening),
		"thinking-dev.webp",
		jira_ticketsAsABlueprint,
		jira_ticketsAsASecondBrain,
		jira_ticketsAsAsync,
		jira_conclusion,
	))
}

const jira_opening = `
In my work I have noticed that my fellow developers' opinions about JIRA tickets (or Kanban cards, etc.) vary widely -
but many view them as a form of admin drudgery.

When I was starting out, I had this view too - but I've since come to view tickets as a vital part of my workflow.
In fact I'd put them at a similar weight to my IDE in terms of value (seriously). Here I'll show you how I use them
in the hopes that it will give you some ideas.`

var jira_ticketsAsABlueprint = section(
	"Tickets as a blueprint",
	markdown(`
When I pick up a ticket, I'll start by making a series of TODOs on it, and then
continue to refine it down into the finest detail I can. For simple tasks, this is
a 10 minute procedure - for more complex ones, it can easily take an hour (or more).

I usually start by making bullet points which list the high level elements involved
(e.g. this db table, this service, this util, etc). Then I will turn this into a
series of explicit steps (e.g. migrate the column for this table, delete this code,
create a new service class, etc).

The rule I have in mind when doing this is to get to the point where there is no
ambiguity about what needs to be done. I will look through the codebase and talk
to people while I do this.

Some notes on this:

1. Planning upfront like this inevitably leads to better design and less rework.
I'll see that steps conflict and rewrite the plan before rewriting the code.
1. Since I've done the thinking upfront, when I do get to coding, I can get into
"flow" more easily and hold on to it for longer.
1. And since I'm in flow for longer, I find that I complete my work more quickly
(even with the added planning time).
1. As I progress through my TODOs, I mark the individual pieces as [DONE]. This keeps
me motivated, and gives me an accurate indication of progress.
1. I think this has made me better at seeing what is involved in a ticket off the
cuff. This helps immensely during stand up and sprint planning.
1. My ability to estimate tasks has improved. And if the task does go over, this TODO
list serves as a good resource for me to explain *why* I need to go over.
1. If I see any tasks that need to be done first, I'll action those upfront. For
example, I might talk with another team to update their API so they can do that
while I am busy. This obviously helps avoid blocks/stucks.
1. TODOs help immensely with context switching. When I switch tickets, I can read my
TODOs and jump back in to coding. Also if I am waiting on something (e.g. a
deployment), then I can easily find something small from these lists to action.
1. I am hoping that these TODOs act as a historical artifact for future devs (they
can see my reasons for writing such bad code. 😉)
1. I've gotten feedback from my manager and colleagues that these TODOs help a lot
to get a sense of my progress (particularly when there are tasks dependent on my
work). This pre-sight helps them to plan and for me to explain my work.

Note that this does not stop unexpected things from coming up during coding, but it
does help reduce them. And when those moments do come up I try to go and plan again
before coding.`),
	withAsideFigure(figure("", image(
		"thinking-dev.webp",
		512, 512,
		"A software developer leaning back in his chair thinking hard about his work. Thought bubbles, casual clothing, brown hair.",
	))),
)

var jira_ticketsAsASecondBrain = section(
	"Tickets as a second brain",
	markdown(`
I often say (only half joking) that I have the memory of a goldfish. I treat my
brain as an ephemeral cache: any value could fall away at any time. My tickets
serve as a backing database for my brain.

Extending this cache-db analogy: I'll write to that persistent storage regularly.
If, while coding, a disparate question pops into my brain - I'll note the question
down in the ticket and then carry on coding. This keeps me in flow. Later I will
look through the ticket and ensure all items are actioned before submitting for
review.
	
This practice serves to enable a kind of agile Plan-Do-Reflect mode of working.
Instead of just barraging through code, uncertain of where I am going, I am making
steady steps with a clear path ahead. It leads to more inner peace (I think).`),
	withAsideFigure(figure("", image(
		"thinking-goldfish.webp",
		512, 512,
		"A goldfish leaning back in his chair thinking hard about his work. Thought bubbles.",
	))),
)

var jira_ticketsAsAsync = section(
	"Tickets as asynchronous communication",
	markdown(`
I work at a fully remote company, and so good asynchronous communication is vital.
If I have a ticket-related question then I will sometimes ping colleagues on a
ticket rather than slack them. My rationale being:

* They get the notification via email and so are not interrupted (by a slack
  notification bell).
* The context for the question is on the ticket itself.
* That conversation, and all other relevant conversations, are all available on
  the ticket.

Although having said that, people's preferences differ, and one should respect the
communication channels others prefer (especially when asking them for help).`),
	withAsideFigure(figure("", image(
		"thinking-cyberpunk.webp",
		512, 512,
		"A man sitting in his chair hacking into the matrix.",
	))),
)

var jira_conclusion = section(
	"Conclusion",
	markdown(`
Anyway, I hope I've given you something to think about. I'm curious if you have
found other good uses for tickets - if so, please do email me.`),
)
