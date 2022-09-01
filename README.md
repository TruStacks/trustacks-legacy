<p><img src="https://i.imgur.com/awXcDpX.png" height="" alt="trustacks logo"/></p>

---

[![Join Discord](https://img.shields.io/discord/958112026687840257?label=Discord&style=for-the-badge)](https://discord.gg/mbkQz82V)
![Go Report Card](https://goreportcard.com/badge/github.com/trustacks/trustacks?style=for-the-badge)

[![asciicast](https://asciinema.org/a/mkAtIDhoMozEsKMM7a7zFr89z.svg)](https://asciinema.org/a/mkAtIDhoMozEsKMM7a7zFr89z)

# Welcome to TruStacks

## What is TruStacks?

TruStacks is a software delivery engine that enables teams to code with "ultra standardized", framework driven developer workflows.

## Why TruStacks?

[VSDPs](https://devops.com/why-you-need-a-value-stream-delivery-platform/) stop at the pipeline. TruStacks is the write once, consume anywhere approach to software delivery pipelines and developer workflows.

### Community

**It takes a village.**

Software delivery cannot be solved alone. By working together as a community, TruStacks provides end-to-end development workflows from code to delivery for the most popular frameworks across web, mobile, IoT, embedded, and many other software disciplines.

## How Does It Work?

1. Deploy a software toolchain.  
   <small>*The toolchain deploys and configures core tooling such as single-sign (OIDC), CI provider (ie. concourse, gitlab, tekton), and secrets management.*</small>
2. Create an application for a target software framework from a compatible<sup>(1)</sup> project in git.  
   <small>*The application config is appeneded to the toolchain configuration*</small>
3. The TruStacks engine CI driver automically creates and configures the application pipeline.
4. Run your pipeline from the CI proivder.  
   <small>*The CI pipeline is pre-built for your framework. Certain parameters such as deploymnet target (ie. k8s, FaaS) are configuration through application parameters.*</small>

*(1) compability depends on the desired framework. (ie. a `Create React App` based workflow requires an un-ejected react project).*

## Try It Out

[Visit the docs](http://docs.trustacks.io/category/getting-started) to learn more and to get started with TruStacks.

*\*TruStacks is `Alpha Software`. The API and tooling are subject to breaking changes in the future. We will make every attempt to minimize breaking changes as we add more features.\**

## Issues & Suggestions

If you run into any issues [create an issue](https://github.com/TruStacks/trustacks/issues/new0) or drop us a line on [discord](https://discord.gg/kGJBuHdv).

Leave a  [suggestion](https://discord.gg/tgpWURqY).

## What's Happening 

Check out our [Trello Board](https://trello.com/b/IwJMgZiO/trustacks-oss) to see what's happening now.