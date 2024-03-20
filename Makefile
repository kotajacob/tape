# tape
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
GO ?= go
GOFLAGS ?=
RM ?= rm -f

all: tape

tape:
	$(GO) build $(GOFLAGS) -o build/tape

clean:
	$(RM) build/tape

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f build/tape $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/tape

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/tape

.DEFAULT_GOAL := all

.PHONY: all tape clean install uninstall
