package pbdump

import "testing"

type SimpleContext map[int]FieldDescriptor

func (c SimpleContext) Field(id int) (Field, bool) {
	v, ok := c[id]
	return v, ok
}

type FieldDescriptor struct {
	name       string
	isRepeated bool
	context    Context
}

func (f FieldDescriptor) Name() string {
	return f.name
}

func (f FieldDescriptor) IsRepeated() bool {
	return f.isRepeated
}

func (f FieldDescriptor) Context() Context {
	return f.context
}

func TestNamingOfMessageWithInt(t *testing.T) {
	m := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{StringerVarint(42)},
		},
		rawPayload: nil,
	}
	c := make(SimpleContext)
	c[1] = FieldDescriptor{
		name: "foo", isRepeated: false, context: nil,
	}
	out := InjectNames(m, c)
	if v, ok := out["foo"]; !ok {
		t.Fatalf("No expected field: '%v'", out)
	} else if i, ok := v.(StringerVarint); !ok {
		t.Fatalf("Expected varint, got '%v'", v)
	} else if i != 42 {
		t.Fatalf("Expected 42, got %v", i)
	}
}

func TestNamingOfMessageWithRepeatedInt(t *testing.T) {
	m := StringerMessage{
		attributes: map[int]StringerRepeated{
			2: StringerRepeated{StringerVarint(42), StringerVarint(55)},
		},
		rawPayload: nil,
	}
	c := make(SimpleContext)
	c[2] = FieldDescriptor{
		name: "foo", isRepeated: true, context: nil,
	}
	out := InjectNames(m, c)
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
	m1 := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{StringerVarint(42)},
		},
		rawPayload: nil,
	}
	m2 := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{StringerString("test")},
		},
		rawPayload: nil,
	}
	root := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{m1},
			2: StringerRepeated{m2},
		},
		rawPayload: nil,
	}

	c1 := make(SimpleContext)
	c2 := make(SimpleContext)
	croot := make(SimpleContext)

	c1[1] = FieldDescriptor{
		name: "id", isRepeated: false, context: nil,
	}
	c2[1] = FieldDescriptor{
		name: "field", isRepeated: false, context: nil,
	}
	croot[1] = FieldDescriptor{
		name: "m1", isRepeated: false, context: &c1,
	}
	croot[2] = FieldDescriptor{
		name: "m2", isRepeated: false, context: &c2,
	}

	out := InjectNames(root, croot)

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
	m1 := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{StringerVarint(42)},
		},
		rawPayload: nil,
	}

	m2 := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{StringerVarint(55)},
		},
		rawPayload: nil,
	}

	root := StringerMessage{
		attributes: map[int]StringerRepeated{
			1: StringerRepeated{m1, m2},
		},
		rawPayload: nil,
	}

	c1 := make(SimpleContext)
	c1[1] = FieldDescriptor{
		name: "id", isRepeated: false, context: nil,
	}

	croot := make(SimpleContext)
	croot[1] = FieldDescriptor{
		name: "submsg", isRepeated: true, context: &c1,
	}

	out := InjectNames(root, croot)
	if submsg, ok := out["submsg"]; !ok {
		t.Fatalf("Expected submsg got: '%v'", out)
	} else if v, ok := submsg.(NamedMessageRepeated); !ok {
		t.Fatalf("Expected repeated, got: '%v'", out)
	} else if len(v) != 2 {
		t.Fatalf("Expected slice of two, got: '%v'", v)
	} else if v[0].(NamedMessage)["id"].(StringerVarint) != 42 || v[1].(NamedMessage)["id"].(StringerVarint) != 55 {
		t.Fatalf("Expected 42 and 55, got '%v'", v)
	}
}
