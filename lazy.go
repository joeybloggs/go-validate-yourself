//go:build go1.21

package validator

import (
	"regexp"
	"sync"
)

func lazyRegexCompile(str string) func() *regexp.Regexp {
	return sync.OnceValue(func() *regexp.Regexp {
		return regexp.MustCompile(str)
	})
}
