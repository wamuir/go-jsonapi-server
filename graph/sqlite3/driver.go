package backend

import (
	"context"
	"database/sql"
	"embed"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite3 driver for database/sql
	"github.com/wamuir/go-jsonapi-server/graph"
)

//go:embed schema/*.sql
//go:embed statements/*.sql
var fs embed.FS

// Connect opens a connection to a SQLite3 database and returns a graph, as
// *graph.Graph.  Argument `dsn` (data source name) is connection string.
func Connect(dsn string) (graph.Graph, error) {

	var g graph.Graph

	g, err := newConnection(dsn)
	if err != nil {
		return nil, err
	}

	return g, nil
}

type connection struct {
	*sql.DB
	closer func() error
}

// Returns a connection object.
func newConnection(dsn string) (*connection, error) {

	var conn connection

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	conn = connection{db, db.Close}

	err = conn.setup()
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

func (conn connection) Close() error {
	return conn.closer()
}

func (conn connection) Transaction(ctx context.Context, readOnly bool) (graph.Tx, error) {

	var tx graph.Tx

	tx, err := conn.newTransaction(ctx, true, readOnly)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (conn connection) setup() error {

	tx, err := conn.newTransaction(context.TODO(), false, false)
	if err != nil {
		return err
	}

	keys := []string{
		"CreateTableVertices.sql",
		"CreateIndexVertices.sql",
		"CreateTableEdges.sql",
		"CreateIndexEdges.sql",
		"CreateIndexEdgesFk.sql",
	}

	for _, k := range keys {
		p := filepath.Join("schema", k)
		data, err := fs.ReadFile(p)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(string(data)); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type transaction struct {
	*sql.Tx
	Prepared map[string]*sql.Stmt
}

func (conn connection) newTransaction(ctx context.Context, prepare, readOnly bool) (*transaction, error) {

	options := sql.TxOptions{
		Isolation: 0,
		ReadOnly:  readOnly,
	}

	t, err := conn.DB.BeginTx(ctx, &options)
	if err != nil {
		return nil, err
	}

	tx := transaction{
		Tx:       t,
		Prepared: make(map[string]*sql.Stmt),
	}
	if prepare {
		keys := []string{
			"CountRelatedVertices.sql",
			"CountVertices.sql",
			"DeleteEdge.sql",
			"DeleteVertex.sql",
			"FindDistinctEdgeKeys.sql",
			"FindEdge.sql",
			"FindEdges.sql",
			"FindVertex.sql",
			"FindVertices.sql",
			"FindVerticesNewest.sql",
			"InsertEdge.sql",
			"InsertVertex.sql",
		}
		for _, k := range keys {
			p := filepath.Join("statements", k)
			data, err := fs.ReadFile(p)
			if err != nil {
				return nil, err
			}

			name := strings.TrimSuffix(k, ".sql")
			tx.Prepared[name], err = tx.Prepare(string(data))
			if err != nil {
				return nil, err
			}
		}
	}

	return &tx, nil
}

func (tx *transaction) Close() error {

	if err := tx.Rollback(); err != nil {
		return err
	}

	for _, statement := range tx.Prepared {
		if err := statement.Close(); err != nil {
			return err
		}

	}

	return nil
}
