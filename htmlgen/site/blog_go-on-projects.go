package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	BlogPosts = append(BlogPosts, blogPost(
		"go-in-projects",
		`Using Go in Enterprise Projects`,
		"I give my view on why Go works well in Enterprise projects and software",
		civil.Date{Year: 2025, Month: time.July, Day: 8},
		go_projects_opening,
		"TODO.jpg",
		go_project_monorepo,
		go_project_clis,
	))
}

var go_projects_opening = markdown(`
I've been using Go professionally for a few years now, and I kind of love it.
I'm not going to argue here that Go is the best language ever made (it depends),
but I want to try and distill here a few interesting attributes of it that make
it amenable to writing enterprise code as part of a team.
`)

var go_project_monorepo = section("Easy monorepos", markdown(`
With Go, its possible to have multiple modules and programs live alongside each
other in their own folders. It is literally just a matter of making a new \'main.go\'
file and/or doing \'go mod init xyz\' in a directory.

This means that if I need to write a new web server, I don't have to go through the
ordeal of making a new repository, figuring out where common API files live, and
then coordinating pull requests between multiple repos forever more. I just do
it all in one repository. And compiling in Go is so quick and easy (\'go build ./...\')
that I don't really have to think about the code cost of new services in a repo.
I'm freed up to think about the best architecture, and to implement it easily.
`))

var go_project_clis = section("Easy CLIs", markdown(`
Writing CLIs in Go is incredibly sensible, especially if you care about making
programs which work well with UNIX pipes, because:
* Its very easy to efficiently stream
over STDIN, and to write CLIs that can stream over gigabytes of input.
* The concurrency model simplifies concurrent processing
with well defined rate limits, concurrency limits, etc.
* You can compile the program into a static binary which you can use in
shell scripts easily.

Having the capability to write really powerful CLIs gives one more options of
how to automate things in a project. For example, data migrations or one-off
data loads are easy to do. Irregular tasks (e.g. things related to one date
in the year) are easy to code and can be saved to the repository for future use.
Common support tasks can be simplified.

DevEx can be automated. For example: I wrote a tool
that could dynamically port-forward, load, and modify config for running our
services locally.

Side note: I've found it incredibly useful to have an internal library
that facilitates streaming CSV files to an anonymous row-handling function with
configurable rate limiting, concurrency, etc. This has come in clutch SO MANY
times in emergency situations where masses of data have needed fixing in a short
amount of time.
`))
