package configs

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Represents the configuration of a secret vault
// The `uri` field is the URI of the vault configuration.
// It must at least include the scheme which instructs
// which driver to use. Currently on sops with env files
// is supported
//
// vaults:
//
//	  name: sops:./path/to/file.env
//		 name:
//			url: sops:./path/to/file.env
//			shared: true
//		 name:
//			url: sops:
//			shared: false
//		    with:
//		      file: ./path/to/file.env
type VaultItem struct {
	Name   string
	Uri    string `yaml:"uri"`
	Shared bool
	driver string
	With   map[string]interface{} `yaml:"with"`
}

func (v *VaultItem) Driver() string {
	if v.driver == "" {
		if v.Uri != "" {
			index := strings.Index(v.Uri, ":")
			if index != -1 {
				v.driver = v.Uri[:index]
			}
		}
	}

	return v.driver
}

type VaultsSection struct {
	data map[string]VaultItem
}

func (v *VaultsSection) Get(name string) (VaultItem, bool) {
	item, ok := v.data[name]
	return item, ok
}

func (v *VaultsSection) Set(name string, item VaultItem) {
	v.data[name] = item
}

func (v *VaultsSection) Has(name string) bool {
	_, ok := v.data[name]
	return ok
}

func (v *VaultsSection) Len() int {
	return len(v.data)
}

func (v *VaultsSection) Names() []string {
	names := make([]string, 0, len(v.data))
	for name := range v.data {
		names = append(names, name)
	}

	return names
}

func (v *VaultsSection) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}

	v.data = make(map[string]VaultItem)

	// vaults:
	//   context:
	//     uri: vault://path/to/secret

	// or
	// vaults:
	//   context: vault://path/to/secret

	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]

		if val.Kind == yaml.ScalarNode {
			v.data[key.Value] = VaultItem{Name: key.Value, Uri: val.Value}
			continue
		} else if val.Kind == yaml.MappingNode {
			var item VaultItem
			if err := val.Decode(&item); err != nil {
				return err
			}

			item.Name = key.Value
			v.data[key.Value] = item
		} else {
			return fmt.Errorf("expected a scalar or mapping node, got %v", val.Kind)
		}
	}

	return nil
}
