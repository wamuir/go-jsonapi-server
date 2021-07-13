SELECT vertices.type,
       vertices.id,
       vertices.attributes,
       vertices.meta
  FROM vertices
 WHERE (vertices.type=$1 AND vertices.id=$2)
