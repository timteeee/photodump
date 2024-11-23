package env_test

import (
	"fmt"
	"os"
	"photodump/internal/env"
	"testing"
)

type Thing struct {
	Foo   int    `env:"FOO"`
	Bar   bool   `env:"BAR"`
	Hello string `env:"HELLO"`
}

func TestParse(t *testing.T) {
	expectedFoo := 69
	expectedBar := true
	expectedHello := "world"

	os.Setenv("FOO", fmt.Sprintf("%d", expectedFoo))
	os.Setenv("BAR", fmt.Sprintf("%t", expectedBar))
	os.Setenv("HELLO", expectedHello)

	thing := &Thing{}

	if err := env.Parse(thing); err != nil {
		t.Fatal(err)
	}

	if thing.Foo != expectedFoo {
		t.Fatalf("wrong value for field Foo\nGot: %d\nExpected: %d", thing.Foo, expectedFoo)
	}
	if thing.Bar != expectedBar {
		t.Fatalf("wrong value for field Foo\nGot: %t\nExpected: %t", thing.Bar, expectedBar)
	}
	if thing.Hello != expectedHello {
		t.Fatalf("wrong value for field Foo\nGot: %s\nExpected: %s", thing.Hello, expectedHello)
	}
}
