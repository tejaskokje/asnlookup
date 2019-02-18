
build: clean
	GOPATH=$(GOPATH):$(PWD) go build src/cmd/asnlookup.go
	@echo "Generated asnlookup binary in $(PWD)"

test: 
	go test ./...

clean:
	rm -rf asnlookup
