BUILD=build

all: clean build

build: go.sum
	@echo "Building ..."
	@go build -o $(BUILD)/go-infer-examples

go.sum: go.mod
	@echo "Ensure dependencies have not been modified"
	@GO111MODULE=on go mod verify

clean:
	@echo "Clean old built"
	@rm -rf $(BUILD)
	@go clean
