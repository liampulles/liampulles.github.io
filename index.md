---
title: Liam Pulles
description: Homepage for Liam Pulles's blog.
---
<article>
    <hr>
    <section class="toc">
        <h1>Welcome!</h1>
        <p>Hi there - if you're interested in my writing, read on. If you want to hire me (or otherwise find out more about me), then you may wish to see my <a href="/biography">biography</a> or my <a href="/code">code</a>.</p>
    </section>
    <hr>
    <section class="toc">
        <h3>Blog posts</h3>
        <table>
        {% for post in site.posts %}
            <tr>
                <th class="toc-date">{{ post.date | date: '%b %d, %Y' }}</th>
                <th><a href="{{ post.url }}">{{ post.title }}</a></th>
            </tr>
        {% endfor %}
        </table>
    </section>
    <hr>
    <section class="toc">
        <h3>Digital restorations</h3>
        <table>
        {% for post in site.digital_restorations %}
            <tr>
                <th class="toc-date">{{ post.date | date: '%b %d, %Y' }}</th>
                <th><a href="{{ post.url }}">{{ post.title }}</a></th>
            </tr>
        {% endfor %}
        </table>
    </section>
    <hr>
</article>