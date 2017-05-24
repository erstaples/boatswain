# boatswain

## Getting Started
Boatswain works in conjunction with a boatswain-values project — you need both. A boatswain-values project is a set of Helm package folders. To learn more about Helm, or to learn more information about putting together a helm package, [visit the Helm repo](https://github.com/kubernetes/helm). The boatswain-values repo should be organized like this:
```
boatswain-values
|
├── deployment
|   |
|   ├── .cloudformation
|   |   └── cloudformation-template.yaml
|   |
|   ├── .globals
|   |   └── values.staging.yaml
|   |   └── values.yaml
|   |
|   ├── .servicemap
|   |   └── staging.yaml
|   |
|   ├── my-helm-package
|   |   |
|   |   ├── templates
|   |   |   └── ace-deployment.yaml
|   |   └── values.yaml
|   |   └── values.env.yaml
|   |   └── values.staging.yaml
|   |   └── values.prod.yaml
|   |
|   ├── another-helm-package
|   |   └── ...

```
To learn more about the boatswain-values repo, [take a look at the example project](https://github.com/MedBridge/boatswain/tree/add-example-values/examples/boatswain-values/deployment).

To learn more about the boatswain-values repo, [take a look at the example project]().

## Config file lives in `${HOME}/.boatswain.yaml`

## `boatswain stage push`
In order to build and deploy automatically to staging using `boatswain stage push`, you need to add a `builds` array in your config with an entry for each project you want to deploy to staging along with an absolute path to the project's root directory (`rootpath`) and the relative path — from the project's root — to the `build.sh` script (`path`). We need both values because `docker build` needs to be run from the project's root directory, but we also need to know where the `build.sh` script lives. See the example config below.

```yaml
ReleasePath: /path/to/boatswain-values/deployment
Builds:
- Name: medbridge
  Path: deployment/build.sh
  RootPath: /Users/<name>/Programming/Php/Medbridge
- Name: medflix
  Path: deployment/build.sh
  RootPath: /Users/<name>/Programming/Php/Medflix/
- Name: ace
  Path: deployment/build.sh
  RootPath: /Users/<name>/Programming/Php/Ace
```
