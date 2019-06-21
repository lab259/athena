GOPATH=$(shell readlink -f $(CURDIR)/../../../../)
GOPATHCMD=PROJECT_ROOT=$(CURDIR) GOPATH=$(GOPATH)
GOCMD=$(GOPATHCMD) go

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

CMDS=$(shell test -d ./cmd/ && ls ./cmd/)
EXAMPLES=$(shell test -d ./examples/ && ls ./examples/)

VERSION := `git describe --exact-match --tags 2> /dev/null || git rev-parse HEAD`
LDFLAGS=-X=main.version=$(VERSION)

build:
	@$(foreach cmd,$(CMDS),$(GOCMD) build "-ldflags=$(LDFLAGS) -s -w" -o ./bin/$(cmd) -v ./cmd/$(cmd) &&) :
	@$(foreach example,$(EXAMPLES),$(GOCMD) build "-ldflags=$(LDFLAGS) -s -w" -o ./bin/$(example) -v ./examples/$(example) &&) :

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find ./* -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"
	@cp "${COVERAGEFILE}" coverage.txt

coverage-html:
	@$(GOCMD) tool cover -html="${COVERAGEFILE}" -o .cover/report.html
	@xdg-open .cover/report.html 2> /dev/null > /dev/null

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-update:
	@$(GOPATHCMD) dep ensure -update -v $(PACKAGE)

vet:
	@$(GOCMD) vet ./...

fmt:
	@$(GOCMD) fmt ./...

get:
	@$(GOCMD) get -u $(PACKAGE)

generate:
	@$(GOCMD) generate ./...

.PHONY: build test test-watch coverage coverage-ci coverage-html dep-ensure dep-update vet fmt get generate