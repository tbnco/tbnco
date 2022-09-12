# TBNCO

Network configuration operator.

## Development

If testing against a live cluster with an already running instance of the operator is required:

- Create test namespaces with the label key-value pair: `tbnco.github.io/environment: dev`
  - This excludes the namespaces from beeing processed by the deployed operator
  - Note: The operator still utilizes cluster-scoped resources
- Run local executable with the following [example configuration](./hack/config/local.yaml)
  - Replace `cacheNamespace` with your test namespace
  - `./bin/manager --config=<config.yaml>`

> Note, that updates to CRDs should be avoided when testing locally.
> Do not use `make install` or `make deploy` against a cluster with an already deployed operator!
> `make run` is ok, as it only executes the manager without further installation operations.
