build:
	go build -o main cmd/main/main.go
run: build
	./main