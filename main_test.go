package main

import (
	"testing"
)

func Test_main(t *testing.T) {
	got := Add(4, 6)
	want := 11

	if got == want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
