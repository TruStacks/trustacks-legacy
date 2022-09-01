---
title: Application Creation
slug: /tutorial/create-application
sidebar_position: 3
---

# Applications

TruStacks applications deploy CI/CD assets that consume the software toolchain.

## Container Registry

The react worfklow for the react-tutorial project requires a container registry.

If you already have a docker registry then you can skip to the [Configuration](#configuration) section.

### Docker Hub

If you do not have a container registry you can use [docker hub](https://hub.docker.com/).

After creating your account follow [this guide](https://docs.docker.com/docker-hub/access-tokens/) to create an access token for use in the next step.

#### Repository host

Certain container registries, such as Azure Container Registry, allow dynamic creation of container repositories. 

Docker Hub and other container registries requires that repositories be created before pushing images.

The format for container images build by the TruStacks CI workflow is:

    <registry-host>/<project-name>:<image-tag>

If you are using Docker Hub or a container that requires the repository to exist before pushing images then create the repository now.

## <a name="configuration"></a> Configuration

To get started, add the following to your `react-tutorial-config.yaml` configuration file:

```yaml
applications:
- name: react-tutorial
  source: https://github.com/trustacks/workflows
  workflow: react
  vars:
    image: "quay.io/trustacks/workflows"
    imageTag: "1.0.0"
    workflow: "react"
    project: "react-tutorial"
    gitRemote: "<your fork's ssh url>"
    registryHost: "<your registry hostname>"
    registryUsername: "<your registry username>"
  secrets:
    gitPrivateKey: |- 
      "<your ssh or deploy key>"
    registryPassword: "<your container registry password or access token>"
```

Let's break down the configuration values:

> `applications` is the list of available applications to create.  
> `applications[*].name` is the name of the application.  
> `applications[*].source` is the the git repository that contains the workflow resources.  
> `applications[*].workflow` is the workflow to use from the workflow source.  
> `applications[*].vars` are plaintext values that are used by the workflow CI/CD build.  
> `applications[*].secrets` are secret values that are used by the workflow CI/CD build.

:::tip

Remember to append your project name to `registryHost` if your container registry does not support dynamically creating container repositories.

*ie. `registry.hub.docker.com/<project-name>`*

:::

:::caution

Use string quotes to ensure `vars` and `secrets` are interpreted as strings. Numerical values will result in errors.

:::

:::caution

Your repository url must use ssh, and you must use an ssh key.

:::

### Workflow Inputs

Inputs are provided throught **vars** and **secrets**. the input configuration values are passed directly to the CI/CD build.

### Creating the application

After adding the application configuration to `react-tutorial-config.yaml`, we are now ready to create the application.

Use the following command create the application: 

```bash
tsctl application create --name react-tutorial --config react-tutorial-config.yaml
```

## Application components

In addition to the [Core Components](/tutorial/install-toolchain#-core-components-at-a-glance) installed in the toolchain, the [react](https://github.com/TruStacks/workflows/tree/main/workflows/react) workflow installs [Argo CD](https://argoproj.github.io/cd/) as a supplemental component.

Argo CD will be used to deploy the application

### Day 1 Automation

In addition to the user provided workflow inputs, toolchain components can provide system inputs.

Components built for the TruStacks eco-system implement Day 1 automation to make consumption seamless.

The Argo CD [component](https://github.com/TruStacks/catalog/blob/main/pkg/components/argocd/component.go) creates a service account, confingures rbac, and exposes its server endpoint and service account secret during installation.

Activites such as creating SSO client secrets, configuring rbac groups, and all other activities related to the controlled component and "zero touch" consumption readiness are completed during installation. 

:::info

Day 1 Automation is covered in greater detail in the Core Concepts.

:::