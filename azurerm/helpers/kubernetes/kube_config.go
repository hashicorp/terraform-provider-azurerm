package kubernetes

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type clusterItem struct {
	Name    string  `yaml:"name"`
	Cluster cluster `yaml:"cluster"`
}

type cluster struct {
	ClusterAuthorityData string `yaml:"certificate-authority-data"`
	Server               string `yaml:"server"`
}

type userItem struct {
	Name string `yaml:"name"`
	User user   `yaml:"user"`
}

type user struct {
	ClientCertificteData string `yaml:"client-certificate-data"`
	Token                string `yaml:"token"`
	ClientKeyData        string `yaml:"client-key-data"`
}

type userItemAAD struct {
	Name string  `yaml:"name"`
	User userAAD `yaml:"user"`
}

type userAAD struct {
	AuthProvider authProvider `yaml:"auth-provider"`
}

type authProvider struct {
	Name   string        `yaml:"name"`
	Config configAzureAD `yaml:"config"`
}

type configAzureAD struct {
	APIServerID string `yaml:"apiserver-id,omitempty"`
	ClientID    string `yaml:"client-id,omitempty"`
	TenantID    string `yaml:"tenant-id,omitempty"`
}

type contextItem struct {
	Name    string  `yaml:"name"`
	Context context `yaml:"context"`
}

type context struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type KubeConfigBase struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []clusterItem          `yaml:"clusters"`
	Contexts       []contextItem          `yaml:"contexts,omitempty"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind,omitempty"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
}

type KubeConfig struct {
	KubeConfigBase `yaml:",inline"`
	Users          []userItem `yaml:"users"`
}

type KubeConfigAAD struct {
	KubeConfigBase `yaml:",inline"`
	Users          []userItemAAD `yaml:"users"`
}

func ParseKubeConfig(config string) (*KubeConfig, error) {
	if config == "" {
		return nil, fmt.Errorf("Cannot parse empty config")
	}

	var kubeConfig KubeConfig

	if err := yaml.Unmarshal([]byte(config), &kubeConfig); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML config with error %+v", err)
	}
	if len(kubeConfig.Clusters) == 0 || len(kubeConfig.Users) == 0 {
		return nil, fmt.Errorf("Config %+v contains no valid clusters or users", kubeConfig)
	}
	u := kubeConfig.Users[0].User
	if u.Token == "" && (u.ClientCertificteData == "" || u.ClientKeyData == "") {
		return nil, fmt.Errorf("Config requires either token or certificate auth for user %+v", u)
	}
	c := kubeConfig.Clusters[0].Cluster
	if c.Server == "" {
		return nil, fmt.Errorf("Config has invalid or non existent server for cluster %+v", c)
	}

	return &kubeConfig, nil
}

func ParseKubeConfigAAD(config string) (*KubeConfigAAD, error) {
	if config == "" {
		return nil, fmt.Errorf("Cannot parse empty config")
	}

	var kubeConfig KubeConfigAAD
	if err := yaml.Unmarshal([]byte(config), &kubeConfig); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML config with error %+v", err)
	}
	if len(kubeConfig.Clusters) == 0 || len(kubeConfig.Users) == 0 {
		return nil, fmt.Errorf("Config %+v contains no valid clusters or users", kubeConfig)
	}

	c := kubeConfig.Clusters[0].Cluster
	if c.Server == "" {
		return nil, fmt.Errorf("Config has invalid or non existent server for cluster %+v", c)
	}

	return &kubeConfig, nil
}
