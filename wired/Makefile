.PHONY: docker-builder vic-cloud vic-gateway

docker-builder:
	docker build -t armbuilder docker-builder/.

all: wired

go_deps:
	echo `go version` && cd $(shell pwd) && go mod download

wired: go_deps
	docker container run  \
	-v "$(shell pwd)":/go/src/kercre123/wire-os/wired \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-v /tmp:/tmp \
	-w /go/src/kercre123/wire-os/wired \
	--user $(UID):$(GID) \
	armbuilder \
	go build  \
	--trimpath \
	-ldflags '-w -s -linkmode internal -extldflags "-static"' \
	-o build/wired \
	main.go

	docker container run \
	-v "$(shell pwd)":/go/src/kercre123/wire-os/wired \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-v /tmp:/tmp \
	-w /go/src/kercre123/wire-os/wired \
	--user $(UID):$(GID) \
	armbuilder \
	upx build/wired
