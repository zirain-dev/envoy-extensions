RM := rm -rf
WASM_PLUGIN_NAME := wasm-delay
DOCKER_IMAGE_HUB := ghcr.io/zirain
DOCKER_IMAGE_NAME := $(WASM_PLUGIN_NAME)
DOCKER_IMAGE_VER := latest
DOCKER_IMAGE := $(DOCKER_IMAGE_HUB)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_VER)

.PHONY : clean
clean:
	docker rmi $(DOCKER_IMAGE) || $(RM) ./$(WASM_PLUGIN_NAME).wasm
	rm -rf ./*.wasm
	

build: clean
	env GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o $(WASM_PLUGIN_NAME).wasm main.go

docker.build: build
	docker build -t $(DOCKER_IMAGE) . --build-arg WASM_PLUGIN_NAME=$(WASM_PLUGIN_NAME) --load

docker.push: docker.build
	docker push $(DOCKER_IMAGE)
	docker rmi $(DOCKER_IMAGE)
	docker pull $(DOCKER_IMAGE)