
.PHONY: deps
deps:
	#go get github.com/rakyll/statik
	go get github.com/go-bindata/go-bindata/...

.PHONY: build
build:
	(cd ui; yarn build)
	#rm -fr statik
	#$(shell go env GOPATH)/bin/statik -src=ui/public
	$(shell go env GOPATH)/bin/go-bindata -o ui.go -fs -prefix ui/public/ ui/public/

.PHONY: http
http:
	go run . http
