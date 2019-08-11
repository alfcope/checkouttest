NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY: dependency build unit-test integration-test docker-it-up docker-down clear docker-rmi docker-rmv

dependency:
	@go get -u github.com/DATA-DOG/godog/cmd/godog
	@go get -v ./...

build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	@docker build . -t local/payments-service

unit-tests:
	@go test -v -short ./...

integration-tests: docker-it-up
	@echo "$(OK_COLOR)==> Running ITs$(NO_COLOR)"
	@go test -v ./tests/integration/...; docker-compose -f ./tests/docker-compose-it.yml down

component-tests: docker-up
	@cd tests/component; godog; docker-compose down

docker-it-up:
	@docker-compose -f ./tests/docker-compose-it.yml up -d

docker-up:
	@docker-compose up -d
	@echo "$(WARN_COLOR)==> Waiting for services to be ready$(NO_COLOR)"
	@chmod +x waitForContainer.sh
	@./waitForContainer.sh

docker-down:
	@docker-compose down

clear: docker-down

swagger-doc:
	@swagger generate spec -m -o ./swagger.json

docker-rmi:
	@docker rmi $$(docker images -f "dangling=true" -q) -f

docker-rmv:
	@docker volume rm $$(docker volume ls -q -f dangling=true)

