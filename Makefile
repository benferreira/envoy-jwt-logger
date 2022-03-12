
.PHONY: build
build: clean
	tinygo build -o jwt-claim-logger.wasm -scheduler=none -target=wasi main.go

.PHONY: build-image
build-image: build
	docker build . -t localhost:32000/jwt-claim-logger-wasm:0.1.0

.PHONY: clean
clean: 
	go mod tidy
	rm jwt-claim-logger.wasm

.PHONY: run
run:
	envoy -c ./_envoy_config/envoy.yaml --concurrency 1 --log-format '%v'

.PHONY: test-integration
test-integration:
	curl localhost:18000 --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOiJzb21lYXVkaWVuY2UifQ.IaAz5CuyzFcDHbTOKWNzDgHd4xkmn-jgfmMn6IUGUFQ'