package conf

import (
	"fmt"
	"reflect"

	. "github.com/expr-lang/expr/checker/nature"
	"github.com/expr-lang/expr/internal/deref"
	"github.com/expr-lang/expr/types"
)

func Env(env any) Nature {
	if env == nil {
		return Nature{
			NatureBase: NatureBase{
				TypeName: reflect.TypeOf(map[string]any{}).String(),
				Strict:   true,
			},
			Type: reflect.TypeOf(map[string]any{}),
		}
	}

	switch env := env.(type) {
	case types.Map:
		return env.Nature()
	}

	v := reflect.ValueOf(env)
	d := deref.Value(v)

	switch d.Kind() {
	case reflect.Struct:
		return Nature{
			NatureBase: NatureBase{
				TypeName: v.Type().String(),
				Strict:   true,
			},
			Type: v.Type(),
		}

	case reflect.Map:
		n := Nature{
			NatureBase: NatureBase{
				TypeName: v.Type().String(),
				Fields:   make(map[string]Nature, v.Len()),
			},
			Type: v.Type(),
		}

		for _, key := range v.MapKeys() {
			elem := v.MapIndex(key)
			if !elem.IsValid() || !elem.CanInterface() {
				panic(fmt.Sprintf("invalid map value: %s", key))
			}

			face := elem.Interface()

			switch face.(type) {
			case types.Map:
				n.Fields[key.String()] = face.(types.Map).Nature()

			default:
				if face == nil {
					n.Fields[key.String()] = Nature{NatureBase: NatureBase{Nil: true}}
					continue
				}
				n.Fields[key.String()] = Nature{
					NatureBase: NatureBase{
						TypeName: reflect.TypeOf(face).String(),
					},
					Type: reflect.TypeOf(face),
				}
			}

		}

		return n
	}

	panic(fmt.Sprintf("unknown type %T", env))
}
