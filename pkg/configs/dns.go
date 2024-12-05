package configs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type DnsDriverItem struct {
	Name   string
	Uri    string                 `yaml:"uri"`
	Acme   bool                   `yaml:"acme"`
	With   map[string]interface{} `yaml:"with"`
	Shared bool                   `yaml:"shared"`
}

type DnsDriverSection struct {
	data map[string]DnsDriverItem
}

func (d *DnsDriverSection) Get(name string) (DnsDriverItem, bool) {
	item, ok := d.data[name]
	return item, ok
}

func (d *DnsDriverSection) Set(name string, item DnsDriverItem) {
	d.data[name] = item
}

func (d *DnsDriverSection) Has(name string) bool {
	_, ok := d.data[name]
	return ok
}

func (d *DnsDriverSection) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}

	d.data = make(map[string]DnsDriverItem)

	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]
		var item DnsDriverItem
		if val.Kind == yaml.ScalarNode {
			item.Name = key.Value
			item.Uri = val.Value
			d.data[key.Value] = item
		} else if val.Kind == yaml.MappingNode {

			if err := val.Decode(&item); err != nil {
				return err
			}

			item.Name = key.Value
			d.data[key.Value] = item
		} else {
			return fmt.Errorf("expected a scalar or mapping node, got %v", val.Kind)
		}
	}

	return nil
}
