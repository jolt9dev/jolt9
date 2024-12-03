package vps

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/docker/cli/cli/connhelper"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	"github.com/jolt9dev/jolt9/pkg/ssh"
)

func CreateConfig() {

}

func SetupServer(config *ssh.Config, sudoPassword string) error {

	client, err := ssh.NewClient(config)
	if err != nil {
		return err
	}

	err = AddSudoerAsUser(client, config.User, sudoPassword)
	if err != nil {
		return err
	}
	err = InstallAftDirectories(client)
	if err != nil {
		return err
	}

	err = InstallDefaultPackages(client)
	if err != nil {
		return err
	}

	err = InstallDocker(client)
	if err != nil {
		return err
	}

	return nil
}

func InstallAftDirectories(client ssh.Client) error {
	dirs := []string{
		"/opt/jolt9",
		"/opt/jolt9/compose",
		"/opt/jolt9/mnt/data",
		"/opt/jolt9/mnt/backup",
		"/opt/jolt9/mnt/etc",
		"/opt/jolt9/mnt/etc/ssl/certs",
		"/opt/jolt9/mnt/logs",
		"/opt/jolt9/scripts",
	}

	return InstallDirectories(client, dirs)
}

func InstallDocker(client ssh.Client) error {
	cmd := `
if [ -x "$(command -v docker)" ]; then
	echo "Docker already installed"
else
	curl -fsSL https://get.docker.com -o ~/get-docker.sh
	sudo sh ~/get-docker.sh
fi

GROUP=$(getent group docker)

if [[ $GROUP == *"$USER"* ]]; then 
	echo "User already in docker group"
else
	if [ -z $GROUP ]; then
		sudo groupadd docker
	fi 
	sudo usermod -aG docker $USER
fi
`

	return client.Shell(nil, os.Stdout, os.Stderr, cmd)
}

func InstallDirectories(client ssh.Client, dirs []string) error {

	cmds := ""

	for _, dir := range dirs {
		cmds += `
if [ -d ` + dir + ` ]; then
	echo "Directory ` + dir + ` already exists"
else
	echo "Creating directory ` + dir + `"
	sudo mkdir -p ` + dir + `
fi
`
	}

	return client.Shell(nil, os.Stdout, os.Stderr, cmds)
}

func AddSudoerAsUser(client ssh.Client, user string, password string) error {
	cmd := `
echo "starting"
if [ -f /etc/sudoers.d/` + user + ` ]; then
	echo "User already exists"	
else
	echo '` + password + `' | sudo -S touch /etc/sudoers.d/` + user + `
	echo "` + user + ` ALL=(ALL) NOPASSWD:ALL" | sudo tee -a /etc/sudoers.d/` + user + `	
fi
`
	println(cmd)
	out, err := client.Output(cmd)
	if err != nil {
		return err
	}

	println(out)

	return nil
}

func InstallDefaultPackages(client ssh.Client) error {

	pkgs := []string{
		"curl",
		"wget",
		"ca-certificates",
		"apt-transport-https",
		"software-properties-common",
		"gnupg",
		"cifs-utils",
		"bat",
		"btop",
		"ripgrep",
		"eza",
		"net-tools",
		"tre-command",
	}

	return InstallAptPackages(client, pkgs)
}

func InstallAptPackages(client ssh.Client, packages []string) error {
	pkgs := strings.Join(packages, `\
    `)

	cmd := `
export debian_frontend=noninteractive
sudo apt-get update && sudo apt-get upgrade -y 
sudo apt-get install -y ` + pkgs

	println(cmd)
	err := client.Shell(nil, os.Stdout, os.Stderr, cmd)
	if err != nil {
		return err
	}

	return nil
}

func GetOsInfo(client ssh.Client) (*OsInfo, error) {
	cmd := `
if [ -f /etc/os-release ]; then
	cat /etc/os-release
elif [ -f /etc/redhat-release ]; then
	cat /etc/redhat-release
else
	echo "ID=unknown"
fi
`
	content, err := client.Output(cmd)
	if err != nil {
		return nil, err
	}

	env, err := godotenv.Unmarshal(content)
	if err != nil {
		return nil, err
	}

	osInfo := &OsInfo{}
	if id, ok := env["ID"]; ok {
		osInfo.Id = id
	}

	if idLike, ok := env["ID_LIKE"]; ok {
		osInfo.IdLike = idLike
	}

	if version, ok := env["VERSION"]; ok {
		osInfo.Version = version
	}

	if versionId, ok := env["VERSION_ID"]; ok {
		osInfo.VersionId = versionId
	}

	if versionCodename, ok := env["VERSION_CODENAME"]; ok {
		osInfo.VersionCodename = versionCodename
	}

	if platformId, ok := env["PLATFORM_ID"]; ok {
		osInfo.PlatformId = platformId
	}

	if prettyName, ok := env["PRETTY_NAME"]; ok {
		osInfo.PrettyName = prettyName
	}

	return osInfo, nil
}

type OsInfo struct {
	Id              string
	IdLike          string
	Version         string
	VersionId       string
	VersionCodename string
	PlatformId      string
	PrettyName      string
}

func InstallAftDockerNetworks(config ssh.Config) error {

	sshHost := `ssh://` + config.User + `@` + config.Host + `:` + strconv.Itoa(config.Port)

	helper, err := connhelper.GetConnectionHelper(sshHost)

	if err != nil {
		return err
	}

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	cli, err := client.NewClientWithOpts(client.WithHost(helper.Host), client.WithHTTPClient(httpClient), client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return err
	}

	hasAftFrontend := false
	hasAftBackend := false
	for _, n := range networks {
		if n.Name == "aft_frontend" {
			hasAftFrontend = true
		}

		if n.Name == "aft_backend" {
			hasAftBackend = true
		}
	}

	println("hasAftBackend", hasAftBackend)

	if !hasAftFrontend {
		_, err := cli.NetworkCreate(context.Background(), "aft_frontend", network.CreateOptions{

			Driver: "bridge",

			IPAM: &network.IPAM{
				Driver: "default",
				Config: []network.IPAMConfig{
					{
						Subnet:  "172.19.0.0/16",
						Gateway: "172.19.0.1",
					},
				},
			},
			Options: map[string]string{
				"com.docker.network.bridge.name":       "aft_frontend",
				"com.docker.network.bridge.enable_icc": "false",
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
