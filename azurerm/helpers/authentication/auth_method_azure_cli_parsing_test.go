package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/adal"
)

func TestAzureCLIParsingAuth_validate(t *testing.T) {
	cases := []struct {
		Description string
		Config      azureCliParsingAuth
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      azureCliParsingAuth{},
			ExpectError: true,
		},
		{
			Description: "Missing Access Token",
			Config: azureCliParsingAuth{
				profile: &azureCLIProfile{
					clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
					subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
					tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				},
			},
			ExpectError: true,
		},
		{
			Description: "Missing Client ID",
			Config: azureCliParsingAuth{
				profile: &azureCLIProfile{
					accessToken:    &adal.Token{},
					subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
					tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				},
			},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: azureCliParsingAuth{
				profile: &azureCLIProfile{
					accessToken: &adal.Token{},
					clientId:    "62e73395-5017-43b6-8ebf-d6c30a514cf1",
					tenantId:    "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				},
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: azureCliParsingAuth{
				profile: &azureCLIProfile{
					accessToken:    &adal.Token{},
					clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
					subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				},
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: azureCliParsingAuth{
				profile: &azureCLIProfile{
					accessToken:    &adal.Token{},
					clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
					subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
					tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				},
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
