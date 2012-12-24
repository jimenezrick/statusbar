PREFIX=/usr/local

.PHONY: all install clean

all:
	go build

install: all
	install dwmstatus $(DESTDIR)$(PREFIX)/bin

clean:
	go clean
