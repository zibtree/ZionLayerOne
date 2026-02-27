.PHONY: build run devnet test lint clean docker

BINARY := ziond
CMD     := ./cmd/ziond
BUILD   := ./bin

build:
	@echo "🔨 Building $(BINARY)..."
	@mkdir -p $(BUILD)
	go build -o $(BUILD)/$(BINARY) $(CMD)
	@echo "✅ Built: $(BUILD)/$(BINARY)"

run: build
	$(BUILD)/$(BINARY) start

devnet: build
	@echo "🌐 Starting local devnet (3 validators)..."
	@mkdir -p ./data/validator{1,2,3}
	$(BUILD)/$(BINARY) start --rpc-port 8545 --validator 0xValidator1 --data-dir ./data/validator1 &
	$(BUILD)/$(BINARY) start --rpc-port 8546 --validator 0xValidator2 --data-dir ./data/validator2 &
	$(BUILD)/$(BINARY) start --rpc-port 8547 --validator 0xValidator3 --data-dir ./data/validator3 &
	@echo "✅ Devnet running on ports 8545, 8546, 8547"
	@echo "   RPC: http://localhost:8545"

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

clean:
	rm -rf $(BUILD) ./data

docker:
	docker build -t zionlayer:latest .

docker-devnet:
	docker-compose up --build

.DEFAULT_GOAL := build
