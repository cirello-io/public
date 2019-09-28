#!/bin/bash -x

git remote add $1 git@github.com:cirello-io/$1.git
git subtree split --rejoin --prefix=$1
git subtree pull --prefix=$1 $1 master
