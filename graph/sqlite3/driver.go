package backend

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // SQLite3 driver for database/sql
	"github.com/wamuir/go-jsonapi-server/graph"
)

// Connect opens a connection to a SQLite3 database and returns a graph, as
// *graph.Graph.  Argument `dsn` (data source name) is connection string.
func Connect(dsn string) (graph.Graph, error) {

	var g graph.Graph

	g, err := newConnection(dsn)
	if err != nil {
		return g, err
	}

	return g, nil
}

type connection struct{ *sql.DB }

// Returns a connection object.
func newConnection(dsn string) (*connection, error) {

	var conn connection

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return &conn, err
	}

	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		return &conn, err
	}

	conn = connection{db}

	err = conn.setup()
	if err != nil {
		return &conn, err
	}

	return &conn, nil
}

func (conn connection) Transaction(ctx context.Context, readOnly bool) (graph.Tx, error) {

	var tx graph.Tx

	tx, err := conn.newTransaction(ctx, true, readOnly)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (conn connection) setup() error {

	tx, err := conn.newTransaction(context.TODO(), false, false)
	if err != nil {
		return err
	}

	keys := []string{
		"CreateTableVertices",
		"CreateIndexVertices",
		"CreateTableEdges",
		"CreateIndexEdges",
	}

	for _, k := range keys {
		var statement string = string(schema[k])
		_, err := tx.Exec(statement)
		if err != nil {
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
		return &transaction{}, err
	}

	tx := transaction{
		Tx:       t,
		Prepared: make(map[string]*sql.Stmt),
	}
	if prepare {
		for name, statement := range statements {
			tx.Prepared[name], err = tx.Prepare(string(statement))
			if err != nil {
				return &transaction{}, err
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
