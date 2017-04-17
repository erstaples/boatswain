# Local
## Bare minimum
### Run
Since this app uses the MedBridge database, the only requirement to run the app is to deploy the MedBridge app to your cluster.
### Develop
Same as run step above, and in addition create a `values.env.yaml` with `deployment.hostPath` pointing to a local ace repository.

*values.env.yaml*
```yaml
deployment:
  hostPath: /Users/<name>/Programing/php/ace-tracker
```
