package authentication

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/adal"
)

func TestAzureValidateBearerAuth(t *testing.T) {
	cases := []struct {
		Description string
		Config      Config
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      Config{},
			ExpectError: true,
		},
		{
			Description: "Missing Access Token",
			Config: Config{
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Client ID",
			Config: Config{
				AccessToken:    &adal.Token{},
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: Config{
				AccessToken: &adal.Token{},
				ClientID:    "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				TenantID:    "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: Config{
				AccessToken:    &adal.Token{},
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: Config{
				AccessToken:    &adal.Token{},
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: false,
		},
	}

	for _, v := range cases {
		err := v.Config.ValidateBearerAuth()

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}
	}
}

func TestAzureValidateServicePrincipal(t *testing.T) {
	cases := []struct {
		Description string
		Config      Config
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      Config{},
			ExpectError: true,
		},
		{
			Description: "Missing Client ID",
			Config: Config{
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				ClientSecret:   "Does Hammer Time have Daylight Savings Time?",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: Config{
				ClientID:     "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				ClientSecret: "Does Hammer Time have Daylight Savings Time?",
				TenantID:     "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:  "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Client Secret",
			Config: Config{
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: Config{
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				ClientSecret:   "Does Hammer Time have Daylight Savings Time?",
				Environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Environment",
			Config: Config{
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				ClientSecret:   "Does Hammer Time have Daylight Savings Time?",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: Config{
				ClientID:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				ClientSecret:   "Does Hammer Time have Daylight Savings Time?",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:    "public",
			},
			ExpectError: false,
		},
	}

	for _, v := range cases {
		err := v.Config.ValidateServicePrincipal()

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}
	}
}

func TestAzureValidateMsi(t *testing.T) {
	cases := []struct {
		Description string
		Config      Config
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      Config{},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: Config{
				MsiEndpoint: "http://localhost:50342/oauth2/token",
				TenantID:    "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment: "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: Config{
				MsiEndpoint:    "http://localhost:50342/oauth2/token",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				Environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Environment",
			Config: Config{
				MsiEndpoint:    "http://localhost:50342/oauth2/token",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing MSI Endpoint",
			Config: Config{
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: Config{
				MsiEndpoint:    "http://localhost:50342/oauth2/token",
				SubscriptionID: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				TenantID:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				Environment:    "public",
			},
			ExpectError: false,
		},
	}

	for _, v := range cases {
		err := v.Config.ValidateMsi()

		if v.ExpectError && err == nil {
			t.Fatalf("Expected an error for %q: didn't get one", v.Description)
		}

		if !v.ExpectError && err != nil {
			t.Fatalf("Expected there to be no error for %q - but got: %v", v.Description, err)
		}
	}
}
