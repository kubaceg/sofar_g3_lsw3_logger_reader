all: build-arm build-x86

build-arm:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o custom_components/sofar_g3_lsw3_logger_reader/sofar-arm

build-x86:
	go build -o custom_components/sofar_g3_lsw3_logger_reader/sofar-x86