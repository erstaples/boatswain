## Boatswain-Values example repo

This is the expected structure of the boatswain-values repo. Click on a folder to learn more about its purpose.

## Project Structure

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
