package somepkg

import (
	"time"

	_ "cirello.io/exp/importguru/testdata/somepkg/subpkg"
)

type SomeType struct {
	t time.Time
}
