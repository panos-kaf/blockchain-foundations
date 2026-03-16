build: marabu no-bootstrap objex

marabu:
	go build -o bin/marabu ./cmd/marabu

no-bootstrap:
	go build -o bin/no-bootstrap ./cmd/no-bootstrap

objex:
	go build -o bin/objectExchange ./cmd/testing/objectExchange
run:
	go run ./cmd/marabu

clean:
	$(RM) bins/*
	$(RM) logs/*

rebuild: clean build
