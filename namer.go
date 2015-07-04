package pbdump

/*
type Context interface {
	Has(id int) bool
	Name(id int) string
	IsRepeated(id int) bool
	Context(id int) Context
}

type NamedValue

type namedVariant struct {
	Value
	names map[int]string
}

func (m *NamedMessage) connect(id int, name string) {
	m.names[id] = name
}

func InjectNames(v Value, c Context) NamedMessage {
	ret := namedVariant{v, make(map[int]string)}
	for id, val := range *v.Message() {
		if !c.Has(id) {
			continue
		}
		if c.Context(id) != nil {
			if c.IsRepeated(id) {
				r := make([]Value, len(val))
				for i, v := range val {
					r[i] = InjectNames(v, c.Context(id))
				}
				ret[c.Name(id)] = r
				ret.connect(id, c.Name(id))
			} else {
				ret[c.Name(id)] = InjectNames(val[0], c.Context(id))
			}
		} else {
			if c.IsRepeated(id) {
				r := make([]Value, len(val))
				for i, v := range val {
					r[i] = convertMessageToString(v)
				}
				ret[c.Name(id)] = r
			} else {
				ret[c.Name(id)] = convertMessageToString(val[0])
			}
		}
	}
	return varint{message: ret}
}

func convertMessageToString(m Value) Value {
	if m.Message() != nil {
		return varint{str: string(m.Payload()), payload: m.Payload()}
	}
	return m
}
*/
