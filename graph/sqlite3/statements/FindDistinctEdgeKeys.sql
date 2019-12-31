SELECT DISTINCT key
  FROM edges
 INNER JOIN vertices
    ON (edges.from_rowid=vertices.rowid)
 WHERE (vertices.type=? AND vertices.id=?)
