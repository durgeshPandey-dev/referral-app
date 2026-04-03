NX := npx nx

PROJECT ?= all
TARGET  ?= serve

help:
	node scripts/help.js

run:
ifeq ($(PROJECT),all)
	$(NX) run-many -t $(TARGET)
else
	$(NX) $(TARGET) $(PROJECT)
endif

affected:
	$(NX) affected -t $(TARGET)

generate:
	go mod tidy -modfile=go.tools.mod && go generate ./...
	npx kubb generate

ci:
	$(NX) run-many -t lint typecheck build test

.PHONY: help run affected generate ci
