package kubernetes

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseKubeConfig(t *testing.T) {
	testCases := []struct {
		sourceFile string
		expected   KubeConfig
	}{
		{
			"user_token_valid.yml",
			KubeConfig{
				APIVersion: "v1",
				Clusters: []clusterItem{
					{
						Name: "test-cluster",
						Cluster: cluster{
							Server: "https://testcluster.net:8080",
						},
					},
				},
				Users: []userItem{
					{
						Name: "test-user",
						User: user{
							Token: "test-token",
						},
					},
				},
				Kind: "Config",
			},
		},
		{
			"user_certs_valid.yml",
			KubeConfig{
				APIVersion: "v1",
				Clusters: []clusterItem{
					{
						Name: "test-cluster",
						Cluster: cluster{
							ClusterAuthorityData: "test-cluster-authority-data",
							Server:               "https://testcluster.org:443",
						},
					},
				},
				Users: []userItem{
					{
						Name: "test-user",
						User: user{
							ClientCertificteData: "test-client-certificate-data",
							ClientKeyData:        "test-client-key-data",
						},
					},
				},
				Contexts: []contextItem{
					{
						Name: "test-cluster",
						Context: context{
							Cluster:   "test-cluster",
							User:      "test-user",
							Namespace: "test-namespace",
						},
					},
				},
				CurrentContext: "test-cluster",
				Kind:           "Config",
				Preferences:    nil,
			},
		},
		{
			"user_both_valid.yml",
			KubeConfig{
				APIVersion: "v1",
				Clusters: []clusterItem{
					{
						Name: "test-cluster",
						Cluster: cluster{
							ClusterAuthorityData: "test-cluster-authority-data",
							Server:               "https://testcluster.org:443",
						},
					},
				},
				Users: []userItem{
					{
						Name: "test-user",
						User: user{
							ClientCertificteData: "test-client-certificate-data",
							ClientKeyData:        "test-client-key-data",
							Token:                "test-token",
						},
					},
				},
				Contexts: []contextItem{
					{
						Name: "test-cluster",
						Context: context{
							Cluster:   "test-cluster",
							User:      "test-user",
							Namespace: "test-namespace",
						},
					},
				},
				CurrentContext: "test-cluster",
				Kind:           "Config",
				Preferences: map[string]interface{}{
					"colors": true,
				},
			},
		},
	}

	for i, test := range testCases {
		encodedConfig := LoadConfig(test.sourceFile)
		if len(encodedConfig) <= 0 {
			t.Fatalf("Test case [%d]: Failed to read config from file %+v",
				i, test.sourceFile)
		}
		result, err := ParseKubeConfig(&encodedConfig)
		if err != nil {
			t.Fatalf("Test case [%d]: Error thrown calling ParseKubeConfig error: '%+v'",
				i, err)
		}
		if !reflect.DeepEqual(test.expected, *result) {
			t.Fatalf("Test case [%d]: Expected '%+v' for config '%+v' - got '%+v'",
				i, test.expected, encodedConfig, *result)
		}
	}
}

func LoadConfig(fileName string) string {
	filePath := filepath.Join("testdata", fileName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
