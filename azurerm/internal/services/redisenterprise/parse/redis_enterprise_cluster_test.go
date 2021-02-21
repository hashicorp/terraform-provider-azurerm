package parse

import (
	"testing"
)

func TestRedisEnterpriseClusterIDFormatter(t *testing.T) {
	actual := NewRedisEnterpriseClusterID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "cluster1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRedisenterpriseClusterID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RedisEnterpriseClusterId
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
			// missing cluster name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Cache/redisEnterprise/cluster1",
			Expected: &RedisEnterpriseClusterId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				Name:           "cluster1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CACHE/REDISENTERPRISE/CLUSTER1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RedisEnterpriseClusterID(v.Input)
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

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
