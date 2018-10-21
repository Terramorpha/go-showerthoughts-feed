win64:
	GOOS="windows" GOARCH="amd64" go build -o winbin
install:
	go build
	sudo cp -f redditFeed /bin/trailingShowerThought
