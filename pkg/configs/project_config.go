package configs

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
