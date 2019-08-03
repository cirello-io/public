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
	"flag"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/xerrors"
)

var verbose = flag.Bool("verbose", false, "")

func main() {
	flag.Parse()
	const dsn = "file:sqlite.db?_auto_vacuum=full&_busy_timeout=5000&_journal=DELETE&_locking=NORMAL&mode=rw&_mutex=full&_secure_delete=true&_sync=EXTRA&_txlock=exclusive"
	db, err := sql.Open("sqlite3", dsn)
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
			for {
				row := db.QueryRow("select v from t1")
				var i int64
				err := row.Scan(&i)
				if isLocked(err) {
					log.Println(i, err, len(notifications), "retrying")
					continue
				}
				log.Println(i, err, len(notifications))
				break
			}
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
						if *verbose {
							log.Println("detected modification", len(notifications), event, event.Name)
						}
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

func isLocked(err error) bool {
	var sqliteErr sqlite3.Error
	if xerrors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrBusy {
		return true
	}
	return false
}
