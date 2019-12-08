package web

import "testing"

func TestValidateAppServiceCertificateID(t *testing.T) {
	cases := []struct {
		ID    string
		Valid bool
	}{
		{
			ID:    "",
			Valid: false,
		},
		{
			ID:    "nonsense",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/certificates/cert1",
			Valid: true,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/Certificates/cert1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %q", tc.ID)
		_, errors := ValidateAppServiceCertificateID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
