SELECT vertices.type,
       vertices.id,
       vertices.attributes,
       vertices.meta
  FROM vertices
 WHERE (vertices.type=?)
 ORDER BY vertices.id ASC
 LIMIT ?
OFFSET ?
