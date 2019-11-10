#!/bin/bash

export REPOS="HumorChecker dynamolock errors goherokuname runner bookmarkd snippetsd cci mlflappygopher kohrah-ani oversight pglock svc gists exp poo alfredemoji pidctl openapigen pgqueue"

for repo in $REPOS
do
	./scripts/push-repo.sh $repo
done

