
build: clean
	GOPATH=$(GOPATH):$(PWD) go build src/cmd/asnlookup.go
	@echo "Generated asnlookup binary in $(PWD)"

test: 
	GOPATH=$(GOPATH):$(PWD) go test ./...

clean:
	rm -rf asnlookup
