package pbdump

import "fmt"

type Kind int

const (
	Complex Kind = iota
	Simple
)

type Field struct {
	Name     string
	Kind     Kind
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

func InjectNames(m StringerMessage, c Context) NamedMessage {
	ret := make(map[string]fmt.Stringer)
	for id, val := range m {
		field, ok := c.Field(id)
		if !ok {
			continue
		}
		if field.Kind == Complex {
			if field.Repeated {
				r := make(NamedMessageRepeated, len(val))
				for i, v := range val {
					r[i] = InjectNames(v.(StringerMessage), field.Context)
				}
				ret[field.Name] = r
			} else {
				ret[field.Name] = InjectNames(val[0].(StringerMessage), field.Context)
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
