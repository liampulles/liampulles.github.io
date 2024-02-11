# General Software Development Cheatsheet
This is a totally subjective, terse and yet broad set of ideas, rules, snippets, and references for my own use. I use this for problem solving, designing, coding, and reviewing.

The layout here is meant to be optimal for scanning, and easy to Ctrl+F, yet also allows for expansion should I need to remind myself of some context.

If this prompts any confusion or ideas that you'd like to discuss - pop me an email at the link at the bottom of the page.

## Design

### Duplication/Abstraction
Is it intrinsically or coincidentally similar? Imagine potential changes and risk of seperate code vs burder of joined code.

General rule: The more "core" the logic is, the more likely to be intrinsically similar it is.

## Code review checklist

Only make it major if the code is broken or not solving the problem.

Adress the code, not the person. Assume best intent.

Try to use a review mode rather than comments.

* Is new code unused? Has old code become unused?
* Are there too many nested ifs; loops; etc.? Can the ifs be inverted, or can functions be split out?
* Should a func be private?
* Are the naming of variables and "knowledge" of the code appropraite to the domain?
* Is there code that needs to be tested?
* Can code be simplified?
* Is it obvious at a glance what the intent of the code is?
* Is there a lot of calling to and fro from another class?
* Can calls for data be done in bulk?
* Is the order of performance as good as it could be?

## Career

### Assessing projects
