#!/bin/sh

# compile a go program statically, so that it can run in any container,
# for example, in alpine containers
CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' $*
