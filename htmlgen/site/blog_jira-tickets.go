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
		civil.Date{Year: 2023, Month: time.June, Day: 6},
		markdown(
			`In my work I have noticed that my fellow developers' opinions about JIRA tickets (or Kanban cards, etc.) vary widely -
			but many view them as a form of admin drudgery.
			
			When I was starting out, I had this view too - but I've since come to view tickets as a vital part of my workflow.
			In fact I'd put them at a similar weight to my IDE in terms of value (seriously). Here I'll show you how I use them
			in the hopes that it will give you some ideas.`,
		),
		nil,
	))
}
