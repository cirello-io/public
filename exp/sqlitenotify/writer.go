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
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/xerrors"
)

var slow = flag.Bool("slow", false, "")
var interactive = flag.Bool("interactive", false, "")

func main() {
	flag.Parse()
	const dsn = "file:sqlite.db?_auto_vacuum=full&_busy_timeout=5000&_journal=DELETE&_locking=NORMAL&mode=rw&_mutex=full&_secure_delete=true&_sync=EXTRA&_txlock=exclusive"
	db, err := sql.Open("sqlite3", dsn)
	check(err)
	db.SetMaxOpenConns(1)
	defer db.Close()

	// typical user operation latency
	tick := time.Tick(150 * time.Millisecond)
	for range tick {
		const retries = 5
		checkLock := func(err error) bool {
			if err == nil {
				return false
			}
			if isLocked(err) {
				log.Println(err, "retrying")
				return true
			}
			check(err)
			return false
		}
		for i := 0; i < retries; i++ {
			tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
				Isolation: sql.LevelSerializable,
			})
			if checkLock(err) {
				continue
			}

			var i int
			row := tx.QueryRow("SELECT v FROM t1 WHERE id = 1")
			checkLock(row.Scan(&i))

			i++

			_, err = tx.Exec("UPDATE t1 SET v = $1 WHERE id = 1", i)
			if checkLock(err) {
				continue
			}

			// typical local file operation latency
			wait := 75 * time.Millisecond
			if *slow {
				// typical network latency
				wait = 250 * time.Millisecond
			}
			time.Sleep(wait)

			if checkLock(tx.Commit()) {
				tx.Rollback()
				continue
			}
			break
		}
		fmt.Print(".")
		if *interactive {
			var i string
			fmt.Scanln(&i)
		}
	}
}

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func isLocked(err error) bool {
	var sqliteErr sqlite3.Error
	if xerrors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrBusy {
		return true
	}
	return false
}
