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
	"database/sql"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:sqlite.db")
	check(err)
	db.SetMaxOpenConns(1)
	defer db.Close()

	tick := time.Tick(1 * time.Nanosecond)
	for range tick {
		stmts := []string{
			"BEGIN TRANSACTION",
			"INSERT INTO t1 VALUES (1)",
			"COMMIT",
		}
		for _, stmt := range stmts {
			res, err := db.Exec(stmt)
			if err != nil {
				spew.Dump(res, err)
			}
		}
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
