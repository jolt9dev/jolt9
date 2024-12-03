package configs_test

import (
	"strings"
	"testing"

	"github.com/jolt9dev/jolt9/pkg/configs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type EnvTestRoot struct {
	Envs configs.EnvsSection `yaml:"envs"`
}

func TestUnmarshalEnvFromRoot(t *testing.T) {
	yamlData := `
envs:
  kv:
    vars:
      name: "VAR1"

  assignment:
    vars:
      - name="VAR2 VALUE"

  seqkv:
    vars:
      - name: "VAR3"
    imports: ["./.env"]
`

	root := &EnvTestRoot{}

	dec := yaml.NewDecoder(strings.NewReader(yamlData))
	err := dec.Decode(&root)
	if err != nil {
		t.Fatal(err)
	}
	envs := root.Envs
	assert.Equal(t, 3, envs.Len())

	kv, ok := envs.Get("kv")
	assert.True(t, ok)
	assert.NotNil(t, kv)
	assert.Equal(t, 1, len(kv.Vars))
	assert.Equal(t, "VAR1", kv.Vars["name"])

	assignment, ok := envs.Get("assignment")
	assert.True(t, ok)
	assert.NotNil(t, assignment)
	assert.Equal(t, 1, len(assignment.Vars))
	assert.Equal(t, "VAR2 VALUE", assignment.Vars["name"])

	seqkv, ok := envs.Get("seqkv")
	assert.True(t, ok)
	assert.NotNil(t, seqkv)
	assert.Equal(t, 1, len(seqkv.Vars))
	assert.Equal(t, "VAR3", seqkv.Vars["name"])
	assert.Equal(t, 1, len(seqkv.Imports))
	assert.Equal(t, "./.env", seqkv.Imports[0])
}
