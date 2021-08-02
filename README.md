# kbrew

[![CI](https://github.com/kbrew-dev/kbrew/actions/workflows/go.yml/badge.svg)](https://github.com/kbrew-dev/kbrew/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kbrew-dev/kbrew)](https://goreportcard.com/report/github.com/kbrew-dev/kbrew)
[![Release Version](https://img.shields.io/github/v/release/kbrew-dev/kbrew?label=kbrew)](https://github.com/kbrew-dev/kbrew/releases/latest)
[![License](https://img.shields.io/github/license/kbrew-dev/kbrew?color=light%20green&logo=github)](https://github.com/kbrew-dev/kbrew/blob/main/LICENSE)

kbrew is to Kubernetes what Homebrew is to MacOS - a simple and easy to use package manager which hides the underlying complexity.

Let's talk in context of an example of at installing Kafka on a Kubernetes cluster
 - You need cert manager & Zookeeper & kube-prometheus-stack for monitoring installed
 - Zookeeper is a operator so you need to create a CR of Zookeeper cluster after installation of operator.
 - Then install Kafka operator
 - Create a CR of Kafka and wait for everything to stabilize.
 - Create ServieMonitor resources to enable prom scraping

With kbrew all of this happens with a single command (This command will change in near future):

```
$ kbrew install kafka-operator
```
## Helm chart or operator or Manifests - all abstracted

Kbrew abstracts the underlying chart or operator or manifest and gives you a recipe to install a stack with all basic configurations done.

## Installation

### Install the pre-compiled binary

```bash
$ curl -sfL https://raw.githubusercontent.com/kbrew-dev/kbrew/main/install.sh | sh
```

### Compiling from source

#### Step 1: Clone the repo

```bash
$ git clone https://github.com/kbrew-dev/kbrew.git
```

#### Step 2: Build binary using make

```bash
$ make
```

## CLI Usage

For a quick demo, watch: https://www.youtube.com/watch?v=pWRZhZgfYSw 

### kbrew install

Installs a recipe in your cluster with all pre & posts steps and applications.

### kbrew search

Searches for a recipe for the given application. Lists all the available recipes if no application name is passed.

### kbrew update

Checks for kbrew updates and upgrades automatically if a newer version is available.
Fetches updates for all the kbrew recipe registries

### kbrew remove 

Partially implemented as of now.

## Terminology

### Recipe

A recipe defines the end to end installation of a set of things along with some custom steps and metadata. 
Checkout kafka-operator.yaml or similar recipes in context of the terms being explained here. 
All public recipes are maintained in the repo https://github.com/kbrew-dev/kbrew-registry 

#### Repository

A repository has a name, type and a URL and based on type, there are two ways you can define a repo:

#### Helm

You use the base URL of the Chart repo and name of chart to be used along with the type as `helm`

```
  repository:
    name: banzaicloud-stable
    url: https://kubernetes-charts.banzaicloud.com
    type: helm
```
#### Raw

You can use a URL to a operator or a RAW yaml file with type `raw`

```
app:
  repository:
    name: cert-manager
    url: https://github.com/jetstack/cert-manager/releases/download/v0.10.1/cert-manager.yaml
    type: raw
```

#### Pre_Install & Post_Install

Every recipe has a bunch of pre & post install activities. There are two types supported today:

##### Apps

Points to other formulas in repo

##### Steps

Custom step - which can be inline commands or scripts etc.
