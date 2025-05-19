ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

run:
	go run $(ROOT_DIR)/cmd/facilitator

test-e2e:
	go test -v $(ROOT_DIR)/test/e2e

