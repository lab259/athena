COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

CMDS=$(shell test -d ./cmd/ && ls ./cmd/)
EXAMPLES=$(shell test -d ./examples/ && ls ./examples/)

VERSION=$(shell git describe --exact-match --tags 2> /dev/null || git rev-parse HEAD)
LDFLAGS=-X=main.version=$(VERSION)

build:
	@$(foreach cmd,$(CMDS),go build "-ldflags=$(LDFLAGS) -s -w" -o ./bin/$(cmd) -v ./cmd/$(cmd) &&) :
	@$(foreach example,$(EXAMPLES),go build "-ldflags=$(LDFLAGS) -s -w" -o ./bin/$(example) -v ./examples/$(example) &&) :

test:
	@go run github.com/onsi/ginkgo/ginkgo --failFast ./...

test-watch:
	@go run github.com/onsi/ginkgo/ginkgo watch -cover -r ./...

coverage-ci:
	@mkdir -p $(COVERDIR)
	@go run github.com/onsi/ginkgo/ginkgo -r -covermode=count --cover --randomizeAllSpecs --randomizeSuites --failOnPending --trace --compilers=2
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name '*.coverprofile' -exec cat {} \; -exec rm -f {} \; | grep -h -v "^mode:" >> ${COVERAGEFILE}

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"
	@cp "${COVERAGEFILE}" coverage.txt

coverage-html:
	@go tool cover -html="${COVERAGEFILE}" -o .cover/report.html
	@xdg-open .cover/report.html 2> /dev/null > /dev/null

vet:
	@go vet ./...

fmt:
	@go fmt ./...

generate:
	@go generate ./...

.PHONY: build test test-watch coverage coverage-ci coverage-html vet fmt get generate