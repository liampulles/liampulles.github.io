package main

import (
	"time"

	"cloud.google.com/go/civil"
)

// This file is essentially a blueprint for the site. What is defined here
// eventually makes its way into HTML.

// Other details are probably in the render functions.

var IndexPage = PageDefinition{
	Type:           Index,
	Short:          "index",
	Title:          "Liam Pulles",
	SEODescription: "Homepage for Liam Pulles's blog.",
	ExtraData: map[string]any{
		"MainHeader": "Welcome!",
		"MainDescription": InlineMarkdown(`Hi there - if you're interested in my writing, read on. 
		    If you want to hire me (or otherwise find out more about me), then you may wish to see
			my [[biography]] or my [[code]].`),
	},
}

var BlogPosts = []PageDefinition{
	blogPost(
		"jira-tickets",
		"How I use JIRA tickets",
		"This article speaks to how I use JIRA tickets or kanban cards: as a blueprint, as a second brain, and as a form of asynchronous communication.",
		civil.Date{Year: 2023, Month: time.June, Day: 6},
	),
}

// ---
// --- Helpers
// ---

type InlineMarkdown string

type PageType int

const (
	Index PageType = iota
	Nav
	BlogPost
	DigitalRestoration
)

type PageDefinition struct {
	Type           PageType
	Short          string
	Title          InlineMarkdown
	SEODescription InlineMarkdown
	ExtraData      map[string]any
}

func blogPost(
	short string,
	title string,
	seoDesc string,
	date civil.Date,
) PageDefinition {
	t := date.In(time.Local)
	return PageDefinition{
		Type:           BlogPost,
		Short:          short,
		Title:          InlineMarkdown(title),
		SEODescription: InlineMarkdown(seoDesc),
		ExtraData: map[string]any{
			"Date": t,
		},
	}
}
