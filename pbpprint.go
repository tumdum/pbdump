package pbdump

import "fmt"

type Type int

const (
	Int32 Type = iota
	Int64
	Uint32
	Uint64
	Double
	String
	Message
)

type Field struct {
	Name     string
	Type     Type
	Repeated bool
	Context  Context
}

type Context interface {
	Field(id int) (Field, bool)
}

type NamedMessage map[string]fmt.Stringer

func (m NamedMessage) String() string {
	return fmt.Sprint(map[string]fmt.Stringer(m))
}

type NamedMessageRepeated []NamedMessage

func (m NamedMessageRepeated) String() string {
	return fmt.Sprint([]NamedMessage(m))
}

func Apply(m StringerMessage, c Context) NamedMessage {
	ret := make(map[string]fmt.Stringer)
	for id, val := range m {
		field, ok := c.Field(id)
		if !ok {
			continue
		}
		if field.Type == Message {
			if field.Repeated {
				r := make(NamedMessageRepeated, len(val))
				for i, v := range val {
					r[i] = Apply(v.(StringerMessage), field.Context)
				}
				ret[field.Name] = r
			} else {
				ret[field.Name] = Apply(val[0].(StringerMessage), field.Context)
			}
		} else {
			if field.Repeated {
				ret[field.Name] = val
			} else {
				ret[field.Name] = val[0]
			}
		}
	}
	return ret
}
