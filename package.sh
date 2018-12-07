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

echo Built 6 platforms

echo Compressing binaries...

tar -czvf ./bin/linux_arm_acdu.tar.gz ./bin/linux_arm/acdu
tar -czvf ./bin/linux_386_acdu.tar.gz ./bin/linux_386/acdu
tar -czvf ./bin/linux_amd64_acdu.tar.gz ./bin/linux_amd64/acdu
tar -czvf ./bin/darwin_amd64_acdu.tar.gz ./bin/darwin_amd64/acdu
zip -r ./bin/windows_386_acdu.zip  ./bin/windows_386/acdu.exe
zip -r ./bin/windows_amd64_acdu.zip ./bin/windows_amd64/acdu.exe

echo Compressed binaries

echo Cleaning up...

rm ./bin/linux_arm/acdu
rm ./bin/linux_386/acdu
rm ./bin/linux_amd64/acdu
rm ./bin/darwin_amd64/acdu
rm  ./bin/windows_386/acdu.exe
rm ./bin/windows_amd64/acdu.exe

echo Cleaned up

echo Branch:      $BRANCH
echo BuildTime:   $NOW
echo Uncommitted: $UNCOMMITTED
echo CommitHash:  $GIT_COMMIT
echo Done