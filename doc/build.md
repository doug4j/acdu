# How to build

Import dependencies

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

go install -ldflags \
    "-X cmd.GitCommit=57b9870 -X cmd.GitBranch=dry-run -X main.GitState=dirty -X main.Version=0.1.0 -X main.BuildDate=2016-08-08T20:50:21Z"

go build -ldflags "-X github.com/doug4j/acdu/cmd.Version=0.9.0-SNAPSHOT github.com/doug4j/acdu/cmd.CommitHash=abc"