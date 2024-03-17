.PHONY: build_linux builddir deploy

build_linux: builddir
	@GOOS=linux GOARCH=amd64 go build  -o ./build/trackrouter-linux-x86_64 ./cmd/telegram_server

builddir:
	@mkdir -p build
