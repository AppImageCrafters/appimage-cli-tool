get-deps:
	go get -v -t -d ./...
build:
	go build -o ./dist/app -v ./app
install:
	install ./app*.AppImage /usr/local/bin/app
