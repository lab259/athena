GOPATH=$(CURDIR)/../../../../
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

EXAMPLES=$(shell ls ./examples/)

VERSION := `git describe --exact-match --tags 2> /dev/null || git rev-parse HEAD`
LDFLAGS=-X=main.version=$(VERSION)

build:
	@$(GOPATHCMD) go build "-ldflags=$(LDFLAGS) -s -w" -o ./bin/athena -v ./athena
	@test -d ./examples && $(foreach example,$(EXAMPLES),$(GOPATHCMD) go build "-ldflags=$(LDFLAGS)" -o ./bin/$(example) -v ./examples/$(example) &&) :

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
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html
	@xdg-open .cover/report.html 2> /dev/null > /dev/null

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-update:
	@$(GOPATHCMD) dep update -v $(PACKAGE)

vet:
	@$(GOPATHCMD) go vet ./...

fmt:
	@$(GOPATHCMD) gofmt -e -s -d *.go

.PHONY: deps deps-ci coverage coverage-ci test test-watch coverage coverage-html