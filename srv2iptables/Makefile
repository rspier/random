
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o srv2iptables -ldflags '-extldflags "-f no-PIC -static"' -tags 'osusergo netgo static_build' .

build-386:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o srv2iptables.386 -ldflags '-extldflags "-f no-PIC -static"' -tags 'osusergo netgo static_build' .
