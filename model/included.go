package model

type Included []Resource

func (c Included) Merge(m Included) Included {

	for _, r := range m {
		c = c.MergeResource(r)
	}

	return c

}

// Merge a resource into the collection, such as for `Included`
func (c Included) MergeResource(r Resource) Included {

	idx, found := c.FindIndex(r)
	if !found {
		c = append(c, r)
		return c
	}

	for key, rel := range r.Relationships {

		if _, ok := c[idx].Relationships[key]; ok {
			// Relationship already exists
			if c[idx].Relationships[key].Data != nil {
				// And data already exists
				continue
			}
		}

		c[idx].Relationships[key] = rel
	}

	return c
}

func (c Included) FindIndex(r Resource) (int, bool) {

	var idx int
	for i := range c {
		if c[i].Type == r.Type && c[i].Identifier == r.Identifier {
			return i, true
		}
	}
	return idx, false
}
