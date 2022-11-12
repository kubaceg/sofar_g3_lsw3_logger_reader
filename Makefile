build-arm:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o sofar-arm

build:
	go build -o sofar