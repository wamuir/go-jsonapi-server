SELECT from_vertex.type,
       from_vertex.id,
       from_vertex.attributes,
       from_vertex.meta,
       to_vertex.type,
       to_vertex.id,
       to_vertex.attributes,
       to_vertex.meta,
       edges.key,
       edges.meta
  FROM edges
 INNER JOIN vertices from_vertex
    ON (edges.from_rowid=from_vertex.rowid)
 INNER JOIN vertices to_vertex
    ON (edges.to_rowid=to_vertex.rowid)
 WHERE (from_vertex.type=? AND from_vertex.id=? AND edges.key=?)
 ORDER BY edges.position ASC
 LIMIT ?
OFFSET ?
