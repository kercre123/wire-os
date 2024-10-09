#!/bin/bash

if [[ -d ./wired ]]; then
	cd wired
fi

# make vector-gobot

cd vector-gobot
GCC="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang" \
GPP="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang++" \
make vector-gobot
cd ..

#make wired
# we don't need a docker container

CC="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang -w" \
CXX="${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/prebuilt/bin/arm-oe-linux-gnueabi-clang++" \
CGO_LDFLAGS="-L$(pwd)/vector-gobot/build -Wl,-rpath-link,${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/sysroot/lib -Wl,-rpath-link,${HOME}/.anki/vicos-sdk/dist/1.1.0-r04/sysroot/usr/lib -latomic" \
CGO_CFLAGS="-I$(pwd)/vector-gobot/include" \
CGO_ENABLED=1 \
GOARCH=arm \
GOARM=7 \
GOOS=linux \
go build  \
--trimpath \
-ldflags '-w -s' \
-o build/wired

upx build/wired
