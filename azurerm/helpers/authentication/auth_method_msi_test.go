package authentication

import "testing"

func TestManagedServiceIdentity_builder(t *testing.T) {
	builder := Builder{
		MsiEndpoint: "https://hello-world",
	}

	method, err := newManagedServiceIdentityAuth(builder)
	if err != nil {
		t.Fatalf("Error building MSI Identity Auth: %+v", err)
	}

	authMethod := method.(managedServiceIdentityAuth)
	if builder.MsiEndpoint != authMethod.endpoint {
		t.Fatalf("Expected MSI Endpoint to be %q but got %q", builder.MsiEndpoint, authMethod.endpoint)
	}
}

func TestManagedServiceIdentity_validate(t *testing.T) {
	cases := []struct {
		Description string
		Config      managedServiceIdentityAuth
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      managedServiceIdentityAuth{},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: managedServiceIdentityAuth{
				endpoint: "https://some-location",
			},
			ExpectError: false,
		},
	}

	for _, v := range cases {
		err := v.Config.validate()

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}
	}
}
