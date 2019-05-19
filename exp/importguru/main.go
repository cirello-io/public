// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"go/build"
	"log"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	pkgName := "./testdata/somepkg"
	dir := "."
	spew.Dump(detectDirectImports(pkgName, dir))
}

func detectDirectImports(pkgName, dir string) (map[string][]string, error) {
	directImports := make(map[string][]string)
	pkg, err := build.Import(pkgName, dir, 0)
	if err != nil {
		return directImports, fmt.Errorf("cannot find package (%s): %v", pkgName, err)
	}
	for _, importPkg := range pkg.Imports {
		di := directImports[importPkg]
		di = append(di, pkg.Name)
		directImports[importPkg] = di
		subImports, err := detectDirectImports(pkg.Name, dir)
		if err == nil {
			for pkg, imports := range subImports {
				di := directImports[pkg]
				di = append(di, imports...)
				directImports[pkg] = di
			}
		}
	}
	return directImports, nil
}

type importPair struct {
	origin, target string
}
