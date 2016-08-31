// Package goherokuname generates Heroku-like random names. Such "block-black-0231".
//
// It supports both decimal and hexadecimal suffixes. It supports customized
// rendering, in which the delimiter, the length of suffix and the acceptable
// characters for suffix are tweakable.
package goherokuname

import (
	"fmt"
	"math/rand"
	"time"
)

var seed int64

func randomSource() *rand.Rand {
	lseed := seed
	if 0 == lseed {
		lseed = time.Now().UnixNano()
	}
	return rand.New(rand.NewSource(lseed))
}

// HaikunateCustom generates a Heroku-like random name in which the delimiter,
// length and acceptable characters for suffix are tweakable.
func HaikunateCustom(delimiter string, toklen int, tokchars string) string {
	r := randomSource()

	noun := nouns[r.Intn(len(nouns))]
	adjective := adjectives[r.Intn(len(adjectives))]

	token := ""
	for i := 0; i < toklen; i++ {
		token += string(tokchars[r.Intn(len(tokchars))])
	}

	return adjective + delimiter + noun + delimiter + token
}

// Haikunate generate standard Heroku-like random name, with "-" as delimiter,
// decimal with 4 digits.
func Haikunate() string {
	return HaikunateCustom("-", 4, "0123456789")
}

// HaikunateHex generate standard Heroku-like random name, with "-" as
// delimiter, hexadecimal with 4 digits.
func HaikunateHex() string {
	return HaikunateCustom("-", 4, "0123456789abcdef")
}

// Ubuntu generates a Ubuntu-like random name in which the delimiter is tweakable.
func Ubuntu(delimiter, letter string) (string, error) {
	r := randomSource()
	chosenLetter := letter[0]

	var filteredAdjectives, filteredNouns []string

	for _, adjective := range adjectives {
		if adjective[0] == chosenLetter {
			filteredAdjectives = append(filteredAdjectives, adjective)
		}
	}

	for _, noun := range nouns {
		if noun[0] == chosenLetter {
			filteredNouns = append(filteredNouns, noun)
		}
	}

	if len(filteredAdjectives) == 0 || len(filteredNouns) == 0 {
		return "", fmt.Errorf("could not find words in dictionary starting with \"%c\" for adjective and/or noun", chosenLetter)
	}

	adjective := filteredAdjectives[r.Intn(len(filteredAdjectives))]
	noun := filteredNouns[r.Intn(len(filteredNouns))]

	return adjective + delimiter + noun, nil
}
