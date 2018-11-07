package authentication

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestServicePrincipalClientCertAuth_builder(t *testing.T) {
	builder := Builder{
		ClientID:           "some-client-id",
		ClientCertPath:     "some-client-cert-path",
		ClientCertPassword: "some-password",
		Environment:        "some-environment",
		SubscriptionID:     "some-subscription-id",
		TenantID:           "some-tenant-id",
	}
	config, err := servicePrincipalClientCertificateAuth{}.build(builder)
	if err != nil {
		t.Fatalf("Error building client cert auth: %s", err)
	}

	servicePrincipal := config.(servicePrincipalClientCertificateAuth)

	if builder.ClientID != servicePrincipal.clientId {
		t.Fatalf("Expected Client ID to be %q but got %q", builder.ClientID, servicePrincipal.clientId)
	}

	if builder.ClientCertPath != servicePrincipal.clientCertPath {
		t.Fatalf("Expected Client Certificate Path to be %q but got %q", builder.ClientCertPath, servicePrincipal.clientCertPath)
	}

	if builder.ClientCertPassword != servicePrincipal.clientCertPassword {
		t.Fatalf("Expected Client Certificate Password to be %q but got %q", builder.ClientCertPassword, servicePrincipal.clientCertPassword)
	}

	if builder.SubscriptionID != servicePrincipal.subscriptionId {
		t.Fatalf("Expected Subscription ID to be %q but got %q", builder.SubscriptionID, servicePrincipal.subscriptionId)
	}

	if builder.TenantID != servicePrincipal.tenantId {
		t.Fatalf("Expected Tenant ID to be %q but got %q", builder.TenantID, servicePrincipal.tenantId)
	}
}

func TestServicePrincipalClientCertAuth_isApplicable(t *testing.T) {
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
				SupportsClientCertAuth: false,
			},
			Valid: false,
		},
		{
			Description: "Feature Toggled on but no cert specified",
			Builder: Builder{
				SupportsClientCertAuth: true,
			},
			Valid: false,
		},
		{
			Description: "Cert specified but feature toggled off",
			Builder: Builder{
				ClientCertPath: "./path/to/file",
			},
			Valid: false,
		},
		{
			Description: "Valid configuration",
			Builder: Builder{
				SupportsClientCertAuth: true,
				ClientCertPath:         "./path/to/file",
			},
			Valid: true,
		},
	}

	for _, v := range cases {
		applicable := servicePrincipalClientCertificateAuth{}.isApplicable(v.Builder)
		if v.Valid != applicable {
			t.Fatalf("Expected %q to be %t but got %t", v.Description, v.Valid, applicable)
		}
	}
}

func TestServicePrincipalClientCertAuth_populateConfig(t *testing.T) {
	config := &Config{}
	err := servicePrincipalClientCertificateAuth{}.populateConfig(config)
	if err != nil {
		t.Fatalf("Error populating config: %s", err)
	}

	if !config.AuthenticatedAsAServicePrincipal {
		t.Fatalf("Expected `AuthenticatedAsAServicePrincipal` to be true but it wasn't")
	}
}

func TestServicePrincipalClientCertAuth_validate(t *testing.T) {
	data := []byte("client-cert-auth")
	filePath := "./example.pfx"
	err := ioutil.WriteFile(filePath, data, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	cases := []struct {
		Description string
		Config      servicePrincipalClientCertificateAuth
		ExpectError bool
	}{
		{
			Description: "Empty Configuration",
			Config:      servicePrincipalClientCertificateAuth{},
			ExpectError: true,
		},
		{
			Description: "Missing Client ID",
			Config: servicePrincipalClientCertificateAuth{
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath: filePath,
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Subscription ID",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				clientCertPath: filePath,
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Client Certificate Path",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Missing Tenant ID",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath: filePath,
			},
			ExpectError: true,
		},
		{
			Description: "File isn't a pfx",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath: "not-valid.txt",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "File does not exist",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath: "does-not-exist.pfx",
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: true,
		},
		{
			Description: "Valid Configuration (Basic)",
			Config: servicePrincipalClientCertificateAuth{
				clientId:       "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId: "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath: filePath,
				tenantId:       "9834f8d0-24b3-41b7-8b8d-c611c461a129",
			},
			ExpectError: false,
		},
		{
			Description: "Valid Configuration (Complete)",
			Config: servicePrincipalClientCertificateAuth{
				clientId:           "62e73395-5017-43b6-8ebf-d6c30a514cf1",
				subscriptionId:     "8e8b5e02-5c13-4822-b7dc-4232afb7e8c2",
				clientCertPath:     filePath,
				clientCertPassword: "Password1234!",
				tenantId:           "9834f8d0-24b3-41b7-8b8d-c611c461a129",
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
