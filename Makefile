.PHONY: build run devnet test lint fmt clean docker docker-devnet

BINARY := ziond
CMD    := ./cmd/ziond
BUILD  := ./bin

build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) $(CMD)
	@echo "Built: $(BUILD)/$(BINARY)"

run: build
	$(BUILD)/$(BINARY) start

devnet: build
	@echo "Starting local 3-node development environment..."
	@mkdir -p ./data/validator1 ./data/validator2 ./data/validator3
	$(BUILD)/$(BINARY) start --rpc-port 8545 --validator 0xValidator1 --data-dir ./data/validator1 & 	$(BUILD)/$(BINARY) start --rpc-port 8546 --validator 0xValidator2 --data-dir ./data/validator2 & 	$(BUILD)/$(BINARY) start --rpc-port 8547 --validator 0xValidator3 --data-dir ./data/validator3 &
	@echo "RPC endpoints: 8545, 8546, 8547"
	@echo "Note: these processes are currently independent proposer nodes, not a BFT cluster."

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

fmt:
	gofmt -w $$(find . -name "*.go" -not -path "./vendor/*")

clean:
	rm -rf $(BUILD) ./data

docker:
	docker build -t zionlayer:latest .

docker-devnet:
	docker compose up --build

.DEFAULT_GOAL := build
