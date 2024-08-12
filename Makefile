.PHONY: all install clean

all: build

build:
	go build

install:
	mkdir -p /usr/local/gd-website-api/bin
	cp gd-website-api /usr/local/gd-website-api/bin
	cp config.ini /usr/local/gd-website-api
	cp -r templates /usr/local/gd-website-api

clean:
	rm gd-website-api
