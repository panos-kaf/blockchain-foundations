SRC=./src

build:
	make -C $(SRC) build

standard-cli:
	make -C $(SRC) standard-cli

standard-headless:
	make -C $(SRC) standard-headless

no-bootstrap-headless:
	make -C $(SRC) no-bootstrap-headless

no-bootstrap-cli:
	make -C $(SRC) no-bootstrap-cli

objex:
	make -C $(SRC) objex
run:
	make -C $(SRC) run

clean:
	$(RM) bin/*
	$(RM) logs/*

rebuild: clean build
