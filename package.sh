#!/usr/bin/env bash

#./install-first.sh

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
export ACDU_TEMP_BUILD_DIR=$(mktemp -d)
export ACDU_DIR=$(pwd)
echo "Branch:         $BRANCH"
echo "BuildTime:      $NOW"
echo "Uncommitted:    $UNCOMMITTED"
echo "CommitHash:     $GIT_COMMIT"
echo "Temp build Dir: $ACDU_TEMP_BUILD_DIR"
echo "ACDU dir:       $ACDU_DIR"

echo Packaging 3 platforms...
env GOOS=linux   GOARCH=amd64 go build -o $ACDU_TEMP_BUILD_DIR/bin/acdu_linux_amd64/acdu       -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
echo "Built linux amd64"
env GOOS=darwin  GOARCH=amd64 go build -o $ACDU_TEMP_BUILD_DIR/bin/acdu_darwin_amd64/acdu      -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
echo "Built darwin amd64"
env GOOS=windows GOARCH=amd64 go build -o $ACDU_TEMP_BUILD_DIR/bin/acdu_windows_amd64/acdu.exe -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"
echo "Built darwin amd64"
echo Built 3 platforms

echo Compressing binaries... 
#Linux 64 bit Intel/AMD
tar -czvf ./bin/acdu_linux_amd64.tar.gz  -C $ACDU_TEMP_BUILD_DIR/bin acdu_linux_amd64
#MacOS
tar -czvf ./bin/acdu_darwin_amd64.tar.gz -C $ACDU_TEMP_BUILD_DIR/bin acdu_darwin_amd64
#windows 64 bit
#rm ./bin/acdu_windows_amd64.zip
pushd $ACDU_TEMP_BUILD_DIR/bin/
zip -r $ACDU_DIR/bin/acdu_windows_amd64.zip ./acdu_windows_amd64
popd

echo Compressed binaries

echo Cleaning up...
rm -r $ACDU_TEMP_BUILD_DIR
echo Cleaned up

echo Done