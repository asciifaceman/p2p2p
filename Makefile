
test: ## Run tests if applicable
	go test github.com/asciifaceman/p2p2p/service --cover

clean: ## clean build dir
	rm -rf target

build-osx: ## Build for OSX
	@echo "Building target/p2p2p ..."
	@GOOS=darwin GOARCH=amd64 go build -o target/p2p2p
	@echo "Done."

build-linux: ## Build for linux
	@echo "Building target/p2p2p ..."
	@GOOS=linux GOARCH=amd64 go build -o target/p2p2p
	@echo "Done."

protos: ## Compiles protos
	@echo "Building protos ..."
	protoc -I service/ service/*.proto --go_out=plugins=grpc:service