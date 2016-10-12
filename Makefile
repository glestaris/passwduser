.PHONY: help \
	test lint

test:
	CGO_ENABLED=0 go test -v .

help:
	@echo '    test ................................ runs tests'
	@echo '    lint ................................ lint the Go code'
	@echo '    help ................................ print this message'

lint:
	gometalinter --disable=gotype --disable=dupl --disable=gas \
		--deadline 100s \
		--tests
