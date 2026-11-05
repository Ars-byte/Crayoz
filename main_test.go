package main

import (
	"testing"

	"Crayoz/src"
)

func TestVersionDefined(t *testing.T) {
	if src.Version == "" {
		t.Fatal("expected src.Version to not be empty")
	}
}
