install:
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build:
	go build

unit-test:
	LOG_LEVEL=DEBUG gotest ./... -tags=unit

test: teardown
	docker-compose -f docker-compose.test.yaml up --detach
	@docker-compose -f docker-compose.test.yaml exec database timeout 5 sh -c 'until nc -z localhost 5432; do sleep 1; done'
	@docker-compose -f docker-compose.test.yaml exec database pg_isready --quiet
	@echo
	LOG_LEVEL=DEBUG gotest ./...

teardown:
	docker-compose -f docker-compose.test.yaml down --remove-orphans --timeout 1 --volumes
