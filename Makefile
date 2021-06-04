build-windows:
	set CGO_ENABLED=0 && set GOOS=windows && set GOARCH=amd64 && \
	go build -o client.exe
build-linux:
	set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=amd64 && \
	go build -o client.exe
router-generate:
	router-annotation --dir=/Users/linqiankai/go/src/ocpp-client/api --output=/Users/linqiankai/go/src/ocpp-client/init