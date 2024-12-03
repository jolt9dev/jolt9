package configs

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type DnsProvider struct {
	Name   string
	Inputs map[string]interface{}
}

type DnsProvidersSection struct {
	data map[string]DnsProvider
}

func (d *DnsProvidersSection) Get(name string) (DnsProvider, bool) {
	item, ok := d.data[name]
	return item, ok
}

func (d *DnsProvidersSection) Set(name string, item DnsProvider) {
	d.data[name] = item
}

func (d *DnsProvidersSection) Has(name string) bool {
	_, ok := d.data[name]
	return ok
}

func (d *DnsProvidersSection) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}

	d.data = make(map[string]DnsProvider)

	for i := 0; i < len(value.Content); i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]

		if val.Kind == yaml.ScalarNode {

		} else if val.Kind == yaml.MappingNode {

			var item DnsProvider
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

type ComposeSection struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type ContextsSection struct {
}

type TraefikSection struct {
	Ingnore bool
	Enabled bool
}

type ContextItem struct {
	Vaults    []string
	Envs      []string
	Dns       string
	SshConfig string
	Servers   []string
}

type UseEnvsSection struct {
	Vars    []string
	Include []string
}

type SecretItem struct {
	Name     string
	Key      string
	Generate bool
	Special  string
	Digits   bool
	Lower    bool
	Upper    bool
}

type UseVaultsSection struct {
	Include []string
	Secrets []SecretItem
}

type ProjectConfig struct {
	Vaults VaultsSection
	Envs   EnvsSection
}
