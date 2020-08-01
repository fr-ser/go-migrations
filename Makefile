install:
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build:
	go build
	mv go-migrations db-migrations

unit-test:
	@echo Run as 'make unit-test args="-s -v"' to pass flags
	LOG_LEVEL=DEBUG gotest ./... -tags=unit ${args}

test: teardown
	docker-compose -f docker-compose.test.yaml up --detach
	@docker-compose -f docker-compose.test.yaml exec database timeout 5 sh -c 'until nc -z localhost 5432; do sleep 1; done'
	@docker-compose -f docker-compose.test.yaml exec database pg_isready --quiet
	@echo
	@echo Run as 'make test args="-count 1"' to pass flags
	@echo
	LOG_LEVEL=DEBUG gotest ./... ${args}

teardown:
	docker-compose -f docker-compose.test.yaml down --remove-orphans --timeout 1 --volumes

start-example:
	go run . start --dc-file ./example/docker-compose.yaml -p ./example/migrations -r

releases:
	env GOOS=linux GOARCH=amd64 go build
	mv go-migrations db-migrations-linux-amd64

	env GOOS=darwin GOARCH=amd64 go build
	mv go-migrations db-migrations-darwin-amd64