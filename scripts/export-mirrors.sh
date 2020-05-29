#!/bin/bash

export REPOS="qrng dynamolock runner bookmarkd oversight pglock pgqueue HumorChecker goherokuname snippetsd cci mlflappygopher kohrah-ani svc gists exp poo alfredemoji pidctl openapigen"

for repo in $REPOS
do
	./scripts/push-repo.sh $repo
done

