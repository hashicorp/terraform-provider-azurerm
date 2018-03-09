package kubernetes

import (
	"encoding/base64"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type ClusterItem struct {
	Name    string  `yaml:"name"`
	Cluster Cluster `yaml:"cluster"`
}

type Cluster struct {
	ClusterAuthorityData string `yaml:"certificate-authority-data"`
	Server               string `yaml:"server"`
}

type UserItem struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

type User struct {
	ClientCertificteData string `yaml:"client-certificate-data"`
	Token                string `yaml:"token"`
	ClientKeyData        string `yaml:"client-key-data"`
}

type ContextItem struct {
	Name    string  `yaml:"name"`
	Context Context `yaml:"context"`
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []ClusterItem          `yaml:"clusters"`
	Users          []UserItem             `yaml:"users"`
	Contexts       []ContextItem          `yaml:"contexts,omitempty"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind,omitempty"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
}

func ParseKubeConfig(config *string) (*KubeConfig, error) {
	if config == nil {
		return nil, fmt.Errorf("Cannot parse nil config")
	}

	configBytes, err := base64.StdEncoding.DecodeString(*config)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode base64 string %+v with error %+v", *config, err)
	}

	configStr := string(configBytes)
	if len(configStr) <= 0 {
		return nil, fmt.Errorf("Decoded config is empty %+v", *config)
	}

	var kubeConfig KubeConfig
	err = yaml.Unmarshal([]byte(configStr), &kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML config with error %+v", err)
	}
	if len(kubeConfig.Clusters) <= 0 || len(kubeConfig.Users) <= 0 {
		return nil, fmt.Errorf("Config %+v contains no valid clusters or users", kubeConfig)
	}

	return &kubeConfig, nil
}
