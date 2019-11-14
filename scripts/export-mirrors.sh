#!/bin/bash

export REPOS="dynamolock runner bookmarkd oversight pglock pgqueue HumorChecker errors goherokuname snippetsd cci mlflappygopher kohrah-ani svc gists exp poo alfredemoji pidctl openapigen"

for repo in $REPOS
do
	./scripts/push-repo.sh $repo
done

