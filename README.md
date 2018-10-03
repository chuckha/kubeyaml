# Validate k8s documents

provide some yaml where the top level must have Kind and ApiVersion. The rest can be inferred


reading yaml and getting a list of keys and the associated values, a map[string]interface{}

Already know all vanilla k8s schemas

look up the schema with the kind/version, a string

(a schema, a map[string]interface{})

for key, value in map[string]interface{}

    if key not in schema: explode
    if value is not the expected type: explode
    if the type of key exists in the schema

## Deploying

1. Make the binary with `make kubeyaml`
2. Build the image with `./scripts/build-image.sh`
3. Push the image with `./scripts/push-image.sh`
4. On the server, restart the kubeyaml service which will pull and restart the image.

### Staging

1. Make the binary
2. `IMAGE_TAG=staging ./scripts/build-image.sh`
3. `IMAGE_TAG=staging ./scripts/push-image.sh`
4. `service kubeyaml-staging restart`

## Updating schemas

1. Run `go run scripts/update-schemas.go`

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
