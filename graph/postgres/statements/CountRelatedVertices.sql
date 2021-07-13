SELECT COUNT(*)
  FROM edges
 INNER JOIN vertices
    ON (edges.from_rowid=vertices.rowid)
 WHERE (vertices.type=$1 AND vertices.id=$2 AND edges.key=$3)
