build-windows:
	set CGO_ENABLED=0 && set GOOS=windows && set GOARCH=amd64 && \
	go build -o client.exe
build-linux:
	set CGO_ENABLED=0 && set GOOS=linux && set GOARCH=amd64 && \
	go build -o client.exe
