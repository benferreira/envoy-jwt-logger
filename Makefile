
.PHONY: build
build:
	tinygo build -o jwt-claim-logger.wasm -scheduler=none -target=wasi main.go


.PHONY: run
run:
	envoy -c envoy.yaml --log-format '%v'