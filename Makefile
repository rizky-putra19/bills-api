.PHONY: help
help: # print all available make commands and their usages.
	@printf "\e[32musage: make [target]\n\n\e[0m"
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: setup
setup: # setup queue name and topic.
	@./scripts/create_topic.sh &> setup.log &
	@./scripts/create_queue.sh &> setup.log &

.PHONY: download-dependencies
download-dependencies: # download project dependencies
	go mod vendor -v

.PHONY: build-docker
build-docker: # create docker image
	docker build -t provider-id-linkqu .

.PHONY: build
build: # build golang application compatible with linux
	GOOS=linux go -o build/provider-id-linkqu cmd/main.go

.PHONY: run
run: # run application
	@./scripts/run.sh