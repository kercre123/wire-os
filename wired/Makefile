.PHONY: docker-builder wired

docker-builder:
	docker build -t armbuilder docker-builder/.

all: wired

go_deps:
	echo `go version` && cd $(shell pwd) && go mod download

wired:
	echo $(shell pwd)
	mkdir -p build
	docker container run  \
	-v "$(shell pwd)":/go/src/kercre123/wire-os/wired \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/kercre123/wire-os/wired \
	--user $(UID):$(GID) \
	armbuilder \
	go build  \
	-o build/wired \
	main.go

	docker container run \
		-v "$(shell pwd)":/go/src/kercre123/wire-os/wired \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/kercre123/wire-os/wired \
	--user $(UID):$(GID) \
	armbuilder \
	upx build/wired