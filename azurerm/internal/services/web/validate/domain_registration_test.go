package validate

import "testing"

func TestValidateDomainRegistrationID(t *testing.T) {
	cases := []struct {
		ID    string
		Valid bool
	}{
		{
			ID:    "",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.DomainRegistration/",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup1/providers/Microsoft.DomainRegistration/domains/testDomain1",
			Valid: true,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/testGroup1/providers/Microsoft.DomainRegistration/Domains/testDomain1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.ID)
		_, errors := DomainRegistrationID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}

func TestValidateDomainRegistrationName(t *testing.T) {
	cases := []struct {
		Name  string
		Valid bool
	}{
		{
			Name:  "",
			Valid: false,
		},
		{
			Name:  "-testDomain",
			Valid: false,
		},
		{
			Name:  "testDomain-",
			Valid: false,
		},
		{
			Name:  "testDomain.",
			Valid: false,
		},
		{
			Name:  "abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz12",
			Valid: false,
		},
		{
			Name:  "1testDomain",
			Valid: true,
		},
		{
			Name:  "testdomain1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Name)
		_, errors := DomainRegistrationName(tc.Name, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
