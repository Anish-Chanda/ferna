SHELL := bash

# ==== Config ====
BACKEND_DIR = backend
GO_APP_NAME = ferna-api
GO_BUILD_DIR = bin

# ==== Backend Commands ====
.PHONY: build-api run-api

build-api:
	mkdir -p $(GO_BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(GO_BUILD_DIR)/$(GO_APP_NAME) .

run-api:
	@bash -c "set -a; . .env.example; set +a; cd $(BACKEND_DIR) && go run ."
