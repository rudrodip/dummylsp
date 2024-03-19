BIN=bin/dummylsp
MAIN=main.go

build:
	@echo "Building ${MAIN}"
	@go build -o $(BIN) $(MAIN)

run: build
	@echo "Running ${BIN}"
	@./$(BIN)

clean:
	@echo "Cleaning up..."
	@rm -f $(BIN)