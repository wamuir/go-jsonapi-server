package schema

import (
	"testing"

	"github.com/wamuir/go-jsonapi-core"
)

func TestValidate(t *testing.T) {

	var d core.Document

	// empty primary data
	d = core.New()
	if r, err := Validate(d); err != nil {
		t.Fatal(err)
	} else if r.Valid() {
		t.Fatalf("Got %v, want %v", r.Valid(), false)
	}

	// valid doc w/resource as primary data
	d = core.New()
	d.Data = core.Resource{
		Type:       "foo",
		Identifier: "bar",
	}
	if r, err := Validate(d); err != nil {
		t.Fatal(err)
	} else if !r.Valid() {
		t.Fatalf("Got %v, want %v", r.Valid(), true)
	}

	// regex for type/id fields
	d = core.New()
	d.Data = core.Resource{
		Type:       "foo/bar",
		Identifier: "/baz",
	}
	if r, err := Validate(d); err != nil {
		t.Fatal(err)
	} else if r.Valid() {
		t.Fatalf("Got %v, want %v", r.Valid(), false)
	}

	// missing type
	d = core.New()
	d.Data = core.Resource{
		Identifier: "/baz",
	}
	if r, err := Validate(d); err != nil {
		t.Fatal(err)
	} else if r.Valid() {
		t.Fatalf("Got %v, want %v", r.Valid(), false)
	}


	// valid doc/collection as primary data
	d = core.New()
	d.Data = core.Collection{
		core.Resource{
			Type:       "foo",
			Identifier: "bar",
		},
		core.Resource{
			Type:       "foo",
			Identifier: "baz",
		},
	}
	if r, err := Validate(d); err != nil {
		t.Fatal(err)
	} else if !r.Valid() {
		t.Fatalf("Got %v, want %v", r.Valid(), true)
	}
}
