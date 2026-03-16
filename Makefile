SRC=./src

build:
	make -C $(SRC) build

marabu:
	make -C $(SRC) marabu

no-bootstrap:
	make -C $(SRC) no-bootstrap

objex:
	make -C $(SRC) objex
run:
	make -C $(SRC) run

clean:
	$(RM) bin/*
	$(RM) logs/*

rebuild: clean build
