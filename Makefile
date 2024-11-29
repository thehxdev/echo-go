GO = go
GOFLAGS = -ldflags='-s -w -buildid=' -trimpath
BIN = echo-go

$(BIN):
	$(GO) build $(GOFLAGS) -o $(BIN) ./...

fmt:
	@$(GO) fmt ./...

clean:
	rm -rf $(BIN)
