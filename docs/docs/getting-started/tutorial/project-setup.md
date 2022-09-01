---
title: Project Setup
slug: /tutorial/project-setup
sidebar_position: 1
---

# Setup

This tutorial will walk you through deploying a [react](https://reactjs.org/) javascript application.

## Goals

- Deploy a TruStacks toolchain
- Create a TruStacks application
- Run the application build and release

## React Tutorial Project

To get started, [fork](https://docs.github.com/en/get-started/quickstart/fork-a-repo) the [react-tutorial](https://github.com/TruStacks/react-tutorial.git) repository.

### About the project

`react-tutorial` is a [Create React App](https://create-react-app.dev/) derived react application. The project uses the [react](https://github.com/TruStacks/workflows/tree/main/workflows/react) **workflow** for build and release.

:::info What is a Workflow?

Workflows encompass vendor and industry best practices for building software in a specific framework. Using a workflow ensures standard, reliable, and repeatable delivery of your software.

:::

### SSH Access

Your fork of the `react-tutorial` project must be accessible over ssh. You can use an existing ssh key or [create a new ssh key](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent) if you don't already have one.

:::caution using deploy keys
 
 Deploy keys must have write access.

:::