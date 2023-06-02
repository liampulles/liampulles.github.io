---
layout: post
title: "How I use JIRA tickets"
description: "This article speaks to how I use JIRA tickets or kanban cards: as a blueprint, as a second brain, and as a form of asynchronous communication."
shareable: true
permalink: /blog/jira-tickets
---
<section>
    <p>In my work I have noticed that my fellow developers' opinions about JIRA tickets (or Kanban cards, etc.) vary widely - but many view them as a form of admin drudgery.</p>
    <p>When I was starting out, I had this view too - but I've since come to view tickets as a vital part of my workflow. In fact I'd put them at a similar weight to my IDE in terms of value (seriously). Here I'll show you how I use them in the hopes that it will give you some ideas.</p>
</section>

<section>
    <h2>Tickets as a blueprint</h2>
    <aside>
        <figure>
            <img src="/images/thinking-dev.webp" height="512" width="512" alt="A software developer leaning back in his chair thinking hard about his work. Thought bubbles, casual clothing, brown hair.">
        </figure>
    </aside>
    <p>When I pick up a ticket, I'll start by making a series of TODOs on it, and then continue to refine it down into the finest detail I can. For simple tasks, this is a 10 minute procedure - for more complex ones, it can easily take an hour (or more).</p>
    <p>I usually start by making bullet points which list the high level elements involved (e.g. this db table, this service, this util, etc). Then I will turn this into a series of explicit steps (e.g. migrate the column for this table, delete this code, create a new service class, etc).</p>
    <p>The rule I have in mind when doing this is to get to the point where there is no ambiguity about what needs to be done. I will look through the codebase and talk to people while I do this.</p>
    <p>Some notes on this:</p>
    <ol>
        <li>Planning upfront like this inevitably leads to better design and less rework. I'll see that steps conflict and rewrite the plan before rewriting the code.</li>
        <li>Since I've done the thinking upfront, when I do get to coding, I can get into "flow" more easily and hold on to it for longer.</li>
        <li>And since I'm in flow for longer, I find that I complete my work more quickly (even with the added planning time).</li>
        <li>As I progress through my TODOs, I mark the individual pieces as <code>[DONE]</code>. This keeps me motivated, and gives me an accurate indication of progress.</li>
        <li>I think this has made me better at seeing what is involved in a ticket off the cuff. This helps immensely during stand up and sprint planning.</li>
        <li>My ability to estimate tasks has improved. And if the task does go over, this TODO list serves as a good resource for me to explain <i>why</i> I need to go over.</li>
        <li>If I see any tasks that need to be done first, I'll action those upfront. For example, I might talk with another team to update their API so they can do that while I am busy. This obviously helps avoid blocks/stucks.</li>
        <li>TODOs help immensely with context switching. When I switch tickets, I can read my TODOs and jump back in to coding. Also if I am waiting on something (e.g. a deployment), then I can easily find something small from these lists to action.</li>
        <li>I am hoping that these TODOs act as a historical artifact for future devs (they can see my reasons for writing such bad code. 😉) </li>
        <li>I've gotten feedback from my manager and colleagues that these TODOs help a lot to get a sense of my progress (particularly when there are tasks dependent on my work). This pre-sight helps them to plan and for me to explain my work.</li>
    </ol>
    <p>Note that this does not stop unexpected things from coming up during coding, but it does help reduce them. And when those moments do come up I try to go and plan again before coding.</p>
</section>

<section>
    <h2>Tickets as a second brain</h2>
    <aside>
        <figure>
            <img src="/images/thinking-goldfish.webp" height="512" width="512" alt="A goldfish leaning back in his chair thinking hard about his work. Thought bubbles.">
        </figure>
    </aside>
    <p>I often say (only half joking) that I have the memory of a goldfish. I treat my brain as an ephemeral cache: any value could fall away at any time. My tickets serve as a backing database for my brain.</p>
    <p>Extending this cache-db analogy: I'll write to that persistent storage regularly. If, while coding, a disparate question pops into my brain - I'll note the question down in the ticket and then carry on coding. This keeps me in flow. Later I will look through the ticket and ensure all items are actioned before submitting for review.</p>
    <p>This practice serves to enable a kind of agile Plan-Do-Reflect mode of working. Instead of just barraging through code, uncertain of where I am going, I am making steady steps with a clear path ahead. It leads to more inner peace (I think).</p>
</section>

<section>
    <h2>Tickets as asynchronous communication</h2>
        <aside>
        <figure>
            <img src="/images/thinking-cyberpunk.webp" height="512" width="512" alt="A man sitting in his chair hacking into the matrix.">
        </figure>
    </aside>
    <p>I work at a fully remote company, and so good asynchronous communication is vital. If I have a ticket-related question then I will sometimes ping colleagues on a ticket rather than slack them. My rationale being:</p>
    <ul>
        <li>They get the notification via email and so are not interrupted (by a slack notification bell).</li>
        <li>The context for the question is on the ticket itself.</li>
        <li>That conversation, and all other relevant conversations, are all available on the ticket.</li>
    </ul>
    <p>Although having said that, people's preferences differ, and one should respect the communication channels others prefer (especially when asking them for help).</p>
</section>

<section>
    <h2>Conclusion</h2>
    <p>Anyway, I hope I've given you something to think about. I'm curious if you have found other good uses for tickets - if so, please do email me.</p>
</section>