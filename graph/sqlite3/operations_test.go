package backend

import (
	"context"
	"testing"
)

func TestInsertEdge(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)

	if err := tx.InsertEdge("typeA", "idA", "typeB", "idB", "key", 0, nil); err != nil {
		t.Fatal(err)
	}
}

func TestInsertVertex(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	if err := tx.InsertVertex("typeA", "idA", nil, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteVertex(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)

	if err := tx.DeleteVertex("typeA", "idA"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteEdge(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeB", "idB", "key", 0, nil)

	if err := tx.DeleteEdge("typeA", "idA", "typeB", "idB", "key"); err != nil {
		t.Fatal(err)
	}
}

func TestCountVertices(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeA", "idB", nil, nil)

	i, err := tx.CountVertices("typeA")
	if err != nil {
		t.Fatal(err)
	} else if i != 2 {
		t.Fatalf("Got %d, want %d", i, 2)
	}
}

func TestFindVertices(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeA", "idB", nil, nil)

	v, err := tx.FindVertices("typeA", 10, 0, "")
	if err != nil {
		t.Fatal(err)
	} else if len(v) != 2 {
		t.Fatalf("%v", v)
	}
}

func TestFindVertex(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)

	_, err := tx.FindVertex("typeA", "idA")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindDistinctEdgeKeys(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)
	_ = tx.InsertVertex("typeC", "idC", nil, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeB", "idB", "keyA", 0, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeC", "idC", "keyB", 0, nil)

	e, err := tx.FindDistinctEdgeKeys("typeA", "idA")
	if err != nil {
		t.Fatal(err)
	} else if len(e) != 2 {
		t.Fatalf("%v", e)
	}
}

func TestCountRelatedVertices(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)
	_ = tx.InsertVertex("typeC", "idC", nil, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeB", "idB", "keyA", 0, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeC", "idC", "keyA", 0, nil)

	e, err := tx.CountRelatedVertices("typeA", "idA", "keyA")
	if err != nil {
		t.Fatal(err)
	} else if e != 2 {
		t.Fatalf("%v", e)
	}
}

func TestFindEdges(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)
	_ = tx.InsertVertex("typeC", "idC", nil, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeB", "idB", "keyA", 0, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeC", "idC", "keyA", 0, nil)

	e, err := tx.FindEdges("typeA", "idA", "keyA", 10, 0)
	if err != nil {
		t.Fatal(err)
	} else if len(e) != 2 {
		t.Fatalf("%v", e)
	}
}

func TestEdge(t *testing.T) {

	g, _ := Connect("file::memory:?cache=shared")
	defer g.Close()

	tx, _ := g.Transaction(context.Background(), false)
	defer tx.Close()

	_ = tx.InsertVertex("typeA", "idA", nil, nil)
	_ = tx.InsertVertex("typeB", "idB", nil, nil)
	_ = tx.InsertEdge("typeA", "idA", "typeB", "idB", "key", 0, nil)

	_, err := tx.FindEdge("typeA", "idA", "typeB", "idB", "key")
	if err != nil {
		t.Fatal(err)
	}
}
