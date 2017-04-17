# Getting started
Follow the instructions located in [Bare minimum](#bare-minimum). Look at [Run](#run) if all you want to do is get an instance of the app running in a local cluster with minimal configuration. [Develop](#develop) will allow you to develop the app in a local cluster. 

Look at [Image tags](#image-tags) for a quick demo of how image tags work in this deployment, how to version an image, and how to version all images.

[Docker images](#docker-images) gives a breakdown of each image in this deployment as well as relevant configuration options you can set in `values.yaml`.

# Local deployment
## Bare minimum
### Run
You will need to create a data directory for the MedBridge database to run locally. Create a `values.env.yaml` file with the following:
* `MedBridge.Persistence.MySQL.HostPath` to a local data directory. If you don't have a data directory, create an empty one.
* database credentials for your local db

*values.env.yaml*
```yaml
MedBridge:
  Persistence:
    MySQL:
      HostPath: /Users/<user>/data/medbridge
  Secrets:
    MedbridgeEdDb:
      User: root
      Password: test123123
```

### Develop
In addition to the run values.env.yaml file above, add `MedBridge.Persistence.Source.Hostpath` pointing to the MedBridge repo:

*values.env.yaml*
```yaml
MedBridge:
  Persistence:
    Source:
      HostPath: /Users/<user>/Programing/php/Medbridge
    MySQL:
      HostPath: /Users/<user>/data/medbridge
  Secrets:
    MedbridgeEdDb:
      User: root
      Password: test123123
```

# Image tags
You can set a single `ImageTag` for all containers, or selectively update `ImageTag`s for specific containers.

```yaml
MedBridge:
  ImageTag: "abc123"
  Images:
    Source:
      ImageTag: "xyz789"
     Migrations:
       ImageTag:
     Nginx:
       ImageTag:
     PHPCLI:
       ImageTag:
     PHP:
       ImageTag: 
     MySQL:
       ImageTag: "5.7"
```
Using the `values.yaml` snippet above, Migrations, Nginx, PHPCLI, and PHP will be assigned the `abc123` image tag, 
MySQL the `5.7` image tag, and Source will be assigned `xyz789`. This will often be the case: source code is going to change far more often than the extensions, apt-get packages, and other software that make up the other images.

# Docker Images
We use several images in this deployment. 
## Source
Contains the source code. Any updates to the code base (with exception to `httpdocs/images`) will require a rebuild of `Source`. Unless we're making infrastructure changes or adding on new extensions, this and `ImageAssets` are the only images that needs to be updated on a continual basis. It is used in an init container, which will mount the code contained in this image into a volume shared by `Nginx`, `PHP-FPM`, `PHP-CLI`, etc.

You may override the init-containers and instead mount the source code:
```yaml
MedBridge:
  Persistence:
    Source:
      HostPath: /Users/<user>/Programing/php/Medbridge
```
This will skip the init-container volume mount step and mount your local source code directory onto the container.

## ImageAssets
Contains `httpdocs/images` directory. Any additions or updates to files in this directory will require a rebuild of `ImageAssets`. Unless we're making infrastructure changes or adding on new extensions, this and `Source` are the only images that needs to be updated on a continual basis.

Like `Source`, you may override the init-containers and instead mount the source code:
```yaml
MedBridge:
  Persistence:
    Source:
      HostPath: /Users/<user>/Programing/php/Medbridge
```

## Filebeat
Handles routing logs to ELK stack. The image tag should not change often, unless we're upgrading our ELK stack. You can toggle logging in a `values.yaml` file:
```yaml
MedBridge:
  Logging: true
```
Leave this property blank and the pod will not spin up a filebeat container, and no logs will be sent to the ELK stack.

## Nginx
Contains Nginx web server. You can configure `sites-enabled` in the `values.yaml` file. This example enables the localmed.com configuration:
```yaml
MedBridge:
  Configs:
    ConfigMaps:
      Nginx:
        SitesEnabled:
          - localmed.com.conf
```
Other configurations are stored in the `sites-available` configmap.

## PHP-CLI
CLI utilities such as PHP unit tests and application CLIs

## NewRelic Daemon
Handles dispatching data to NewRelic. You can switch NewRelic reporting on by including `MedBridge.NewRelic.AppName`:

```yaml
MedBridge:
  NewRelic:
    AppName: medbridge.io
```
Or switch NewRelic reporting off by setting it to a falsy value:
```yaml
MedBridge:
  NewRelic:
    AppName:
```

## PHP-FPM
Contains all extensions, etc to run the site.
