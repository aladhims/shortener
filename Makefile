.PHONY: setup
setup: ## generate go files from the defined protobuf files
	protoc --go_out=plugins=grpc:./pkg/shorten ./proto/shorten.proto
	protoc --go_out=plugins=grpc:./pkg/user ./proto/user.proto
	protoc --go_out=plugins=grpc:./pkg/notification ./proto/notification.proto

.PHONY: start-compose
start-compose: ## start the application using docker-compose
	docker-compose -f deployments/docker-compose.yml  up -d user-db
	docker-compose -f deployments/docker-compose.yml  up

.PHONY: stop-compose
stop-compose: ## stop the application using docker-compose
	docker-compose -f deployments/docker-compose.yml down

.PHONY: stop-compose-rm
stop-compose-rm: ## stop the application using docker-compose and remove the images that have been built
	docker-compose -f deployments/docker-compose.yml down --rmi all

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'