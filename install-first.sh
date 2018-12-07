#!/usr/bin/env bash

echo "Installing as if this is first time (including an update to dependencies)..."

go get github.com/inconshreveable/mousetrap
go get github.com/mitchellh/go-homedir
go get github.com/spf13/cobra
go get github.com/spf13/viper
go get github.com/spf13/pflag
go get github.com/cpuguy83/go-md2man/md2man
go get k8s.io/api/core/v1
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/client-go/kubernetes
go get k8s.io/client-go/kubernetes/typed/core/v1
go get k8s.io/client-go/tools/clientcmd

./install.sh