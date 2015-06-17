package pbdump

import "testing"

type SimpleContext map[int]Field

func (c SimpleContext) Field(id int) (Field, bool) {
	v, ok := c[id]
	return v, ok
}

func TestNamingOfMessageWithInt(t *testing.T) {
	m := make(StringerMessage)
	m[1] = StringerRepeated{StringerVarint(42)}
	c := make(SimpleContext)
	c[1] = Field{
		Name: "foo", Type: Int64,
		Repeated: false, Context: nil,
	}
	out := Apply(m, c)
	if v, ok := out["foo"]; !ok {
		t.Fatalf("No expected field: '%v'", out)
	} else if i, ok := v.(StringerVarint); !ok {
		t.Fatalf("Expected varint, got '%v'", v)
	} else if i != 42 {
		t.Fatalf("Expected 42, got %v", i)
	}
}

func TestNamingOfMessageWithRepeatedInt(t *testing.T) {
	m := make(StringerMessage)
	m[2] = StringerRepeated{
		StringerVarint(42), StringerVarint(55),
	}
	c := make(SimpleContext)
	c[2] = Field{
		Name: "foo", Type: Int64,
		Repeated: true, Context: nil,
	}
	out := Apply(m, c)
	if v, ok := out["foo"]; !ok {
		t.Fatalf("No expected field: '%v'", out)
	} else if s, ok := v.(StringerRepeated); !ok {
		t.Fatalf("Expected repeated, got '%v'", v)
	} else if len(s) != 2 {
		t.Fatalf("Expected slice of 2, got '%v'", s)
	} else if s[0].(StringerVarint) != 42 || s[1].(StringerVarint) != 55 {
		t.Fatalf("Unexpected slice content: '%v'", s)
	}
}

func TestNamingOfMessageWithMessage(t *testing.T) {
	m1 := make(StringerMessage)
	m2 := make(StringerMessage)
	root := make(StringerMessage)

	m1[1] = StringerRepeated{StringerVarint(42)}
	m2[1] = StringerRepeated{StringerString("test")}
	root[1] = StringerRepeated{m1}
	root[2] = StringerRepeated{m2}

	c1 := make(SimpleContext)
	c2 := make(SimpleContext)
	croot := make(SimpleContext)

	c1[1] = Field{
		Name: "id", Type: Uint64,
		Repeated: false, Context: nil,
	}
	c2[1] = Field{
		Name: "field", Type: String,
		Repeated: false, Context: nil,
	}
	croot[1] = Field{
		Name: "m1", Type: Message,
		Repeated: false, Context: &c1,
	}
	croot[2] = Field{
		Name: "m2", Type: Message,
		Repeated: false, Context: &c2,
	}

	out := Apply(root, croot)

	if v1, ok := out["m1"]; !ok {
		t.Fatalf("Missing field m1: '%v'", out)
	} else if id, ok := v1.(NamedMessage)["id"]; !ok {
		t.Fatalf("Expected field m1:id: '%v'", out)
	} else if id.(StringerVarint) != 42 {
		t.Fatalf("Expected 42, got: '%v'", id)
	}

	if v2, ok := out["m2"]; !ok {
		t.Fatalf("Missing field m2: '%v'", out)
	} else if field, ok := v2.(NamedMessage)["field"]; !ok {
		t.Fatalf("Expected field m2:field: '%v'", out)
	} else if field.(StringerString) != "test" {
		t.Fatalf("Expected \"test\", got: '%v'", out)
	}
}

func TestNamingOfMessageWithRepeatedMessage(t *testing.T) {
	m1 := make(StringerMessage)
	m1[1] = StringerRepeated{StringerVarint(42)}

	m2 := make(StringerMessage)
	m2[1] = StringerRepeated{StringerVarint(55)}

	root := make(StringerMessage)
	root[1] = StringerRepeated{m1, m2}

	c1 := make(SimpleContext)
	c1[1] = Field{
		Name: "id", Type: Uint64,
		Repeated: false, Context: nil,
	}

	croot := make(SimpleContext)
	croot[1] = Field{
		Name: "submsg", Type: Message,
		Repeated: true, Context: &c1,
	}

	out := Apply(root, croot)
	if submsg, ok := out["submsg"]; !ok {
		t.Fatalf("Expected submsg got: '%v'", out)
	} else if v, ok := submsg.(NamedMessageRepeated); !ok {
		t.Fatalf("Expected repeated, got: '%v'", out)
	} else if len(v) != 2 {
		t.Fatalf("Expected slice of two, got: '%v'", v)
	} else if v[0]["id"].(StringerVarint) != 42 || v[1]["id"].(StringerVarint) != 55 {
		t.Fatalf("Expected 42 and 55, got '%v'", v)
	}
}
