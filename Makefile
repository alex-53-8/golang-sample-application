BINARY=build/rest-application


test-unit:
	go test -v ./... -skip=TestIntegration/test/integration

test-integration:
	go test -v ./... -run=TestIntegration/test/integration

test-all:
	go test -v ./...

build:
	go build -o ${BINARY} .

run:
	export $$(cat .env | xargs) && go run .

clean:
	go clean
	rm -rf build
