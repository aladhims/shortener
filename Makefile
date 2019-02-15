.PHONY: setup
setup: ## generate go files from the defined protobuf files
	protoc --go_out=plugins=grpc:./pkg/shorten ./proto/shorten.proto
	protoc --go_out=plugins=grpc:./pkg/user ./proto/user.proto
	protoc --go_out=plugins=grpc:./pkg/notification ./proto/notification.proto


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'