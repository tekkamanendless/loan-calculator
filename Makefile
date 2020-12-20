.PHONY: clean
clean:
	rm -rf build

build:
	mkdir -p build

web: build
	cp cmd/web/assets/index.html build/
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" build/
	GOOS=js GOARCH=wasm go build -o build/web.wasm cmd/web/main.go

run: web
	cd build && python -m SimpleHTTPServer 8000
