#!/bin/bash
ROOTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BINDIR=$ROOTDIR/bin

# export GOPATH
export GOPATH=$ROOTDIR

# update PATH
if [ -d "$BINDIR" ] && [[ ":$PATH:" != *":$BINDIR:"* ]]; then
    PATH="${PATH:+"$PATH:"}$BINDIR"
fi

# import missing libraries
go get github.com/google/uuid
go get github.com/pkg/errors
go get gopkg.in/yaml.v2
go get gopkg.in/abiosoft/ishell.v2
