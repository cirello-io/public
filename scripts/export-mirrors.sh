#!/bin/bash

export REPOS="HumorChecker dynamolock errors goherokuname runner bookmarkd snippetsd cci mlflappygopher kohrah-ani oversight pglock junk svc gists exp poo alfredemoji pidctl openapigen"

for repo in $REPOS
do
	./scripts/push-repo.sh $repo
done

