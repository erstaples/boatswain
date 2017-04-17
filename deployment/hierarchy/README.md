# Local
## Bare minimum
### Run
create a `values.env.yaml` with `Neo4j.Persistence.Type` set to `hostPath` and `Neo4j.Persistence.Path` set to the hierarchy data directory. If no data directory exists, create it and copy data from the flash drive into this folder.

*values.env.yaml*
```yaml
Neo4j:
  Persistence:
    Type: hostPath
    Path: /Users/<user>/data/hierarchy

```

### Develop
Coming soon