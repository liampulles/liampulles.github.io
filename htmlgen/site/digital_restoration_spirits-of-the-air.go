package site

import (
	"time"

	"cloud.google.com/go/civil"
)

func init() {
	DigitalRestorations = append(DigitalRestorations, digitalRestoration(
		"spirits-of-the-air",
		"Spirits of the Air, Gremlins of the Clouds",
		"Presents a restored version of a Spirits of the Air, Gremlins of the Clouds theatrical poster. Describes the restoration process.",
		civil.Date{Year: 2024, Month: time.August, Day: 23},
		restorationImage(
			"spirits-of-the-air",
			"https://drive.google.com/file/d/1DE3_tOsvRWao9t8yMbPWtzIeeLXv7spr/view?usp=sharing",
			750, 1000,
			"Spirits of the Air, Gremlins of the Clouds poster",
		),
		spirits_of_the_air_desc,
	))
}

var spirits_of_the_air_desc = markdown(`
*Spirits* is an interesting film... it is quite arty, with strange production design
and an eerie ethereal score. And it has some great posters.

The restoration process:
1. I found the original image online via Google Images. It was at a much lower
   resolution, and with a black background (instead of clouds).
1. I used [Upscayl](https://www.upscayl.org/) (using the digital art model) to bring up
   the resolution **8 times**. This digital art model really is excellent - it did
   introduce a few hallucinated artifacts, but I found these quite interesting
   and so left them in.
1. I removed the black background (around the square-ish middle) to leave it
   transparent, keeping the text.
1. I found a nice high-res cloud image to put on a bottom layer. It is actually a
   sky from the Australian outback, for what its worth.
1. Although I did some feathering when I removed the background, there was still
   a little jaggedness on the outline of the centerpiece I was not comfortable with,
   so I went over the edges with a bit of tactical motion blurring to smooth them
   out.
1. I denoised the cloud backing very slightly - I'd rather lose a little detail
   here then for the noise to conflict with the very smooth interior artwork.

I hired the same printing and framing studio I've used
[previously](/mishima) to do an A0 print on premium matte paper. I'm quite pleased
with the final result, though if I had to pick flaws I'd probably raise the
brightness on the cloud background a bit to bring it closer to "real" sky blue.

It really does feel like it brings more air into my bedroom, somehow. :) 
`)
