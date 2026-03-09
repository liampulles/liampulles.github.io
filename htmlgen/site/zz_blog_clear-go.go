package site

// func init() {
// 	BlogPosts = append(BlogPosts, blogPost(
// 		"go-story-enrollment",
// 		`Enrollment: A Go story.`,
// 		"A story about a failed enrollment and how Go saved the day.",
// 		civil.Date{Year: 2025, Month: time.October, Day: 14},
// 		markdown(clearGo_opening),
// 		"thinking-dev.webp",
// 		clearGo_section,
// 	))
// }

const clearGo_opening = `
*This is an entirely fictional story. Never happened. Nope.*
`

var clearGo_section = section(
	`Revelations`,
	markdown(`
> Jane: So no one is enrolled?

> Bob: None. Because the Albatross enrollment is failing, none of the subsequent
integration services were triggered.

> Jane: Not one of 20000 users?

> Bob: None. They're all stuck in a failing retry loop in NATS. And we can't enroll
anyone into Intercom or Google Workspace or MyEMS until they've been enrolled
in Albatross, because Albatross is responsible for figuring out those API calls.

> Jane: Why is Albatross not working? Can we fix it?

> Bob: Its broken in multiple places, because we haven't catered for real production
data. We tested this in staging using last year's data, and finance (as well as HR)
have now changed the groups they're using in prod for 2015. 
It would take some time for me to
find those groups and to add them to our system, but having not tested them, I
don't really know how the system will act when it tries to enroll them.

> Jane: So no one will be able to login tomorrow morning?

> Bob: Correct.

...

> Alex: Okay. Putting the process aside for a moment, lets ask what the bare 
minimum result is
here. We don't need users in Intercom for now - there is a support mail link in
MyEMS. What else can go?

> Bob: We don't need everything up to date in Google Workspace, but we do need all
our users there so they can sign into MyEMS using single sign on.

> Alex: Could we switch everyone to email/password login, hypothetically?

> Bob: ... Hypothetically, we could turn that setting on for 
everyone in the DB, and I could
write a script to trigger a forgotten password email for everyone. But the users
not currently on GW will fail to load into MyEMS, because MyEMS needs to fetch
some data from GW for page rendering.

> Alex: Damn, okay. What is the bare minimum we need to update on MyEMS, besides
getting our news users on?

> Jane: I figure that we don't need to update the finance groups off the bat, we
have two weeks before salaries go out and that is the only absolutely critical
finance need of MyEMS.

> Bob: I agree, but note that everyone needs to be in *some* finance group
off the bat, otherwise page loads will fail. In terms of HR groups, those dictate
line manager relationships. Some of those relationships are now invalid due to
employee movements and departures, and we do need to update them.

> Alex: And if we just give everyone one catch-all HR group for now?

> Bob: ... I can't think of a reason right now why that wouldn't work, but I do
think there may be unexpected behaviour from everyone being in one group,
particularly for the line manager.

> Alex: And if we just make me the line manager, and I don't login for now?

> Bob: Its risky, but I can't think of a better 80/20 solution.

> Alex: Okay, lets go with that for now.

> Jane: It occurs to me that we don't need *everyone* ready for tomorrow. A good
portion of the company will still be on holiday.

> Bob: Yes but- well actually, I guess if Alex is everyone's line manager, then its
okay if we only update and enroll some of the people.

> Alex: Great - are there any other side effects of doing a partial enrollment?

> Bob: God knows. Its probably safest to disable the profiles of whoever we don't
enroll. In principle that should eliminate almost all interactions between them
and the rest of the system.

> Alex: Okay, lets do that.

> Jane: Frankly, there are also some people who are more important to enroll then
others. The executive team are the first priority, and our janitorial staff
don't even use MyEMS, for example.

> Alex: True. Okay, Bob: do you think you could write something to do all
this?

> Bob: To be honest, I haven't kept track of everything we discussed...

> Jane: I've written it all down in a ticket. I'll share a link with you.

> Alex: Thank you Jane. Bob?

> Bob: ...can you give me an hour to think about it and we'll chat later?

> Alex: Okay lets do that. I need to put the kids to bed anyway. See you all
again at 21:00?

*Everyone agrees and signs off the Google Meet call.*

Bob rest his headphones down on the desk, walks a few steps to the kitchen, 
and turns the kettle on. Then he waddles over to the bathroom, and craters his
skull into his palms as he sits on the porcelain.
`),
	withAsideFigure(figure("", image(
		"thinking-dev.webp",
		512, 512,
		"A software developer leaning back in his chair thinking hard about his work. Thought bubbles, casual clothing, brown hair.",
	))),
)
