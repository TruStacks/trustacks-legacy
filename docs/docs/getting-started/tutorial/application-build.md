---
title: Application Build & Release
slug: /tutorial/application build
sidebar_position: 4
---

## Running the application build

During the application creation, the concourse ci driver configured a pipeline for us to build the react-tutorial application. 

Next we will run the application build in concourse.

First, navigate to [concourse](http://concourse.local.gd:8081) and confirm that the react-tutorial pipeline was created.

<img src="/img/concourse-pipelines.jpg" />

You should see the `react-tutorial` application in the browser.

Next, click the blank gray square in the pipeline widget to navigate to the pipeline builds page.

## Run The Build

Click the "+" plus button at the top right to run the application build.

<img src="/img/concourse-builds.jpg" />

The build will start and you will see a series of resources in the browser.

<img src="/img/concourse-build-started.jpg" />

Once the `task: build` widget turns yellow, click on it to expand the build output

## Application Deployment

Once the build turns green, signifying a successfuly pipeline run, the application is up and running.

Navigate to the application in [Argo CD](http://argo-cd.local.gd:8081/applications/react-tutorial?view=tree&resource=) to view the deployed application assets.

<img src="/img/argocd-application.jpg" />

## View the Application

The tutorial uses an "in-cluster" deployment of the react application.

Run the following command to create a kubectl port forwarding proxy

    kubectl port-forward svc/staging-react-tutorial -n react-tutorial 54321:8080

Navigate to the [proxy address](http://localhost:54321) to view the deployed react application.

<img src="/img/react-landing-page.jpg" />
