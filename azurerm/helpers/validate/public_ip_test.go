package validate

import "testing"

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
			Value:    acctest.RandString(80),
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := PublicIpDomainNameLabel(tc.Value, "azurerm_public_ip")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Azure RM Public IP Domain Name Label to trigger a validation error")
		}
	}
}
