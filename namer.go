package pbdump

import "fmt"

type Context interface {
	Has(id int) bool
	Name(id int) string
	IsRepeated(id int) bool
	Context(id int) Context
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
		if !c.Has(id) {
			continue
		}
		if c.Context(id) != nil {
			if c.IsRepeated(id) {
				r := make(NamedMessageRepeated, len(val))
				for i, v := range val {
					r[i] = InjectNames(v.(StringerMessage), c.Context(id))
				}
				ret[c.Name(id)] = r
			} else {
				ret[c.Name(id)] = InjectNames(val[0].(StringerMessage), c.Context(id))
			}
		} else {
			if c.IsRepeated(id) {
				ret[c.Name(id)] = val
			} else {
				ret[c.Name(id)] = val[0]
			}
		}
	}
	return ret
}
