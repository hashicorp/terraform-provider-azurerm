package authentication

import "testing"

func TestServicePrincipalClientSecretAuth_builder(t *testing.T) {
	builder := Builder{
		ClientID:       "some-client-id",
		ClientSecret:   "some-client-secret",
		Environment:    "some-environment",
		SubscriptionID: "some-subscription-id",
		TenantID:       "some-tenant-id",
	}
	config := newServicePrincipalClientSecretAuth(builder)
	servicePrincipal := config.(servicePrincipalClientSecretAuth)

	if builder.ClientID != servicePrincipal.clientId {
		t.Fatalf("Expected Client ID to be %q but got %q", builder.ClientID, servicePrincipal.clientId)
	}

	if builder.ClientSecret != servicePrincipal.clientSecret {
		t.Fatalf("Expected Client Secret to be %q but got %q", builder.ClientSecret, servicePrincipal.clientSecret)
	}

	if builder.Environment != servicePrincipal.environment {
		t.Fatalf("Expected Environment to be %q but got %q", builder.Environment, servicePrincipal.environment)
	}

	if builder.SubscriptionID != servicePrincipal.subscriptionId {
		t.Fatalf("Expected Subscription ID to be %q but got %q", builder.SubscriptionID, servicePrincipal.subscriptionId)
	}

	if builder.TenantID != servicePrincipal.tenantId {
		t.Fatalf("Expected Tenant ID to be %q but got %q", builder.TenantID, servicePrincipal.tenantId)
	}
}

func TestServicePrincipalClientSecretAuth_validate(t *testing.T) {
	cases := []struct {
		Description string
		Config      servicePrincipalClientSecretAuth
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      servicePrincipalClientSecretAuth{},
			ExpectError: true,
		},
		{
			Description: "Missing Client ID",
			Config: servicePrincipalClientSecretAuth{
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientSecret:   "Does Hammer Time have Daylight Savings Time?",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: servicePrincipalClientSecretAuth{
				clientId:     "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				clientSecret: "Does Hammer Time have Daylight Savings Time?",
				tenantId:     "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				environment:  "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Client Secret",
			Config: servicePrincipalClientSecretAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: servicePrincipalClientSecretAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientSecret:   "Does Hammer Time have Daylight Savings Time?",
				environment:    "public",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Environment",
			Config: servicePrincipalClientSecretAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientSecret:   "Does Hammer Time have Daylight Savings Time?",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration",
			Config: servicePrincipalClientSecretAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientSecret:   "Does Hammer Time have Daylight Savings Time?",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
				environment:    "public",
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
