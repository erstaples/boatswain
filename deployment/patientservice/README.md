# Local
## Bare minimum
### Run
create a `values.env.yaml` with `MySQL.Persistence.Path` set to the patient service data directory. If the directory doesn't exist, create it.

*values.env.yaml*
```yaml
MySQL:
  Persistence:
    Path: /Users/<user>/data/patient_service/

```

### Develop
Coming soon