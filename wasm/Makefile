WASM_PLUGINS := $(wildcard */.) 

docker.push:
	@for plugin in $(WASM_PLUGINS); do \
		echo "Pushing $${plugin}"; \
		pushd $${plugin}; \
		make docker.push || exit 1; \
		popd; \
	done
	@echo "Pushing all plugins"