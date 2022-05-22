package main

import (
	"helloWorld/config"
	"helloWorld/database"
	"testing"
)

func Test_main(t *testing.T) {
	got := Add(4, 6)
	want := 10

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

	got = Multiply(4, 6)
	want = 24

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}

	received := compareString("rajat", "rajat")
	expected := true

	if received != expected {
		t.Errorf("received %t, expected %t", received, expected)
	}
	config.InitConfig()
	config.InitConfiguration()
	database.ConnectToDatabase()
	insertIntoDatabase(database.Get(), "Sanjay Sin")

}
