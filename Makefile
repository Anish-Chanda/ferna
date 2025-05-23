# ==== Config ====
FRONTEND_DIR = app
BACKEND_DIR = backend
GO_APP_NAME = ferna-server
GO_BUILD_DIR = bin
GO_MAIN = main.go

# ==== Targets ====

# Common stuff
all: flutter-pub go-build flutter-build

test:
	@echo "Running Go tests..."
	cd $(BACKEND_DIR) && go test ./...
	@echo "Running Flutter tests..."
	cd $(FRONTEND_DIR) && flutter test

clean: flutter-clean go-clean

# Flutter (Frontend)
.PHONY: flutter-run flutter-clean flutter-build flutter-pub

flutter-run:
	cd $(FRONTEND_DIR) && flutter run

flutter-clean:
	cd $(FRONTEND_DIR) && flutter clean

flutter-build:
	cd $(FRONTEND_DIR) && flutter build apk
	cp $(FRONTEND_DIR)/build/app/outputs/flutter-apk/app-release.apk $(GO_BUILD_DIR)/ferna.apk

flutter-pub:
	cd $(FRONTEND_DIR) && flutter pub get

# Go (Backend)
.PHONY: go-run go-build go-lint go-fmt go-clean

go-run:
	cd $(BACKEND_DIR) && go run $(GO_MAIN)

go-build:
	mkdir -p $(GO_BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o ../$(GO_BUILD_DIR)/$(GO_APP_NAME) $(GO_MAIN)

go-lint:
	cd $(BACKEND_DIR) && go vet ./...

go-fmt:
	cd $(BACKEND_DIR) && go fmt ./...

go-clean:
	rm -rf $(GO_BUILD_DIR)

# Convenience
.PHONY: all clean
