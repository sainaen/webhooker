SOURCE = $(wildcard *.go)
TAG ?= $(shell git describe --tags)
GOBUILD = go build -ldflags '-s -w'

ALL = $(foreach suffix,win.exe linux osx,\
		build/webhooker-$(suffix))

all: $(ALL)

clean:
	rm -f $(ALL)

test:
	go test
	cram tests/cram.t

win.exe = windows
osx = darwin
build/webhooker-%: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=amd64 $(GOBUILD) -o $@

release: $(ALL)
ifndef desc
	@echo "Run it as 'make release desc=tralala'"
else
	github-release release -u piranha -r webhooker -t "$(TAG)" -n "$(TAG)" --description "$(desc)"
	@for x in $(ALL); do \
		echo "Uploading $$x" && \
		github-release upload -u piranha \
                              -r webhooker \
                              -t $(TAG) \
                              -f "$$x" \
                              -n "$$(basename $$x)"; \
	done
endif
