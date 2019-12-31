SELECT vertices.type,
       vertices.id,
       vertices.attributes,
       vertices.meta
  FROM vertices
 WHERE (vertices.type=? AND vertices.id=?)
