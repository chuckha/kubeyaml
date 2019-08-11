# kubeyaml

## Webserver

### Docker

0. `docker run --network host registry.hub.docker.com/chuckdha/kubeyaml:latest`
0. Visit http://localhost:9000

### Manually

#### Requirements

* Go 1.11
* Go < 1.11 will be like any other go project without vendor or dependency management

1. `make kubeyaml-server`
0. `./kubeyaml-server`
0. Visit http://localhost:9000


## CLI

### Installing

#### Requirements

* Go 1.11
* Probably other versions of Go. Haven't tested.

1. `go get github.com/chuckha/kubeyaml/cmd/kubeyaml`

### Examples

#### Validate against recent versions of kubernetes

`cat test-yaml/deployment.yaml | kubeyaml`

#### Validate against one recent version of kubernetes

`cat test-yaml/deployment.yaml | kubeyaml -versions 1.12`

#### Validate against two recent versions of kubernetes

`cat test-yaml/pod.yaml | kubeyaml -versions 1.11,1.12`

#### Be quiet and rely on exit codes

```
`cat test-yaml/pod.yaml | kubeyaml -versions 1.12 -silent`
```


# Infra notes for kubeyaml.com

## Deploying

1. Make the binary with `make kubeyaml`
2. Build the image with `./scripts/build-image.sh`
3. Push the image with `./scripts/push-image.sh`
4. On the server, restart the kubeyaml service which will pull and restart the image.

### Staging

1. Make the binary
2. `IMAGE_TAG=staging ./scripts/build-image.sh`
3. `IMAGE_TAG=staging ./scripts/push-image.sh`
4. `service kube-validate-staging restart`

# Updating schemas

0. Add new minors to scripts/update-schemas.go
1. Run `go run scripts/update-schemas.go`
2. Update lookup.go to add new minor versions
3. Update cmd/generate/main.go
4. Update cmd/server/main.go

## TLS

use certbot on the host

```
sudo cerbot --nginx
```

and fill out the details


# Why didn't you just...

## Use a json schema validator

Most json schema validators do not validate YAML against json schemas.

## Generate go objects from the swagger spec using go-swagger

I did try this but then realized go can't do dynamic object lookup and only loads objects that are directly referenced.

## Write it in python

I did but then deleted it because I am not a good python programmer.
