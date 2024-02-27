package site

var BiographyPage = page(rootTmpl, "biography", root(
	"Biography",
	"Biography of Liam Pulles, speaking to key experience and giving contact links.",
	article("Biography", nil,
		section("",
			markdown(biography_body),
			withAsideFigure(figure("", image("profile.jpg", 300, 300, ""))),
		),
	),
	withConnectFooter,
))

const biography_body = `
I am a software developer currently working in **Johannesburg,
South Africa** for the **[Valenture Institute](https://www.valentureinstitute.com/)**.
I completed my honors in Computer Science at the University of the Witwatersrand
(my research focusing on demosaicking and decorrelation stretch methods) in 2017.

In terms of languages:
* I'm most familiar with **Go**, **Java**, and **Python**.
* I have moderate experience with **Shell** scripting, **Javascript**, and
  **Typescript**.
* I have some experience with **CSS**, **Elixir**, and **C**.

I have studied **Image Processing**, **Computer Vision**, **Robotics**, and
**Machine Learning** - among other subjects.

I am familiar with **Kubernetes**, **Docker**, **Spring**, and **Ansible** -
among other technologies and frameworks.

My hobbies include working on small coding and Linux related projects, as well
as viewing a variety of films (I suppose I could be called a "film buff").`
