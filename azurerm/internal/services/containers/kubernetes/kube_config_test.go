package kubernetes

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseKubeConfig(t *testing.T) {
	testCases := []struct {
		sourceFile string
		expected   KubeConfig
		checkFunc  func(expected KubeConfig, config string) (bool, error)
	}{
		{
			"user_with_token.yml",
			KubeConfig{
				KubeConfigBase: KubeConfigBase{
					APIVersion: "v1",
					Clusters: []clusterItem{
						{
							Name: "test-cluster",
							Cluster: cluster{
								Server: "https://testcluster.net:8080",
							},
						},
					},
					Kind: "Config",
				},
				Users: []userItem{
					{
						Name: "test-user",
						User: user{
							Token: "test-token",
						},
					},
				},
			},
			isValidConfig,
		},
		{
			"user_with_cert.yml",
			KubeConfig{
				KubeConfigBase: KubeConfigBase{
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
				Users: []userItem{
					{
						Name: "test-user",
						User: user{
							ClientCertificteData: "test-client-certificate-data",
							ClientKeyData:        "test-client-key-data",
						},
					},
				},
			},
			isValidConfig,
		},
		{
			"user_with_cert_token.yml",
			KubeConfig{
				KubeConfigBase: KubeConfigBase{
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
			},
			isValidConfig,
		},
		{
			"user_with_no_auth.yml",
			KubeConfig{},
			isInvalidConfig,
		},
		{
			"no_cluster.yml",
			KubeConfig{},
			isInvalidConfig,
		},
		{
			"no_user.yml",
			KubeConfig{},
			isInvalidConfig,
		},
		{
			"user_with_partial_auth.yml",
			KubeConfig{},
			isInvalidConfig,
		},
		{
			"cluster_with_no_server.yml",
			KubeConfig{},
			isInvalidConfig,
		},
	}

	for i, test := range testCases {
		encodedConfig := LoadConfig(test.sourceFile)
		if len(encodedConfig) == 0 {
			t.Fatalf("Test case [%d]: Failed to read config from file '%+v' \n",
				i, test.sourceFile)
		}
		if success, err := test.checkFunc(test.expected, encodedConfig); !success {
			t.Fatalf("Test case [%d]: Failed, config '%+v' with error: '%+v'",
				i, test.sourceFile, err)
		}
	}
}

func isValidConfig(expected KubeConfig, encodedConfig string) (bool, error) {
	result, err := ParseKubeConfig(encodedConfig)
	if err != nil {
		return false, err
	}

	if !reflect.DeepEqual(expected, *result) {
		return false, fmt.Errorf("expected '%+v but got '%+v' with encoded config '%+v'",
			expected, *result, encodedConfig)
	}
	return true, nil
}

func isInvalidConfig(_ KubeConfig, encodedConfig string) (bool, error) {
	if _, err := ParseKubeConfig(encodedConfig); err == nil {
		return false, fmt.Errorf("expected test to throw error but didn't")
	}
	return true, nil
}

func LoadConfig(fileName string) string {
	filePath := filepath.Join("testdata", fileName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return ""
	}

	return string(bytes)
}
