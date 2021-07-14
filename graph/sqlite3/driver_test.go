package backend

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {

	g, err := Connect("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := g.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTransation(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, err := g.Transaction(context.Background(), false)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := tx.Close(); err != nil {
		t.Fatal(err)
	}
}
