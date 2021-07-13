SELECT DISTINCT key
  FROM edges
 INNER JOIN vertices
    ON (edges.from_rowid=vertices.rowid)
 WHERE (vertices.type=$1 AND vertices.id=$2)
