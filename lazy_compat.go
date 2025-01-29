//go:build !go1.21

package validator

import (
	"regexp"
	"sync"
)

// Copied and adapted from go1.21 stdlib's sync.OnceValue for backwards compatibility:
// OnceValue returns a function that invokes f only once and returns the value
// returned by f. The returned function may be called concurrently.
//
// If f panics, the returned function will panic with the same value on every call.
func onceValue(f func() *regexp.Regexp) func() *regexp.Regexp {
	var (
		once   sync.Once
		valid  bool
		p      interface{}
		result *regexp.Regexp
	)
	g := func() {
		defer func() {
			p = recover()
			if !valid {
				panic(p)
			}
		}()
		result = f()
		f = nil
		valid = true
	}
	return func() *regexp.Regexp {
		once.Do(g)
		if !valid {
			panic(p)
		}
		return result
	}
}

func lazyRegexCompile(str string) func() *regexp.Regexp {
	return onceValue(func() *regexp.Regexp {
		return regexp.MustCompile(str)
	})
}
