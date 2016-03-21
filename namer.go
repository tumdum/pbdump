package pbdump

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Field interface {
	Name() string
	IsRepeated() bool
	Context() Context
}

type Context interface {
	Field(id int) (Field, bool)
}

type NamedMessage map[string]fmt.Stringer

func (m NamedMessage) String() string {
	return fmt.Sprint(map[string]fmt.Stringer(m))
}

func (m NamedMessage) MarshalJSON() ([]byte, error) {
	b, err := json.MarshalIndent(map[string]fmt.Stringer(m), "", "\t")
	if err != nil {
		return nil, err
	}
	return b, nil
}

type NamedMessageRepeated []fmt.Stringer

func (m NamedMessageRepeated) String() string {
	return fmt.Sprint([]fmt.Stringer(m))
}

func InjectNames(m StringerMessage, c Context) NamedMessage {
	ret := make(map[string]fmt.Stringer)
	for id, val := range m.attributes {
		field, ok := c.Field(id)
		if !ok || strings.HasPrefix(field.Name(), "oper_") {
			continue
		}
		if field.Context() != nil {
			if field.IsRepeated() {
				r := make(NamedMessageRepeated, len(val))
				for i, v := range val {
					switch v := v.(type) {
					case StringerString, StringerDouble, StringerVarint:
						r[i] = v
					case StringerMessage:
						r[i] = InjectNames(v, field.Context())
					default:
						panic(v)
					}
				}
				ret[field.Name()] = r
			} else {
				log.Printf("%#v %T\n", val[0], val[0])
				switch v := val[0].(type) {
				case StringerString, StringerDouble, StringerVarint:
					ret[field.Name()] = v
				case StringerMessage:
					ret[field.Name()] = InjectNames(val[0].(StringerMessage), field.Context())
				default:
					panic(v)
				}
			}
		} else {
			if field.IsRepeated() {
				ret[field.Name()] = val
			} else {
				ret[field.Name()] = convertUnexpectedMsgToString(val[0])
			}
		}
	}
	return ret
}

func convertUnexpectedMsgToString(v fmt.Stringer) fmt.Stringer {
	if v, ok := v.(StringerMessage); ok {
		return StringerString(v.rawPayload)
	}
	return v
}
