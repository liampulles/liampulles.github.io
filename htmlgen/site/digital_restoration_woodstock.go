package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	DigitalRestorations = append(DigitalRestorations, digitalRestoration(
		"woodstock-restoration",
		"Woodstock: 3 Days of Peace and Music",
		"Presents a restored version of an old Woodstock festival poster. Describes the restoration process.",
		civil.Date{Year: 2020, Month: time.October, Day: 1},
		image("3-days-of-peace-and-music.png", 7026, 9933, "Woodstock poster"),
		woodstock_desc,
	))
}

var woodstock_desc = markdown(`
The genesis of this restoration was my sisters birthday. Having asked her what
she might like, she noted that she might like a print from a digital poster on
Etsy...

That seemed like a good option to me - So I found this woodstock poster for a
reasonable price, got a digital download link for it and was... horrified.
Though the seller had advertised that the print was suitable for an A1 print,
the cheap scan they gave my was barely suitable for printing on an A4.

Anyway, the image I've restored is suitable for a 300 DPI print for an A1 size
poster. ;)`)
