BIN=../bin
MARABU=./cmd/marabu

build: standard-cli standard-headless no-bootstrap-cli no-bootstrap-headless test-obj test-obj-simple test-pset2

standard-cli:
	go build -tags "standard cli" -o $(BIN)/marabu-standard-cli $(MARABU)

standard-headless:
	go build -tags "standard headless" -o $(BIN)/marabu-standard-headless $(MARABU)

no-bootstrap-cli:
	go build -tags "no_bootstrap cli" -o $(BIN)/marabu-no-bootstrap-cli $(MARABU)

no-bootstrap-headless:
	go build -tags "no_bootstrap headless" -o $(BIN)/marabu-no-bootstrap-headless $(MARABU)

tests: test-obj test-obj-simple

test-obj-simple:
	go build -o $(BIN)/objectSimple ./cmd/testing/objectSimple

test-obj:
	go build -o $(BIN)/object ./cmd/testing/object

test-pset2:
	go build -o $(BIN)/pset2grader ./cmd/testing/pset2

run:
	go run $(MARABU)
