package configs_test

import (
	"strings"
	"testing"

	"github.com/jolt9dev/jolt9/pkg/configs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type VaultTestRoot struct {
	Vaults configs.VaultsSection `yaml:"vaults"`
}

func TestVaultSection(t *testing.T) {
	yamlData := `
vaults:
  default: "sops:./path/to/file.env"

  target1:
    uri: "sops:./path/to/target1.env"
    shared: true

  target2:
    uri: "sops:"
    with:
      file: "./path/to/target2.env"	
`

	root := &VaultTestRoot{}

	dec := yaml.NewDecoder(strings.NewReader(yamlData))
	err := dec.Decode(&root)

	if err != nil {
		t.Fatal(err)
	}

	vaults := root.Vaults

	assert.Equal(t, 3, vaults.Len())

	vault, ok := vaults.Get("default")
	assert.True(t, ok)
	assert.NotNil(t, vault)
	assert.Equal(t, "sops:./path/to/file.env", vault.Uri)
	assert.False(t, vault.Shared)

	vault1, ok := vaults.Get("target1")
	assert.True(t, ok)
	assert.NotNil(t, vault1)
	assert.Equal(t, "sops:./path/to/target1.env", vault1.Uri)
	assert.True(t, vault1.Shared)

	vault2, ok := vaults.Get("target2")
	assert.True(t, ok)
	assert.NotNil(t, vault2)
	assert.Equal(t, "sops:", vault2.Uri)
	assert.False(t, vault2.Shared)
	assert.Equal(t, "./path/to/target2.env", vault2.With["file"])
}
