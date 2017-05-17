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
  rootpath: "/Users/<name>/Programming/Php/CourseRecommender/var/www"
- name: "ace"
  path: "deployment/build.sh"
  rootpath: "/Users/<name>/Programming/Php/ace-tracker"
```
