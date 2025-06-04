.PHONY : all clean build docker.build docker.push

RM := rm -rf
WASM_PLUGIN_NAME := filter-state
DOCKER_IMAGE_HUB := ghcr.io/zirain
DOCKER_IMAGE_NAME := filter-state
DOCKER_IMAGE := $(DOCKER_IMAGE_HUB)/$(DOCKER_IMAGE_NAME)
VERSION_LIST := 0.0.1

all: docker.push

clean:
	$(RM) *.wasm
	
build: clean
	$(foreach VERSION, $(VERSION_LIST), GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -ldflags '-X main.injectionVersion=$(VERSION)' -o $(WASM_PLUGIN_NAME)-$(VERSION).wasm plugin.go;)

docker.build: build
	$(foreach VERSION, $(VERSION_LIST), docker buildx build . -t $(DOCKER_IMAGE):$(VERSION) --build-arg WASM_BINARY=$(WASM_PLUGIN_NAME)-$(VERSION).wasm --load;)

docker.push: docker.build
	$(foreach VERSION, $(VERSION_LIST), docker push $(DOCKER_IMAGE):$(VERSION);)