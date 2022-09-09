---
title: Toolchain Installation
slug: /tutorial/install-toolchain
sidebar_position: 2
---

# Software Toolchain

The toolchain contains the software components required to build and release applications.

## Configuration

To install a new toolchain we need to create a toolchain configuration.

Create a file named `react-tutorial-config.yaml` and add the following values:

```yaml
name: react-tutorial
source: https://github.com/trustacks/toolchain
parameters:
  sso: authentik
  ci: concourse
  network: private
  ingressPort: "8081" # change this value if you used a different port for your k3d loadbalaner
  tls: "false"
```

Let's break down the configuration values:

> `name` is the name of the toolchain.  
> `source` is the the git repository that contains the toolchain resources.  
> `parameters` are [values](https://github.com/TruStacks/catalog/blob/main/pkg/catalog/catalog.yaml) that will be passed to the software components during installation.

:::tip

Configuration parameters must be strings. Ensure that numbers and booleans are string quoted.

:::

:::caution in case you missed it

If you changed your ingress port during the [k3d cluster creation](http://localhost:3000/installation#cluster-creation) to something other than `8081`, make sure to update `parameters.ingressPort` in the yaml configuration file before proceeding.

:::

## <a id="core-components-aag"></a> Core Components At-A-Glance

TruStacks toolchains must include a set of core components as dependencies.

The core componets are:

- CI Provider
- SSO Provider

Toolchains can include any number of supplemental components that provide a wide range of software delivery capabilities.

The toolchain used in this tutorial uses [authentik](https://goauthentik.io/) for SSO and [concourse](https://concourse-ci.org/) for CI.

:::tip don't like the tools?

TruStacks Engine is built to be extensible. If it can be deployed with helm, then it can likely be built into the TruStacks eco-system.

Drop us a suggestion on [discord](https://discord.gg/tgpWURqY), or check our project board to see [what's happening now](https://trello.com/b/IwJMgZiO/trustacks-oss).

:::

:::info

Parameters, Core Components and other in-depth topics are covered in greater detail in the Core Concepts.

:::


## Install

Now that we have our configuration file we can install our toolchain.

Run the following command to install the toolchain:

    tsctl toolchain install --config react-tutorial-config.yaml

Check the status of the services with the following command. Wait until all service are in the `Running` state:

    kubectl get po -n trustacks-toolchain-react-tutorial  

Example output:

    NAME                                READY   STATUS    RESTARTS   AGE
    dind-88db84bd6-8bxxf                1/1     Running   0          91s
    authentik-worker-78c5d654c4-jlv45   1/1     Running   0          87s
    authentik-postgresql-0              1/1     Running   0          87s
    authentik-redis-master-0            1/1     Running   0          87s
    authentik-server-786cb79b-xq76s     1/1     Running   0          87s
    concourse-worker-0                  1/1     Running   0          22s
    concourse-worker-1                  1/1     Running   0          22s
    concourse-postgresql-0              1/1     Running   0          22s
    concourse-web-747c56c56f-b94ql      2/2     Running   0          22s

:::caution air gapped installation

Air-Gapped environemnts are not currently supported.

:::

## SSO Setup

Navigate to [authentik](http://authentik.local.gd:8081/if/flow/initial-setup/) to configure the `akadmin` user.

<img src="/img/authentik-initial-setup.jpg" />

You should see the following page after navigating to the address. Enter an email and password and click continue.

After clicking continue you will be taken to the authentik landing page. The sso provider configuration is now complete.

:::tip changed port?

Navigate to `http://authentik.local.gd:<port>/if/flow/initial-setup/` if you changed your loadbalancer port.

:::

## CI Provider

Now that your sso provider is configured, you can log in to [concourse](http://concourse.local.gd:8081).

<img src="/img/concourse-login.jpg" />

Click the login page at the top right of the page or click the login link in the hero image.

<img src="/img/concourse-sso.jpg" />

Click `sso` to sign in using the sso provider.

:::tip changed port?

Navigate to `http://concourse.local.gd:<port>/if/flow/initial-setup/` if you changed your loadbalancer port.

:::