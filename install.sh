#!/usr/bin/env bash


if [ -z "$(git status --porcelain)" ]; then 
  export UNCOMMITTED=true
  # Working directory clean
else 
  export UNCOMMITTED=true
  # Uncommitted changes
fi

export GIT_COMMIT=$(git rev-list -1 HEAD) && \
  BRANCH=$(git branch | grep \* | cut -d ' ' -f2)
  NOW=$(date -u '+%Y-%m-%d_%I:%M:%S%p') && \
  go install -v -ldflags "-X github.com/doug4j/acdu/cmd.CommitHash=$GIT_COMMIT -X github.com/doug4j/acdu/cmd.BuildTime=$NOW -X github.com/doug4j/acdu/cmd.Branch=$BRANCH -X github.com/doug4j/acdu/cmd.HasUncommitted=$UNCOMMITTED"

echo CommitHash: $GIT_COMMIT
echo Done