.PHONY: all install clean

all: build

build:
	go build

install:
	mkdir -p /usr/local/gd-website-api
	cp DeepLX /usr/local/gd-website-api
	cp -r templates /usr/local/gd-website-api

clean:
	rm DeepLX
