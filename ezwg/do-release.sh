#!/bin/sh
set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <tag>"
    echo "Release version required as argument"
    exit 1
fi

VERSION="$1"
GIT_COMMIT=$(git rev-list -1 HEAD)
BUILD_DATE=$(date)

RELEASE_FILE=RELEASE.md

LDFLAGS="-s -w \
    -X \"github.com/cirello-io/ezwg.GIT_COMMIT=$GIT_COMMIT\" \
    -X \"github.com/cirello-io/ezwg.VERSION=$VERSION\" \
    -X \"github.com/cirello-io/ezwg.BUILD_DATE=$BUILD_DATE\"\
"

# get release information
if ! test -f $RELEASE_FILE || head -n 1 $RELEASE_FILE | grep -vq $VERSION; then
    # file doesn't exist or is for old version, replace
    printf "$VERSION\n\n\n" > $RELEASE_FILE
fi

vim "+ normal G $" $RELEASE_FILE


# build
mkdir -p dist

export GOOS=linux
export CGO_ENABLED=0

GOARCH=arm GOARM=5 go build -ldflags="$LDFLAGS" cmd/ezwg.go
# upx -q ezwg
mv ezwg dist/ezwg-linux-arm5

GOARCH=amd64 go build -ldflags="$LDFLAGS" cmd/ezwg.go
# upx -q ezwg
mv ezwg dist/ezwg-linux-amd64

hub release create \
    --draft \
    -a dist/ezwg-linux-arm5#"ezwg linux-arm5" \
    -a dist/ezwg-linux-amd64#"ezwg linux-amd64" \
    -F $RELEASE_FILE \
    $1
