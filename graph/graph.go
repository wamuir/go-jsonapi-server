package graph

import (
	"context"
	"errors"
)

type Edge struct {
	From Vertex
	To   Vertex
	Key  string
	Meta []byte
}

type Vertex struct {
	Type       string
	Identifier string
	Attributes []byte
	Meta       []byte
}

var (
	ErrNoRows   = errors.New("no rows in result set")
	ErrConflict = errors.New("unique constraint violation in graph")
)

type Graph interface {
	Close() error
	Transaction(ctx context.Context, readOnly bool) (Tx, error)
}

type Tx interface {
	Close() error
	Commit() error
	CountRelatedVertices(fromVertexType, fromVertexID, key string) (int64, error)
	CountVertices(vertexType string) (int64, error)
	DeleteEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string) error
	DeleteVertex(vertexType, vertexID string) error
	FindDistinctEdgeKeys(fromVertexType, fromVertexID string) ([]string, error)
	FindEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string) (Edge, error)
	FindEdges(fromVertexType, fromVertexID, key string, limit, offset int64) ([]Edge, error)
	FindVertex(vertexType, vertexID string) (Vertex, error)
	FindVertices(vertexType string, limit, offset int64, sort string) ([]Vertex, error)
	InsertEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string, position int, meta []byte) error
	InsertVertex(vertexType, vertexID string, attributes, meta []byte) error
}
