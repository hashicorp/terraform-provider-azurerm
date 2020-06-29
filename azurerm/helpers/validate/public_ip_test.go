package validate

import (
	"testing"
)

func TestPublicIpDomainNameLabel(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "tEsting123",
			ErrCount: 1,
		},
		{
			Value:    "testing123!",
			ErrCount: 1,
		},
		{
			Value:    "testing123-",
			ErrCount: 1,
		},
		{
			Value:    "k2345678-1-2345678-2-2345678-3-2345678-4-2345678-5-2345678-6-23",
			ErrCount: 0,
		},
		{
			Value:    "k2345678-1-2345678-2-2345678-3-2345678-4-2345678-5-2345678-6-234",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := PublicIpDomainNameLabel(tc.Value, "azurerm_public_ip")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Public IP Domain Name Label %s to trigger a validation error", tc.Value)
		}
	}
}
