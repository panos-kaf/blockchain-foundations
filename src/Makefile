BIN=../bin

build: marabu no-bootstrap objex

marabu:
	go build -o $(BIN)/marabu ./cmd/marabu

no-bootstrap:
	go build -o $(BIN)/no-bootstrap ./cmd/no-bootstrap

objex:
	go build -o $(BIN)/objectExchange ./cmd/testing/objectExchange

run:
	go run ./cmd/marabu
