SWAGGER_SPECFILE=docs/spec.json
BUILD_DIR=build/
LD_FLAGS='-s -w -extldflags "-static"'

test:
	go test -v -gcflags=-l ./...

vendor-module:
	go mod vendor

generate-swagger: vendor-module
	GO111MODULE=off CGO_ENABLED=0 swagger generate spec --scan-models -o ${SWAGGER_SPECFILE}
	GO111MODULE=off go get -u github.com/gobuffalo/packr/v2/packr2

pack-data:
	packr2

serve-swagger:
	swagger serve --flavor=swagger ${SWAGGER_SPECFILE}

download-tools:
	GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

build:
	go build -tags netgo -ldflags=${LD_FLAGS} -o ${BUILD_DIR}multimedia

.PHONY: build test