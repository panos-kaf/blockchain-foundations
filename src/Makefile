BIN=../bin
MARABU=./cmd/marabu

build: standard-cli standard-headless no-bootstrap-cli no-bootstrap-headless objex

standard-cli:
	go build -tags "standard cli" -o $(BIN)/marabu-standard-cli $(MARABU)

standard-headless:
	go build -tags "standard headless" -o $(BIN)/marabu-standard-headless $(MARABU)

no-bootstrap-cli:
	go build -tags "no_bootstrap cli" -o $(BIN)/marabu-no-bootstrap-cli $(MARABU)

no-bootstrap-headless:
	go build -tags "no_bootstrap headless" -o $(BIN)/marabu-no-bootstrap-headless $(MARABU)

objex:
	go build -o $(BIN)/objectExchange ./cmd/testing/objectExchange

run:
	go run $(MARABU)
