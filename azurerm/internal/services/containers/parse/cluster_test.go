package parse

import (
	"testing"
)

func TestKubernetesClusterID(t *testing.T) {
	testData := []struct {
		input    string
		expected *ClusterId
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups",
			expected: nil,
		},
		{
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello",
			expected: nil,
		},
		{
			// wrong case
			input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/managedclusters/cluster1",
			expected: nil,
		},
		{
			input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/hello/managedClusters/cluster1",
			expected: &ClusterId{
				ManagedClusterName: "cluster1",
				ResourceGroup:      "hello",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)
		actual, err := ClusterID(v.input)

		// if we get something there shouldn't be an error
		if v.expected != nil && err == nil {
			continue
		}

		// if nothing's expected we should get an error
		if v.expected == nil && err != nil {
			continue
		}

		if v.expected == nil && actual == nil {
			continue
		}

		if v.expected == nil && actual != nil {
			t.Fatalf("Expected nothing but got %+v", actual)
		}
		if v.expected != nil && actual == nil {
			t.Fatalf("Expected %+v but got nil", actual)
		}

		if v.expected.ResourceGroup != actual.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but got %q", v.expected.ResourceGroup, actual.ResourceGroup)
		}
		if v.expected.ManagedClusterName != actual.ManagedClusterName {
			t.Fatalf("Expected Name to be %q but got %q", v.expected.ManagedClusterName, actual.ManagedClusterName)
		}
	}
}
