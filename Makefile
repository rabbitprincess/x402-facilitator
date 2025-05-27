ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

run:
	go run $(ROOT_DIR)/cmd/facilitator \
		--config $(ROOT_DIR)/config.toml

test-e2e:
	go test -v $(ROOT_DIR)/test/e2e

generate-api:
	swag init -g api/server.go -o api/swagger --parseDependency

generate-abi:
	abigen --abi $(ROOT_DIR)/scheme/evm/eip3009/eip3009.abi \
		--pkg eip3009 \
		--out $(ROOT_DIR)/scheme/evm/eip3009/eip3009.go