set -e

<<<<<<< HEAD
GOAMD64=v4 go build -o builds/ecpdksap-ll-latest
=======
# go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest
GOARCH=arm64 go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest-arm64
# GOARCH=amd64 go build -ldflags "-s -w" -o builds/ecpdksap-ll-latest-amd64
>>>>>>> 2148bedb7d8057781bb079e4c09aa2b638954b28
