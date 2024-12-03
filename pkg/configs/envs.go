package configs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type EnvsSection struct {
	data map[string]EnvItem
}

func (e *EnvsSection) Get(name string) (EnvItem, bool) {
	item, ok := e.data[name]
	return item, ok
}

func (e *EnvsSection) Set(name string, item EnvItem) {
	e.data[name] = item
}

func (e *EnvsSection) Has(name string) bool {
	_, ok := e.data[name]
	return ok
}

func (e *EnvsSection) Len() int {
	return len(e.data)
}

func (e *EnvsSection) UnmarshalYAML(value *yaml.Node) error {

	// envs:
	//   context:
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}

	e.data = make(map[string]EnvItem)

	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i]
		println(key.Value)
		val := value.Content[i+1]

		var item EnvItem
		if err := val.Decode(&item); err != nil {
			return err
		}

		e.data[key.Value] = item
	}

	return nil
}

type EnvItem struct {
	Vars    map[string]string
	Imports []string
}

func (e *EnvItem) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}
	// envs:
	//  context:
	//    vars:
	//      key: value
	// or
	// envs:
	//   name:
	//     vars:
	//       - key=value

	e.Vars = make(map[string]string)
	e.Imports = make([]string, 0)

	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]

		switch key.Value {
		case "vars":
			if val.Kind == yaml.MappingNode {
				for j := 0; j < len(val.Content); j += 2 {
					k := val.Content[j]
					v := val.Content[j+1]

					if (k.Kind != yaml.ScalarNode) || (v.Kind != yaml.ScalarNode) {
						return fmt.Errorf("expected scalar nodes, got %v and %v", k.Kind, v.Kind)
					}

					e.Vars[k.Value] = v.Value
				}
			} else if val.Kind == yaml.SequenceNode {
				for j := 0; j < len(val.Content); j++ {
					item := val.Content[j]

					if item.Kind == yaml.ScalarNode {
						kk, vv, err := ParseKeyValue(item.Value)
						if err != nil {
							return err
						}

						e.Vars[kk] = vv
						continue
					}

					if item.Kind != yaml.MappingNode {
						return fmt.Errorf("expected a mapping node or scalar string node, got %v", item.Kind)
					}

					for k := 0; k < len(item.Content); k += 2 {
						kk := item.Content[k]
						vv := item.Content[k+1]

						if (kk.Kind != yaml.ScalarNode) || (vv.Kind != yaml.ScalarNode) {
							return fmt.Errorf("expected scalar nodes, got %v and %v", kk.Kind, vv.Kind)
						}

						e.Vars[kk.Value] = vv.Value
					}
				}
			} else {
				return fmt.Errorf("expected a mapping or sequence node, got %v", val.Kind)
			}
		case "imports":
			if val.Kind != yaml.SequenceNode {
				return fmt.Errorf("expected a sequence node, got %v", val.Kind)
			}

			for j := 0; j < len(val.Content); j++ {
				item := val.Content[j]

				if item.Kind != yaml.ScalarNode {
					return fmt.Errorf("expected a scalar node, got %v", item.Kind)
				}

				e.Imports = append(e.Imports, item.Value)
			}
		default:
			return fmt.Errorf("unexpected key %q", key.Value)
		}
	}

	return nil
}
