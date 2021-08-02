
[![CI](https://github.com/kbrew-dev/kbrew/actions/workflows/go.yml/badge.svg)](https://github.com/kbrew-dev/kbrew/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kbrew-dev/kbrew)](https://goreportcard.com/report/github.com/kbrew-dev/kbrew)
[![Release Version](https://img.shields.io/github/v/release/kbrew-dev/kbrew?label=kbrew)](https://github.com/kbrew-dev/kbrew/releases/latest)
[![License](https://img.shields.io/github/license/kbrew-dev/kbrew?color=light%20green&logo=github)](https://github.com/kbrew-dev/kbrew/blob/main/LICENSE)

![kbrew-logo](./images/kbrew-logo.png)




kbrew is to Kubernetes what Homebrew is to MacOS - a simple and easy to use package manager which hides the underlying complexity.

Let's take the example of installing Kafka on a Kubernetes cluster. If you are a developer trying this on a non prod environment, you want a quick and simple way to set it up. But a typical process looks like this:
 - You need cert-manager & Zookeeper & kube-prometheus-stack for monitoring installed
 - Zookeeper is an operator so you need to create a CR of Zookeeper cluster after installation of the operator.
 - Then install Kafka operator
 - Create a CR of Kafka and wait for everything to stabilize.
 - Create ServieMonitor resources to enable prom scraping

With kbrew, all of this happens with a single command (This command will change in near future):

```
$ kbrew install kafka-operator
```


## Helm chart or operator or Manifests - all abstracted

kbrew abstracts the underlying chart or operator or manifest and gives you a recipe to install a stack with all specific configurations done. You can trust that the recipe `just works`!

## Installation

### Install the pre-compiled binary

```bash
$ curl -sfL https://raw.githubusercontent.com/kbrew-dev/kbrew-release/main/install.sh | sh
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

```
$ kbrew --help
TODO: Long description

Usage:
  kbrew [command]

Available Commands:
  analytics   Manage analytics setting
  completion  Output shell completion code for the specified shell
  help        Help about any command
  info        Describe application
  install     Install application
  remove      Remove application
  search      Search application
  update      Update kbrew and recipe registries
  version     Print version information

Flags:
  -c, --config string       config file (default is $HOME/.kbrew.yaml)
      --config-dir string   config dir (default is $HOME/.kbrew)
      --debug               enable debug logs
  -h, --help                help for kbrew
  -n, --namespace string    namespace

Use "kbrew [command] --help" for more information about a command.
```

For a quick demo, watch: https://www.youtube.com/watch?v=pWRZhZgfYSw 

### Commonly used commands

#### kbrew search

Searches for a recipe for the given application. Lists all the available recipes if no application name is passed.

#### kbrew info

Prints applications details including registry and dependency information. 

#### kbrew install

Installs a recipe in your cluster with all pre & posts steps and applications.

#### kbrew update

Checks for kbrew updates and upgrades automatically if a newer version is available. Fetches updates for all the kbrew recipe registries

#### kbrew remove 

Uninstalls the application and its dependencies.


## Recipes

A kbrew recipe is a simple YAML file that declares the installation process of a Kubernetes app. It allows to *brew* Helm charts or vanilla Kubernetes manifests with scripts, also managing dependencies with other recipes.

Recipes can be grouped togther in structured directory called `Registry`. kbrew uses the [kbrew-registry](https://github.com/kbrew-dev/kbrew-registry/) by default. Any other resistry can be referred with the `--config-dir` flag.

### Recipe structure

A recipe looks like the below YAML
```
apiVersion: v1
kind: kbrew
app:
  repository:
    url: https://raw.githubusercontent.com/repo/manifest.yaml
    type: raw
  args:
    Deployment.nginx.spec.replicas: 4
  namespace: default
  version: v0.17.0
  pre_install:
  - apps:
      - OtherApp
  - steps:
      - echo "installing app"
  post_install:
  - steps:
      - echo "done installing"
  pre_cleanup:
  - steps
      - echo "deleting prerequisite"
  post_cleanup:
  - steps:
      - echo "app deleted"
```

`app` is the declaration of how a Kubernetes app - a Helm chart or a vanilla manifest - will get installed.

* `repository` : defines the source of the app
    - `url` : location of a Helm chart or a Kubernetes YAML manifest
    - `type`: can be `helm` or `raw`
* `args` : kbrew allows to modify the app via arguments that can modify the Helm chart values or manifest field values. See [Arguments](#Arguments) section for more
* `namespace` : Kubernetes namespace where the app should be installed. If not specified, `default` is considered
* `pre_install` : An app may require a few things to be present before installation. As in the case of kafka, it requires cert-manager and zookeeper. kbrew allows to specify this using `pre_install`. This can be a list of other recipe names or shell scripts to be run before the installation of the app
* `post_install` : After the finish of the app installation, it may require to do some operations or install other app. Citing the Kafka example again, a Kafka CR needs to be installed. This can be achieved by specifying  list of other recipe names or shell scripts under `post_install`.
* `pre_cleanup` and `post_cleanup` : While kbrew manages the order of removal of the dependencies and the app itself, there can be things that need manual removal. This typical could be resources manuaally installed with pre/post install steps, like a CR. The removal steps can be shell scripts specified with `pre_cleanup` and `post_cleanup`. The working section below well describes the order.


### Arguments

kbrew supports passing arguments to recipes as [Go templates](https://pkg.go.dev/text/template). All the functions from the [Sprig library](http://masterminds.github.io/sprig/) and the [lookup](https://helm.sh/docs/chart_template_guide/functions_and_pipelines/#using-the-lookup-function) & [include](https://helm.sh/docs/howto/charts_tips_and_tricks/#using-the-include-function) functions from Helm are supported.

**Helm app**: Arguments to a helm app can be the key-value pairs offerred by the chart in it's values.yaml file.
**Raw app**: These arguments patch the manifest of a raw app and can be specified in the format: `<Kind>.<Name>.<FieldPath>: <value>`. For example, to change `spec.replicas` of a `Deployment` named `nginx`, specify `Deployment.nginx.spec.replicas`



## Working of Kbrew Recipe

The process of how kbrew mangaes the installation of an app according to the recipe specification is depicted below. As can be seen, kbrew takes care of the order of pre/post actions.

![kbrew-install](./images/kbrew-install.png)


Similarly, while removing an app, kbrew takes care of the order of removal the dependent apps and the cleanup steps specified via `pre/post_cleanup` in the recipe.

![kbrew-install](./images/kbrew-remove.png)