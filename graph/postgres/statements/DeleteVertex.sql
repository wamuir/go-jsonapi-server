DELETE
  FROM vertices
 WHERE (vertices.type=$1 AND vertices.id=$2)
