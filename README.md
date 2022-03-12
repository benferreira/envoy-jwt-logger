
# Envoy JWT Logger

WASM plugin for Envoy that logs JWT claims

## Requirements

* Golang - 1.17 or higher
* TinyGo
* Envoy
* Istio - 1.13 or higher for sidecar injection

## Build

```sh
make build
```

## Run

```sh
make run
```

## Istio deployment

This includes a sample Istio deployment in `istio/httpbin.yaml` that injects the WASM plugin in an `httpbin` app in the `default` namespace.
