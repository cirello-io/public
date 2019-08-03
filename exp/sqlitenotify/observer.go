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

	"github.com/fsnotify/fsnotify"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:sqlite.db")
	check(err)
	db.SetMaxOpenConns(1)
	defer db.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	notifications := make(chan struct{}, 1024000)
	go func() {
		for range notifications {
			row := db.QueryRow("select count(*) from t1")
			var rowCount int64
			err := row.Scan(&rowCount)
			log.Println(rowCount, err)
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					select {
					case notifications <- struct{}{}:
						log.Println("detected modification", len(notifications), event, event.Name)
					default:
						log.Println("notification overflow", len(notifications))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	<-done
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
