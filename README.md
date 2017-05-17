# boatswain

## Config file lives in `${HOME}/.boatswain.yaml`

## `boatswain stage`
In order to build and deploy automatically to staging, you need to add a `build` array in your config with each application you want to deploy to staging listed along with an absolute path to the project's root directory (`rootpath`) and the relative path, from the project's root, to the `build.sh` script (`path`). This is because `docker build` needs to be run from the project root dir.

```yaml
release: /path/to/boatswain/deployment
builds:
- name: medbridge
  path: deployment/build.sh
  rootpath: /Users/<name>/Programming/Php/Medbridge
- name: medflix
  path: deployment/build.sh
  rootpath: /Users/<name>/Programming/Php/Medflix/
- name: ace
  path: deployment/build.sh
  rootpath: /Users/<name>/Programming/Php/Ace
```
