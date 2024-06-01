package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	DigitalRestorations = append(DigitalRestorations, digitalRestoration(
		"2001-restoration",
		"2001: A Space Odyssey",
		"Presents a restored version of an 2001 theatrical poster. Describes the restoration process.",
		civil.Date{Year: 2024, Month: time.June, Day: 1},
		restorationImage("2001", 8938, 13862, "2001 poster"),
		twoThousandOne_desc,
	))
}

var twoThousandOne_desc = markdown(`
Lately I've been looking to liven up my walls a bit with some art.
There is a professional and reasonably priced [printing shop](https://www.pictureframingstudio.co.za/)
local to me who I've previously used to print a *Monkey Island* poster, so I
thought I should enlist them to print a film poster for me.

I scoured google images for a high resolution and striking poster for a film I
like and discovered the above 2001 poster. The main changes I made were to:

* Run it through [Upscayl](https://www.upscayl.org/) (using the Ultra
  sharp model)
* Level the image out so it ranges from full black to full white
* Redo the white border
* Apply a fake grain, which consists of a couple layers of motion blurred
  random noise
* White out the text and MGM logo

There is enough detail here to get a very sharp A1 print, you can possibly
stretch it even to an A0. I had mine printed and framed on an A1 size with a
speckled glossy paper, which worked out very nicely.

It did come out a bit darker than I expected - you might consider raising the
mid-level a bit as the star child is a difficult to make out
in my print.
`)
