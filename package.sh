#!/usr/bin/env bash

./install-first.sh

if [ -z "$(git status --porcelain)" ]; then 
  export UNCOMMITTED=false
  # Working directory clean
else 
  export UNCOMMITTED=true
  # Uncommitted changes
fi

export GIT_COMMIT=$(git rev-list -1 HEAD) 
export BRANCH=$(git branch | grep \* | cut -d ' ' -f2)
export NOW=$(date -u '+%Y-%m-%d_%I:%M:%S%p')

echo Packaging 6 platforms...

env GOOS=linux GOARCH=arm go build -o ./bin/linux_arm/acdu -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
env GOOS=linux GOARCH=386 go build -o ./bin/linux_386/acdu -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
env GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/acdu -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
env GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/acdu -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
env GOOS=windows GOARCH=386 go build -o ./bin/windows_386/acdu.exe -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
env GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/acdu.exe -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"

echo Branch:      $BRANCH
echo BuildTime:   $NOW
echo Uncommitted: $UNCOMMITTED
echo CommitHash:  $GIT_COMMIT
echo Done building 6 platforms