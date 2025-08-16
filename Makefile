WASM_DIR = assets/wasm
WASM_PROJECTS = wasm-projects/cfd-but-wasm

.PHONY: dev build-wasm watch-wasm serve-jekyll
dev:
	$(MAKE) -j2 watch-wasm serve-jekyll

build-wasm:
	@for project in $(WASM_PROJECTS); do \
		$(MAKE) -C $$project build TARGET_DIR=../../$(WASM_DIR) PROJECT_NAME=$$(basename $$project); \
	done

watch-wasm:
	@for project in $(WASM_PROJECTS); do \
		$(MAKE) -C $$project watch TARGET_DIR=../../$(WASM_DIR) PROJECT_NAME=$$(basename $$project) & \
	done; wait

serve-jekyll:
	bundle exec jekyll serve --livereload

clean:
	rm -rf $(WASM_DIR)/*
