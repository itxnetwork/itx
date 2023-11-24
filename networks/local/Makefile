ARCH := $(shell uname -m)

all:
	docker build --build-arg ARCH=$(ARCH) --no-cache --tag itx/node ../.. -f itx/Dockerfile

.PHONY: all
