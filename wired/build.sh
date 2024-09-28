#!/bin/bash

if [[ -d ./wired ]]; then
	cd wired
fi

#make wired
# we don't need a docker container

CC="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang -w" \
CXX="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang++" \
GOARCH=arm \
GOARM=7 \
GOOS=linux \
go build  \
--trimpath \
-ldflags '-w -s -linkmode internal -extldflags "-static"' \
-o build/wired
