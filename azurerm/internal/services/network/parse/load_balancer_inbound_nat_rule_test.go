package parse

import "testing"

func TestLoadBalancerInboundNATRuleIDParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *LoadBalancerInboundNATRuleId
	}{
		{
			// load balancer id
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1",
			expected: nil,
		},
		{
			// lower-case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/inboundnatrules/rule1",
			expected: nil,
		},
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/loadBalancers/lb1/inboundNatRules/rule1",
			expected: &LoadBalancerInboundNATRuleId{
				ResourceGroup:      "group1",
				LoadBalancerName:   "lb1",
				InboundNatRuleName: "rule1",
			},
		},
		{
			// title case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/Loadbalancers/lb1/Inboundnatrules/rule1",
			expected: nil,
		},
		{
			// pascal case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/LoadBalancers/lb1/InboundNatRules/rule1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := LoadBalancerInboundNATRuleID(test.input)
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

		if actual.InboundNatRuleName != test.expected.InboundNatRuleName {
			t.Fatalf("Expected name to be %q but was %q", test.expected.InboundNatRuleName, actual.InboundNatRuleName)
		}
	}
}
