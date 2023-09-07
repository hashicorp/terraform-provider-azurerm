package parse

import "testing"

func TestNewOrganizationID(t *testing.T) {
	id, err := NewOrganizationID("subdomain", "domainSuffix", "organizationId")
	if err != nil {
		t.Fatalf("Got error for New Organization ID: %+v", err)
	}
	actual := id.ID()
	expected := "https://subdomain.domainSuffix/api/organizations/organizationId"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseNestedItemID(t *testing.T) {
	cases := []struct {
		Input       string
		Expected    OrganizationId
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https",
			ExpectError: true,
		},
		{
			Input:       "https://",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix/",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix/api",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix/api/organizations",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix/api/domains/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: true,
		},
		{
			Input:       "https://subdomain.domainSuffix/api/organizations/fdf067c93bbb4b22bff4d8b7a9a56217",
			ExpectError: false,
			Expected: OrganizationId{
				OrganizationId: "fdf067c93bbb4b22bff4d8b7a9a56217",
				DomainSuffix:   "domainSuffix",
				SubDomain:      "subdomain",
			},
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing %q", tc.Input)

		orgId, err := ParseOrganizationID(tc.Input, "subdomain", "domainSuffix")
		if err != nil {
			if tc.ExpectError {
				t.Logf("[DEBUG]   --> [Received Expected Error]: %+v", err)
				continue
			}

			t.Fatalf("Got error for ID '%s': %+v", tc.Input, err)
		}

		if orgId == nil {
			t.Fatalf("Expected a SecretID to be parsed for ID '%s', got nil.", tc.Input)
		}

		if tc.Expected.SubDomain != orgId.SubDomain {
			t.Fatalf("Expected 'SubDomain' to be '%s', got '%s' for ID '%s'", tc.Expected.SubDomain, orgId.SubDomain, tc.Input)
		}

		if tc.Expected.DomainSuffix != orgId.DomainSuffix {
			t.Fatalf("Expected 'DomainSuffix' to be '%s', got '%s' for ID '%s'", tc.Expected.DomainSuffix, orgId.DomainSuffix, tc.Input)
		}

		if tc.Expected.OrganizationId != orgId.OrganizationId {
			t.Fatalf("Expected 'OrganizationId' to be '%s', got '%s' for ID '%s'", tc.Expected.OrganizationId, orgId.OrganizationId, tc.Input)
		}

		if tc.Input != orgId.ID() {
			t.Fatalf("Expected 'ID()' to be '%s', got '%s'", tc.Input, orgId.ID())
		}
		t.Logf("[DEBUG]   --> [Valid Value]: %+v", tc.Input)
	}
}
