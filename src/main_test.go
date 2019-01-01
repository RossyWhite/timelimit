package main

import "testing"

func TestRun(t *testing.T) {
	err := myHandler()

	if err != nil {
		t.Error(err)
	}
}
