SELECT vertices.type,
       vertices.id,
       vertices.attributes,
       vertices.meta
  FROM vertices
 WHERE (vertices.type=$1)
 ORDER BY vertices.id ASC
 LIMIT $2
OFFSET $3
