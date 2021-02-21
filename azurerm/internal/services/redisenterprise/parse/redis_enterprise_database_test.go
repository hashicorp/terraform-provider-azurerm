package parse

import (
	"testing"
)

func TestRedisEnterpriseDatabaseIDFormatter(t *testing.T) {
	actual := NewRedisEnterpriseDatabaseID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "cluster1", "database1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRedisEnterpriseDatabaseID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RedisEnterpriseDatabaseId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// missing subscriptions
			Input: "/",
			Error: true,
		},
		{
			// missing value for subscriptions
			Input: "/subscriptions/",
			Error: true,
		},
		{
			// missing resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},
		{
			// missing value for resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},
		{
			// missing redisEnterprise
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/",
			Error: true,
		},
		{
			// missing value for redisEnterprise
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/",
			Error: true,
		},
		{
			// missing databases
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1/",
			Error: true,
		},
		{
			// missing value for databases
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1/databases/database1",
			Expected: &RedisEnterpriseDatabaseId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				ClusterName:    "cluster1",
				Name:           "database1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CACHE/REDISENTERPRISE/CLUSTER1/DATABASES/DATABASE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RedisEnterpriseDatabaseID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
