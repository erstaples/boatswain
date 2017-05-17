# boatswain

## Config file lives in ${HOME}/.boatswain.yaml

```yaml
release: /path/to/boatswain/deployment
builds:
- name: "medbridge"
  path: "deployment/build.sh"
  rootpath: "/Users/<name>/Programming/Php/Medbridge"
- name: "medflix"
  path: "deployment/build.sh"
  rootpath: "/Users/<name>/Programming/Php/Medflix/"
- name: "ace"
  path: "deployment/build.sh"
  rootpath: "/Users/<name>/Programming/Php/Ace"
```
