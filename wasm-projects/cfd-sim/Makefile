.PHONY: build dev

build:
	mkdir -p ../../assets/wasm/cfd-sim
	tinygo build -o ../../assets/wasm/cfd-sim/main.wasm -target wasm .
	cp $$(tinygo env TINYGOROOT)/targets/wasm_exec.js ../../assets/wasm/cfd-sim/

dev:
	find . -name "*.go" | entr -r make build
