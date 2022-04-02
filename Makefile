
.PHONY: build build-image clean run test-integration

build: clean
	tinygo build -o jwt-claim-logger.wasm -scheduler=none -target=wasi main.go


build-image: build
	docker build . -t localhost:32000/jwt-claim-logger-wasm:0.1.0


clean: 
	go mod tidy
	rm -f jwt-claim-logger.wasm


run:
	envoy -c ./_envoy_config/envoy.yaml --concurrency 1 --log-format '%v'


test-integration:
	curl localhost:18000 --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJhdWQiOiJzb21lYXVkaWVuY2UifQ.IaAz5CuyzFcDHbTOKWNzDgHd4xkmn-jgfmMn6IUGUFQ'