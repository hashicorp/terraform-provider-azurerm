package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/adal"
)

func TestAzureCLIParsingAuth_isApplicable(t *testing.T) {
	cases := []struct {
		Description string
		Builder     Builder
		Valid       bool
	}{
		{
			Description: "Empty Configuration",
			Builder:     Builder{},
			Valid:       false,
		},
		{
			Description: "Feature Toggled off",
			Builder: Builder{
				SupportsAzureCliCloudShellParsing: false,
			},
			Valid: false,
		},
		{
			Description: "Feature Toggled on",
			Builder: Builder{
				SupportsAzureCliCloudShellParsing: true,
			},
			Valid: true,
		},
	}

	for _, v := range cases {
		applicable := azureCliParsingAuth{}.isApplicable(v.Builder)
		if v.Valid != applicable {
			t.Fatalf("Expected %q to be %t but got %t", v.Description, v.Valid, applicable)
		}
	}
}

func TestAzureCLIParsingAuth_populateConfig(t *testing.T) {
	config := &Config{}
	auth := azureCliParsingAuth{
		profile: &azureCLIProfile{
			clientId:       "some-subscription-id",
			environment:    "dimension-c137",
			subscriptionId: "some-subscription-id",
			tenantId:       "some-tenant-id",
		},
	}

	err := auth.populateConfig(config)
	if err != nil {
		t.Fatalf("Error populating config: %s", err)
	}

	if auth.profile.clientId != config.ClientID {
		t.Fatalf("Expected Client ID to be %q but got %q", auth.profile.tenantId, config.TenantID)
	}

	if auth.profile.environment != config.Environment {
		t.Fatalf("Expected Environment to be %q but got %q", auth.profile.tenantId, config.TenantID)
	}

	if auth.profile.subscriptionId != config.SubscriptionID {
		t.Fatalf("Expected Subscription ID to be %q but got %q", auth.profile.tenantId, config.TenantID)
	}

	if auth.profile.tenantId != config.TenantID {
		t.Fatalf("Expected Tenant ID to be %q but got %q", auth.profile.tenantId, config.TenantID)
	}
}

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
