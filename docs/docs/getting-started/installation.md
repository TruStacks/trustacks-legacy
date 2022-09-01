---
title: Installation
sidebar_position: 1
slug: /installation
---

# Installation

:::info before you get started

If you get stuck or encounter any issues reach out to us on [discord](https://discord.gg/kGJBuHdv).

:::

## Docker

TruStacks requires a container runtime for certain features. 

### Windows Users

Native windows support is currently being tested. To get started now you must install a WSL2 based linux distro.

Follow [this guide](https://docs.microsoft.com/en-us/windows/wsl/install) to install the default Ubuntu based WSL distro on windows 10/11 (home or pro). Make sure to use the WSL 2 backend instead of Hyper-V when installing docker.

### Install

Navigate to the install guide for your operating system to install docker:

[linux](https://docs.docker.com/desktop/install/linux-install/)  
[osx](https://docs.docker.com/desktop/install/mac-install/)  
[windows](https://docs.docker.com/desktop/install/windows-install/)

## k3d

### Install

In order to use TruStacks we need a kubernetes cluster. Let's install [k3d](https://k3d.io) to get started on your local machine. 

K3d deploys a rancher k3s cluster into a docker container, so you will need to have docker installed. Follow the [installation guide](https://k3d.io/#installation) to get k3d installed on your local machine.

### Cluster Creation

Once k3d is installed, use the following command to create a cluster:

    k3d cluster create -p "8081:80@loadbalancer" trustacks

This command will create a new k3d cluster named `trustacks`. The `-p` option will create a loadbalancer on `8081`. This loadbalancer will be used later in the guide to access the toolchain components.

:::tip

If port `8081` is already in use on your machine then feel free to use a different port.

:::

Once the cluster is created check the output of `docker ps` and confirm that you have the `k3d-trustacks-serverlb` and `k3d-trustacks-server-0` containers.

    CONTAINER ID IMAGE                            COMMAND                  CREATED          STATUS          PORTS                                           NAMES
    3e0600614d9f ghcr.io/k3d-io/k3d-tools:5.4.3   "/app/k3d-tools noop"    27 seconds ago   Up 26 seconds                                                   k3d-trustacks-tools
    6b4aaee146ba ghcr.io/k3d-io/k3d-proxy:5.4.3   "/bin/sh -c nginx-pr…"   28 seconds ago   Up 19 seconds   0.0.0.0:8081->80/tcp, 0.0.0.0:44345->6443/tcp   k3d-trustacks-serverlb
    9eda9a6fc566 rancher/k3s:v1.23.6-k3s1         "/bin/k3s server --t…"   28 seconds ago   Up 24 seconds                                                   k3d-trustacks-server-0

Run the following command to confirm that the kubectl client was installed and that the cluster api is healthy:

    kubectl get ns

The command should return the following output:

    NAME              STATUS   AGE
    kube-system       Active   20s
    default           Active   20s
    kube-public       Active   20s
    kube-node-lease   Active   20s

## TruStacks Client

Get the latest TruStacks client for linux or mac from the [releases page](https://github.com/TruStacks/trustacks/releases).

Extract the `tsctl` binary and install it in your system path.

:::tip windows users

Download the `linux` distro and install the binary in your WSL2 distro.

:::