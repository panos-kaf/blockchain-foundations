build:
	go build ./cmd/marabu

run:
	go run ./cmd/marabu

clean:
	$(RM) marabu
	$(RM) logs/*

rebuild: clean build
