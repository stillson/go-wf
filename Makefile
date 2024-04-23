

build:
	go build -o wf

test:
	go test ./...

test-j:
	go test -json ./...

cover:
	go test -covermode count -coverprofile=coverage.data  ./...

report: cover
	go tool cover -html=coverage.data -o cover.html
	go tool cover -func=coverage.data -o cover.txt

check-format:
	test -z $$(go fmt ./...)

check: check-format static-check
	go vet ./...

#setup: install-go init-go install-lint
setup: install-lint copy-hooks

copy-hooks:
	chmod +x scripts/hooks/*
	cp -r scripts/hooks .git/.

install-lint:
	sudo curl -sSfL \
	https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh\
	| sh -s -- -b $$(go env GOPATH)/bin latest

static-check:
	$$(go env GOPATH)/bin/golangci-lint run

clean:
	go clean
	rm -f wf

GO_VERSION :=1.21.6

install-go:
	wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
	sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
	rm go$(GO_VERSION).linux-amd64.tar.gz

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.bashrc
