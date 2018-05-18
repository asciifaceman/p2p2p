
test: ## Run tests if applicable
  go test github.com/asciifaceman/p2p2p --cover

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
  #cd comms && protoc --go_out=plugins=grpc,import_path=comms:. *.proto
  protoc -I comms/ comms/*.proto --go_out=plugins=grpc:comms