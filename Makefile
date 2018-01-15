
.PHONY: build
build:
	go build yaml2json.go
	go build json2yaml.go

.PHONY: deps
deps:
	go get "gopkg.in/yaml.v2"
	go get "github.com/ghodss/yaml"

.PHONY: clean
clean:
	- rm json2yaml
	- rm yaml2json
	-rm error.log
