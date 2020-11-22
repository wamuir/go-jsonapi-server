package backend

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/wamuir/go-jsonapi-server/graph"
)

func (tx *transaction) InsertEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string, position int, meta []byte) error {

	result, err := tx.Prepared["InsertEdge"].Exec(
		key,
		position,
		string(meta),
		fromVertexType,
		fromVertexID,
		toVertexType,
		toVertexID,
	)
	sqliteErr, ok := err.(sqlite3.Error)
	if ok && sqliteErr.Code == sqlite3.ErrConstraint {
		return graph.ErrConflict
	} else if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return graph.ErrNoRows
	}

	return nil
}

func (tx *transaction) InsertVertex(vertexType, vertexID string, attributes, meta []byte) error {

	result, err := tx.Prepared["InsertVertex"].Exec(
		vertexType,
		vertexID,
		string(attributes),
		string(meta),
	)
	sqliteErr, ok := err.(sqlite3.Error)
	if ok && sqliteErr.Code == sqlite3.ErrConstraint {
		return graph.ErrConflict
	} else if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return graph.ErrNoRows
	}

	return nil
}

func (tx *transaction) DeleteVertex(vertexType, vertexID string) error {

	result, err := tx.Prepared["DeleteVertex"].Exec(
		vertexType,
		vertexID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return graph.ErrNoRows
	}

	return nil
}

func (tx *transaction) DeleteEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string) error {

	result, err := tx.Prepared["DeleteEdge"].Exec(
		fromVertexType,
		fromVertexID,
		key,
		toVertexType,
		toVertexID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	} else if count == 0 {
		return graph.ErrNoRows
	}

	return nil
}

func (tx *transaction) CountVertices(vertexType string) (int64, error) {

	var count int64

	result := tx.Prepared["CountVertices"].QueryRow(
		vertexType,
	)

	err := result.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (tx *transaction) FindVertices(vertexType string, limit, offset int64, sort string) ([]graph.Vertex, error) {

	var vertices []graph.Vertex

	key := "FindVertices"
	if sort == "newest" {
		key = "FindVerticesNewest"
	}

	rows, err := tx.Prepared[key].Query(
		vertexType,
		limit,
		offset,
	)
	if err != nil {
		return vertices, err
	}

	defer rows.Close()

	for rows.Next() {

		var vertex graph.Vertex

		err := rows.Scan(
			&vertex.Type,
			&vertex.Identifier,
			&vertex.Attributes,
			&vertex.Meta,
		)
		if err != nil {
			return vertices, err
		}

		vertices = append(vertices, vertex)
	}

	return vertices, nil
}

func (tx *transaction) FindVertex(vertexType, vertexID string) (graph.Vertex, error) {

	var vertex graph.Vertex

	row := tx.Prepared["FindVertex"].QueryRow(
		vertexType,
		vertexID,
	)

	err := row.Scan(
		&vertex.Type,
		&vertex.Identifier,
		&vertex.Attributes,
		&vertex.Meta,
	)
	if err == sql.ErrNoRows {
		return vertex, graph.ErrNoRows
	} else if err != nil {
		return vertex, err
	}

	return vertex, nil
}

func (tx *transaction) FindDistinctEdgeKeys(fromVertexType, fromVertexID string) ([]string, error) {

	var keys []string

	rows, err := tx.Prepared["FindDistinctEdgeKeys"].Query(
		fromVertexType,
		fromVertexID,
	)
	if err != nil {
		return keys, err
	}

	defer rows.Close()

	for rows.Next() {
		var key string

		err = rows.Scan(&key)
		if err != nil {
			return keys, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (tx *transaction) CountRelatedVertices(fromVertexType, fromVertexID, key string) (int64, error) {

	var count int64

	result := tx.Prepared["CountRelatedVertices"].QueryRow(
		fromVertexType,
		fromVertexID,
		key,
	)

	err := result.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (tx *transaction) FindEdges(fromVertexType, fromVertexID, key string, limit, offset int64) ([]graph.Edge, error) {

	var edges []graph.Edge

	rows, err := tx.Prepared["FindEdges"].Query(
		fromVertexType,
		fromVertexID,
		key,
		limit,
		offset,
	)
	if err != nil {
		return edges, err
	}

	defer rows.Close()

	for rows.Next() {

		var edge graph.Edge

		err = rows.Scan(
			&edge.From.Type,
			&edge.From.Identifier,
			&edge.From.Attributes,
			&edge.From.Meta,
			&edge.To.Type,
			&edge.To.Identifier,
			&edge.To.Attributes,
			&edge.To.Meta,
			&edge.Key,
			&edge.Meta,
		)
		if err != nil {
			return edges, err
		}

		edges = append(edges, edge)

	}

	return edges, nil
}

func (tx *transaction) FindEdge(fromVertexType, fromVertexID, toVertexType, toVertexID, key string) (graph.Edge, error) {

	var edge graph.Edge

	row := tx.Prepared["FindEdge"].QueryRow(
		fromVertexType,
		fromVertexID,
		toVertexType,
		toVertexID,
		key,
	)

	err := row.Scan(
		&edge.From.Type,
		&edge.From.Identifier,
		&edge.From.Attributes,
		&edge.From.Meta,
		&edge.To.Type,
		&edge.To.Identifier,
		&edge.To.Attributes,
		&edge.To.Meta,
		&edge.Key,
		&edge.Meta,
	)
	if err == sql.ErrNoRows {
		return edge, graph.ErrNoRows
	} else if err != nil {
		return edge, err
	}

	return edge, nil
}
