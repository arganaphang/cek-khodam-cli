set dotenv-load

# build -> build application
build:
	go build -o main

# run -> application
run:
	./main

# dev -> run build then run it
dev: 
	watchexec -r -c -e go -- just build run
