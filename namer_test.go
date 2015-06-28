package pbdump

/*
import "testing"

type Field struct {
	name       string
	isRepeated bool
	context    Context
}

func (f Field) Name() string {
	return f.name
}

func (f Field) IsRepeated() bool {
	return f.isRepeated
}

func (f Field) Context() Context {
	return f.context
}

type SimpleContext map[int]Field

func (c SimpleContext) Has(id int) bool {
	_, ok := c[id]
	return ok
}

func (c SimpleContext) Name(id int) string {
	return c[id].name
}

func (c SimpleContext) IsRepeated(id int) bool {
	return c[id].isRepeated
}

func (c SimpleContext) Context(id int) Context {
	return c[id].context
}

func TestNamingOfMessageWithInt(t *testing.T) {
	m := StringerMessage{make(map[int]StringerRepeated), nil}
	m.attributes[1] = StringerRepeated{StringerVarint(42)}
	c := make(SimpleContext)
	c[1] = Field{name: "foo", isRepeated: false, context: nil}
	out := InjectNames(m, c)

	if v, ok := out["foo"]; !ok {
		t.Fatalf("No expected field: '%v'", out)
	} else if i, ok := v.(StringerVarint); !ok {
		t.Fatalf("Expected varint, got '%v'", v)
	} else if i != 42 {
		t.Fatalf("Expected 42, got %v", i)
	}
}

func TestNamingOfMessageWithStringParsedAsMessage(t *testing.T) {
	m := StringerMessage{make(map[int]StringerRepeated), nil}
	sub := StringerMessage{make(map[int]StringerRepeated), []byte("test")}
	m.attributes[1] = StringerRepeated{sub}

	c := make(SimpleContext)
	c[1] = Field{name: "foo", isRepeated: false, context: nil}

	out := InjectNames(m, c)
	v := out["foo"]
	if s, ok := v.(StringerString); !ok {
		t.Fatalf("Expected string, got '%#v'", v)
	} else if s != "test" {
		t.Fatalf("Expected 'test', got '%v'", s)
	}
}

func TestNamingOfMessageWithRepeatedInt(t *testing.T) {
	m := StringerMessage{make(map[int]StringerRepeated), nil}
	m.attributes[2] = StringerRepeated{
		StringerVarint(42), StringerVarint(55),
	}
	c := make(SimpleContext)
	c[2] = Field{
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
	m1 := StringerMessage{make(map[int]StringerRepeated), nil}
	m2 := StringerMessage{make(map[int]StringerRepeated), nil}
	root := StringerMessage{make(map[int]StringerRepeated), nil}

	m1.attributes[1] = StringerRepeated{StringerVarint(42)}
	m2.attributes[1] = StringerRepeated{StringerString("test")}
	root.attributes[1] = StringerRepeated{m1}
	root.attributes[2] = StringerRepeated{m2}

	c1 := make(SimpleContext)
	c2 := make(SimpleContext)
	croot := make(SimpleContext)

	c1[1] = Field{
		name: "id", isRepeated: false, context: nil,
	}
	c2[1] = Field{
		name: "field", isRepeated: false, context: nil,
	}
	croot[1] = Field{
		name: "m1", isRepeated: false, context: &c1,
	}
	croot[2] = Field{
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
	m1 := StringerMessage{make(map[int]StringerRepeated), nil}
	m1.attributes[1] = StringerRepeated{StringerVarint(42)}

	m2 := StringerMessage{make(map[int]StringerRepeated), nil}
	m2.attributes[1] = StringerRepeated{StringerVarint(55)}

	root := StringerMessage{make(map[int]StringerRepeated), nil}
	root.attributes[1] = StringerRepeated{m1, m2}

	c1 := make(SimpleContext)
	c1[1] = Field{
		name: "id", isRepeated: false, context: nil,
	}

	croot := make(SimpleContext)
	croot[1] = Field{
		name: "submsg", isRepeated: true, context: &c1,
	}

	out := InjectNames(root, croot)
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
*/
