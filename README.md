# kbrew

kbrew is to Kubernetes what Homebrew is to MacOS - a simple and easy to use package manager which hides the underlying complexity.

Let's talk in context of an example of at installing Kafka on a Kubernetes cluster
 - You need cert manager & Zookeeper & kube-prometheus-stack for monitoring installed
 - Zookeeper is a operator so you need to create a CR of Zookeeper cluster after installation of operator.
 - Then install Kafka operator
 - Create a CR of Kafka and wait for everything to stabilize.
 - Create ServieMonitor resources to enable prom scraping

With kbrew all of this happens with a single command (This command will change in near future):

```
$ kbrew install --config=./recipes/kafka-operator.yaml kafka-operator
```
## Helm chart or operator or Manifests - all abstracted

Kbrew abstracts the underlying chart or operator or manifest and gives you a recipe to install a stack with all basic configurations done.

## CLI Usage

For a quick demo, watch: https://www.youtube.com/watch?v=pWRZhZgfYSw 

### kbrew install

Installs a recipe in your cluster with all pre & posts steps and applications.

### kbrew search

TBD

### kbrew remove 

Partially implemented as of now.

## Terminology

### Recipe

A recipe defines the end to end installation of a set of things along with some custom steps and metadata. Checkout kafka-operator.yaml in context of following terms explained:

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