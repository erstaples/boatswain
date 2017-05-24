# Hydra Helm Package
This is a fairly basic Helm package. It includes a templates folder, a Chart.yaml file, and a set of values files that can be used in different environments: 
* values.env.yaml: development
* values.prod.yaml: production
* values.yaml: sets the default values for all environments, which can be overridden in the environment-specific values files. For more information on developing Helm packages, [visit the Helm repo](https://github.com/kubernetes/helm)

## Deploying Hydra
Boatswain controls deployments and stagings. A few examples of deploying this package using boatswain: 


 Command | Function 
 --- | ---
 `boatswain release hydra` | Release Hydra to the default environment, development. The values.env.yaml file will be used
 `boatswain release hydra --environment production` | Release hydra to the staging environment. values.staging.yaml will be used
`boatswain stage push hydra <myPackageId>` | Builds and deploys Hydra to the staging environment and generates an ingress at the host defined in [.servicemaps](#todo-link-here) `Ingress.Template`: `mypackageid.staging.my-domain.com`