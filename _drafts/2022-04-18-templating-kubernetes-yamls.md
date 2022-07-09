---
layout: post
title: "Reinventing the wheel: Build your own Kubernetes templating app."
shareable: true
---
<section>
    <p>So you've got a Kubernetes cluster running, and a repository (or several) with spec YAMLs to deploy (perhaps using Flux(?)) - but you're at a point where you have to constantly update multiple YAML files whenever a small piece of your DevOps environment changes? If so, then congratulations - you're reading the right article. ;)</p>
</section>

<section>
    <h2>Common solutions to a common problem</h2>
    <p>You're not alone in having this problem, and it is well serviced with a few existing solutions. Most notable among these are Helm and Kustomize, but there are many other options as well.</p>
    <p>Why the title of this article then? Well, I don't think any of these projects are as good as an app you build for yourself - hear me out...</p>
    <p>The question of whether to config an open source solutions vs building your own is (as always) a matter of tradeoffs. In general, the open source option is preferable, because it has all the features you really need and requires you to expend minimal effort in configuring it.</p>
    <p>But with this particular problem, I think the scales tip the other way, and this really boils down to:</p>
    <ol type="a">
        <li>Writing code to template files is a very well established thing; not hard.</li>
        <li>Using your own code to specify Kubernetes specs opens up a lot of interesting DevOps opportunities, which we will get to later.</li>
    </ol>
    <p>If you're convinced or at least intrigued, hang on and I'll show you how to implement a a basic k8s templating app.</p>
</section>

<section>
    <h2>What we will build</h2>
    <p>For the sake of brevity, we'll build a relatively simple app which will read our own custom microservice YAML and generate the associated deployment, service, and config YAMLs.</p>
    <p>With this type of application, you can start small, and that is what I'd recommend - take an SRE approach where you automate the generation and updating of the most mundane, known, repeatable configs and build on it over time.</p>
    <p>I'll be building this app in Go, but all languages have decent templating libraries so don't feel constrained - use whatever makes sense for your team and project.</p>
</section>

<section>
    <h2>Out custom spec</h2>
    <p></p>
</section>