package parse

import "testing"

func TestLoadBalancerRuleIDParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *LoadBalancingRuleId
	}{
		{
			// load balancer id
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1",
			expected: nil,
		},
		{
			// lower-case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/loadbalancingrules/rule1",
			expected: nil,
		},
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/loadBalancingRules/rule1",
			expected: &LoadBalancingRuleId{
				ResourceGroup:    "group1",
				LoadBalancerName: "lb1",
				Name:             "rule1",
			},
		},
		{
			// title case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/Loadbalancers/lb1/Loadbalancingrules/rule1",
			expected: nil,
		},
		{
			// pascal case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/LoadBalancers/lb1/LoadBalancingRules/rule1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := LoadBalancingRuleID(test.input)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		if actual.ResourceGroup != test.expected.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but was %q", test.expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.LoadBalancerName != test.expected.LoadBalancerName {
			t.Fatalf("Expected LoadBalancerName to be %q but was %q", test.expected.LoadBalancerName, actual.LoadBalancerName)
		}

		if actual.Name != test.expected.Name {
			t.Fatalf("Expected name to be %q but was %q", test.expected.Name, actual.Name)
		}
	}
}
