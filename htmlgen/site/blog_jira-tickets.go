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
	))
}
