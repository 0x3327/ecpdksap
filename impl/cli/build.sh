set -e

# go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest
GOARCH=arm64 go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest-arm64
# GOARCH=amd64 go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest-amd64