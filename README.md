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

1. Make the binary with `make kube-validate`
2. Build the image with `./scripts/build-image.sh`
3. Push the image with `./scripts/push-image.sh`
4. On the server, restart the kube-validate service which will pull and restart the image.

## TLS

use certbot on the host

```
sudo cerbot --nginx
```

and fill out the details
