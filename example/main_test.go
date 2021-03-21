package main

import "testing"

func TestFive (t *testing.T) {
	five := five()
	if five != 5 {
		t.Fatalf("Error, expected 5 but got %d", five)
	}
}
