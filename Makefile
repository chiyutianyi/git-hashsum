SHELL = /bin/sh

REVISION := $(shell git describe --tags --match 'v*' --always --dirty 2>/dev/null)
REVISIONDATE := $(shell git log -1 --pretty=format:'%ad' --date short 2>/dev/null)
PKG := github.com/chiyutianyi/git-hashsum
LDFLAGS = -s -w
ifneq ($(strip $(REVISION)),) # Use git clone
	LDFLAGS += -X $(PKG).revision=$(REVISION) \
		   -X $(PKG).revisionDate=$(REVISIONDATE)
endif

all: git-hashsum

git-hashsum:
	go build -ldflags="$(LDFLAGS)" -o bin/git-hashsum ./cmd
