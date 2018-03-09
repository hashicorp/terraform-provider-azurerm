package kubernetes

import (
	"reflect"
	"testing"
)

func TestParseKubeConfig(t *testing.T) {
	testCases := []struct {
		encoded  string
		expected KubeConfig
	}{
		{
			`YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIHNlcnZlcjogaHR0cHM6Ly90ZXN0Y2x1c3Rlci5uZXQ6ODA4MAogIG5hbWU6IHRlc3QtY2x1c3Rlcgp1c2VyczoKLSBuYW1lOiB0ZXN0LXVzZXIKICB1c2VyOgogICAgdG9rZW46IHRlc3QtdG9rZW4Ka2luZDogQ29uZmln`,
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
			}},
		{
			`YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiB0ZXN0LWNsdXN0ZXItYXV0aG9yaXR5LWRhdGEKICAgIHNlcnZlcjogaHR0cHM6Ly90ZXN0Y2x1c3Rlci5vcmc6NDQzCiAgbmFtZTogdGVzdC1jbHVzdGVyCmNvbnRleHRzOgotIGNvbnRleHQ6CiAgICBjbHVzdGVyOiB0ZXN0LWNsdXN0ZXIKICAgIHVzZXI6IHRlc3QtdXNlcgogICAgbmFtZXNwYWNlOiB0ZXN0LW5hbWVzcGFjZQogIG5hbWU6IHRlc3QtY2x1c3RlcgpjdXJyZW50LWNvbnRleHQ6IHRlc3QtY2x1c3Rlcgp1c2VyczoKLSBuYW1lOiB0ZXN0LXVzZXIKICB1c2VyOgogICAgY2xpZW50LWNlcnRpZmljYXRlLWRhdGE6IHRlc3QtY2xpZW50LWNlcnRpZmljYXRlLWRhdGEKICAgIGNsaWVudC1rZXktZGF0YTogdGVzdC1jbGllbnQta2V5LWRhdGEKa2luZDogQ29uZmln`,
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
			}},
		{
			`YXBpVmVyc2lvbjogdjEKY2x1c3RlcnM6Ci0gY2x1c3RlcjoKICAgIGNlcnRpZmljYXRlLWF1dGhvcml0eS1kYXRhOiB0ZXN0LWNsdXN0ZXItYXV0aG9yaXR5LWRhdGEKICAgIHNlcnZlcjogaHR0cHM6Ly90ZXN0Y2x1c3Rlci5vcmc6NDQzCiAgbmFtZTogdGVzdC1jbHVzdGVyCmNvbnRleHRzOgotIGNvbnRleHQ6CiAgICBjbHVzdGVyOiB0ZXN0LWNsdXN0ZXIKICAgIHVzZXI6IHRlc3QtdXNlcgogICAgbmFtZXNwYWNlOiB0ZXN0LW5hbWVzcGFjZQogIG5hbWU6IHRlc3QtY2x1c3RlcgpjdXJyZW50LWNvbnRleHQ6IHRlc3QtY2x1c3Rlcgp1c2VyczoKLSBuYW1lOiB0ZXN0LXVzZXIKICB1c2VyOgogICAgY2xpZW50LWNlcnRpZmljYXRlLWRhdGE6IHRlc3QtY2xpZW50LWNlcnRpZmljYXRlLWRhdGEKICAgIGNsaWVudC1rZXktZGF0YTogdGVzdC1jbGllbnQta2V5LWRhdGEKICAgIHRva2VuOiB0ZXN0LXRva2VuCmtpbmQ6IENvbmZpZwpwcmVmZXJlbmNlczoKICBjb2xvcnM6IHRydWU=`,
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
			}},
	}

	for i, test := range testCases {
		result, err := ParseKubeConfig(&test.encoded)
		if err != nil {
			t.Fatalf("Error thrown calling ParseKubeConfig in test case %d error: '%+v'",
				i, err)
		}
		if !reflect.DeepEqual(test.expected, *result) {
			t.Fatalf("Test case [%d]: Expected '%+v' for config '%+v' - got '%+v'",
				i, test.expected, test.encoded, *result)
		}
	}
}
