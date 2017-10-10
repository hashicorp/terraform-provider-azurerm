package azurerm

import "testing"

func TestResourceAzureRMNetworkSecurityRuleProtocol_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "tcp",
			ErrCount: 0,
		},
		{
			Value:    "TCP",
			ErrCount: 0,
		},
		{
			Value:    "*",
			ErrCount: 0,
		},
		{
			Value:    "Udp",
			ErrCount: 0,
		},
		{
			Value:    "Tcp",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateNetworkSecurityRuleProtocol(tc.Value, "azurerm_network_security_rule")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Network Security Rule protocol to trigger a validation error")
		}
	}
}
