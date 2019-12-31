SELECT COUNT(*)
  FROM edges
 INNER JOIN vertices
    ON (edges.from_rowid=vertices.rowid)
 WHERE (vertices.type=? AND vertices.id=? AND edges.key=?)
