package validator

import (
	"regexp"
	"sync"
	"testing"
)

// TestLazyRegexCompile_Basic tests that lazyRegexCompile compiles the regex only once and caches the result.
func TestLazyRegexCompile_Basic(t *testing.T) {
	alphaRegexString := "^[a-zA-Z]+$"
	alphaRegex := lazyRegexCompile(alphaRegexString)

	callCount := 0
	originalFunc := alphaRegex
	alphaRegex = func() *regexp.Regexp {
		callCount++
		return originalFunc()
	}

	// Call the function multiple times
	for i := 0; i < 10; i++ {
		result := alphaRegex()
		if result == nil {
			t.Fatalf("Expected non-nil result")
		}
		if !result.MatchString("test") {
			t.Fatalf("Expected regex to match 'test'")
		}
	}

	if callCount != 10 {
		t.Fatalf("Expected call count to be 10, got %d", callCount)
	}
}

// TestLazyRegexCompile_Concurrent tests that lazyRegexCompile works correctly when called concurrently.
func TestLazyRegexCompile_Concurrent(t *testing.T) {
	alphaRegexString := "^[a-zA-Z]+$"
	alphaRegex := lazyRegexCompile(alphaRegexString)

	var wg sync.WaitGroup
	const numGoroutines = 100

	// Use a map to ensure all results point to the same instance
	results := make(map[*regexp.Regexp]bool)
	var mu sync.Mutex

	// Call the function concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := alphaRegex()
			if result == nil {
				t.Errorf("Expected non-nil result")
			}
			mu.Lock()
			results[result] = true
			mu.Unlock()
		}()
	}
	wg.Wait()

	if len(results) != 1 {
		t.Fatalf("Expected one unique regex instance, got %d", len(results))
	}
}

// TestLazyRegexCompile_Panic tests that if the regex compilation panics, the panic value is propagated consistently.
func TestLazyRegexCompile_Panic(t *testing.T) {
	faultyRegexString := "[a-z"
	alphaRegex := lazyRegexCompile(faultyRegexString)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected a panic, but none occurred")
		}
	}()

	// Call the function, which should panic
	alphaRegex()
}
