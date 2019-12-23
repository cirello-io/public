// Copyright 2019 github.com/ucirello and cirello.io. All rights reserved.
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

// Package qrng is quantum Random Number Generator Client for ANU.edu.au.
package qrng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const maxLength = 1024

// Uint8 returns byte-sized random numbers.
func Uint8(length int) ([]uint8, error) {
	if length < 1 {
		return nil, fmt.Errorf("length is too small: %v", length)
	} else if length > maxLength {
		return nil, fmt.Errorf("length is too large: %v", length)
	}
	var r struct {
		Data []uint8 `json:"data"`
	}
	if err := read("uint8", length, 0, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

// Uint16 returns word-sized random numbers.
func Uint16(length int) ([]uint16, error) {
	if length < 1 {
		return nil, fmt.Errorf("length is too small: %v", length)
	} else if length > maxLength {
		return nil, fmt.Errorf("length is too large: %v", length)
	}
	var r struct {
		Data []uint16 `json:"data"`
	}
	if err := read("uint16", length, 0, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

// Hex16 returns word-sized hexadecimal random number blocks.
func Hex16(length, blockSize int) ([]string, error) {
	if length < 1 {
		return nil, fmt.Errorf("length is too small: %v", length)
	} else if length > maxLength {
		return nil, fmt.Errorf("length is too large: %v", length)
	}
	if blockSize < 1 {
		return nil, fmt.Errorf("block size is too small: %v", blockSize)
	} else if blockSize > maxLength {
		return nil, fmt.Errorf("block size is too large: %v", blockSize)
	}
	var r struct {
		Data []string `json:"data"`
	}
	if err := read("hex16", length, blockSize, &r); err != nil {
		return nil, err
	}
	return r.Data, nil
}

// Reader is a global, shared instance of a quantum based random number
// generator.
var Reader io.Reader = &reader{}

type reader struct{}

func (r *reader) Read(p []byte) (n int, err error) {
	// fast path
	if len(p) <= maxLength {
		b, err := Uint8(len(p))
		copy(p, b)
		return len(b), err
	}
	var b []byte
	for c := len(p); c > 0; c -= maxLength {
		b, err = Uint8(maxLength)
		copy(p[len(p)-c:], b)
		n += len(b)
		if err != nil {
			break
		}
	}
	return n, err

}

// Read is a helper function that calls Reader.Read using io.ReadFull.
func Read(p []byte) (n int, err error) {
	return io.ReadFull(Reader, p)
}

func read(t string, l, bs int, target interface{}) error {
	url := fmt.Sprintf("https://qrng.anu.edu.au/API/jsonI.php?type=%s&length=%d&size=%d", t, l, bs)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("cannot load random numbers: %w", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return fmt.Errorf("cannot parse response: %w", err)
	}
	return nil
}
