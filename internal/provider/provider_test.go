// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func TestProvider(t *testing.T) {
	if err := TestAzureProvider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestDataSourcesSupportCustomTimeouts(t *testing.T) {
	provider := TestAzureProvider()
	for dataSourceName, dataSource := range provider.DataSourcesMap {
		t.Run(fmt.Sprintf("DataSource/%s", dataSourceName), func(t *testing.T) {
			t.Logf("[DEBUG] Testing Data Source %q..", dataSourceName)

			if dataSource.Timeouts == nil {
				t.Fatalf("Data Source %q has no timeouts block defined!", dataSourceName)
			}

			// Azure's bespoke enough we want to be explicit about the timeouts for each value
			if dataSource.Timeouts.Default != nil {
				t.Fatalf("Data Source %q defines a Default timeout when it shouldn't!", dataSourceName)
			}

			// Data Sources must have a Read
			if dataSource.Timeouts.Read == nil {
				t.Fatalf("Data Source %q doesn't define a Read timeout", dataSourceName)
			}

			// but shouldn't have anything else
			if dataSource.Timeouts.Create != nil {
				t.Fatalf("Data Source %q defines a Create timeout when it shouldn't!", dataSourceName)
			}

			if dataSource.Timeouts.Update != nil {
				t.Fatalf("Data Source %q defines a Update timeout when it shouldn't!", dataSourceName)
			}

			if dataSource.Timeouts.Delete != nil {
				t.Fatalf("Data Source %q defines a Delete timeout when it shouldn't!", dataSourceName)
			}
		})
	}
}

func TestResourcesSupportCustomTimeouts(t *testing.T) {
	provider := TestAzureProvider()
	for resourceName, resource := range provider.ResourcesMap {
		t.Run(fmt.Sprintf("Resource/%s", resourceName), func(t *testing.T) {
			t.Logf("[DEBUG] Testing Resource %q..", resourceName)

			if resource.Timeouts == nil {
				t.Fatalf("Resource %q has no timeouts block defined!", resourceName)
			}

			// Azure's bespoke enough we want to be explicit about the timeouts for each value
			if resource.Timeouts.Default != nil {
				t.Fatalf("Resource %q defines a Default timeout when it shouldn't!", resourceName)
			}

			// every Resource has to have a Create, Read & Destroy timeout

			//lint:ignore SA1019 SDKv2 migration  - staticcheck's own linter directives are currently being ignored under golanci-lint
			if resource.Timeouts.Create == nil && resource.Create != nil { //nolint:staticcheck
				t.Fatalf("Resource %q defines a Create method but no Create Timeout", resourceName)
			}
			if resource.Timeouts.Delete == nil && resource.Delete != nil { //nolint:staticcheck
				t.Fatalf("Resource %q defines a Delete method but no Delete Timeout", resourceName)
			}
			if resource.Timeouts.Read == nil {
				t.Fatalf("Resource %q doesn't define a Read timeout", resourceName)
			} else if *resource.Timeouts.Read > 5*time.Minute {
				exceptionResources := map[string]bool{
					// The key vault item resources have longer read timeout for mitigating issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/11059.
					"azurerm_key_vault_key":         true,
					"azurerm_key_vault_secret":      true,
					"azurerm_key_vault_certificate": true,
				}
				if !exceptionResources[resourceName] {
					t.Fatalf("Read timeouts shouldn't be more than 5 minutes, this indicates a bug which needs to be fixed")
				}
			}

			// Optional
			if resource.Timeouts.Update == nil && resource.Update != nil { //nolint:staticcheck
				t.Fatalf("Resource %q defines a Update method but no Update Timeout", resourceName)
			}
		})
	}
}

func TestProvider_impl(t *testing.T) {
	_ = AzureProvider()
}

func TestProvider_counts(t *testing.T) {
	// @tombuildsstuff: this is less a unit test and more a useful placeholder tbh
	provider := TestAzureProvider()
	log.Printf("Data Sources: %d", len(provider.DataSourcesMap))
	log.Printf("Resources:    %d", len(provider.ResourcesMap))
	log.Printf("-----------------")
	log.Printf("Total:        %d", len(provider.ResourcesMap)+len(provider.DataSourcesMap))
}

func TestAccProvider_cliAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Azure CLI authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		authConfig := &auth.Credentials{
			Environment:                       *env,
			EnableAuthenticatingUsingAzureCLI: true,
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_clientCertificateAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("ARM_CLIENT_ID") == "" {
		t.Skip("ARM_CLIENT_ID not set")
	}
	if os.Getenv("ARM_CLIENT_CERTIFICATE") == "" && os.Getenv("ARM_CLIENT_CERTIFICATE_PATH") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE or ARM_CLIENT_CERTIFICATE_PATH not set")
	}
	if os.Getenv("ARM_CLIENT_CERTIFICATE_PASSWORD") == "" {
		t.Skip("ARM_CLIENT_CERTIFICATE_PASSWORD not set")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Client Certificate authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		var certData []byte
		if encodedCert := d.Get("client_certificate").(string); encodedCert != "" {
			var err error
			certData, err = decodeCertificate(encodedCert)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		authConfig := &auth.Credentials{
			Environment: *env,
			TenantID:    d.Get("tenant_id").(string),
			ClientID:    d.Get("client_id").(string),
			EnableAuthenticatingUsingClientCertificate: true,
			ClientCertificateData:                      certData,
			ClientCertificatePath:                      d.Get("client_certificate_path").(string),
			ClientCertificatePassword:                  d.Get("client_certificate_password").(string),
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_clientSecretAuth(t *testing.T) {
	t.Run("fromEnvironment", testAccProvider_clientSecretAuthFromEnvironment)
	t.Run("fromFiles", testAccProvider_clientSecretAuthFromFiles)
}

func testAccProvider_clientSecretAuthFromEnvironment(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("ARM_CLIENT_ID") == "" {
		t.Skip("ARM_CLIENT_ID not set")
	}
	if os.Getenv("ARM_CLIENT_SECRET") == "" {
		t.Skip("ARM_CLIENT_SECRET not set")
	}

	// Ensure we are running using the expected env-vars
	// t.SetEnv does automatic cleanup / resets the values after the test
	t.Setenv("ARM_CLIENT_ID_FILE_PATH", "")
	t.Setenv("ARM_CLIENT_SECRET_FILE_PATH", "")

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Client Certificate authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		clientSecret, err := getClientSecret(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                           *env,
			TenantID:                              d.Get("tenant_id").(string),
			ClientID:                              *clientId,
			EnableAuthenticatingUsingClientSecret: true,
			ClientSecret:                          *clientSecret,
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func testAccProvider_clientSecretAuthFromFiles(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("ARM_CLIENT_ID_FILE_PATH") == "" {
		t.Skip("ARM_CLIENT_ID_FILE_PATH not set")
	}
	if os.Getenv("ARM_CLIENT_SECRET_FILE_PATH") == "" {
		t.Skip("ARM_CLIENT_SECRET_FILE_PATH not set")
	}

	// Ensure we are running using the expected env-vars
	// t.SetEnv does automatic cleanup / resets the values after the test
	t.Setenv("ARM_CLIENT_ID", "")
	t.Setenv("ARM_CLIENT_SECRET", "")

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Client Certificate authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		clientSecret, err := getClientSecret(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                           *env,
			TenantID:                              d.Get("tenant_id").(string),
			ClientID:                              *clientId,
			EnableAuthenticatingUsingClientSecret: true,
			ClientSecret:                          *clientSecret,
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_genericOidcAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("ARM_OIDC_TOKEN") == "" && os.Getenv("ARM_OIDC_TOKEN_FILE_PATH") == "" {
		t.Skip("ARM_OIDC_TOKEN or ARM_OIDC_TOKEN_FILE_PATH not set")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Client Certificate authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		oidcToken, err := getOidcToken(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                   *env,
			TenantID:                      d.Get("tenant_id").(string),
			ClientID:                      d.Get("client_id").(string),
			EnableAuthenticationUsingOIDC: true,
			OIDCAssertionToken:            *oidcToken,
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func TestAccProvider_githubOidcAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN") == "" {
		t.Skip("ACTIONS_ID_TOKEN_REQUEST_TOKEN not set")
	}
	if os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL") == "" {
		t.Skip("ACTIONS_ID_TOKEN_REQUEST_URL not set")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only Client Certificate authentication
	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		envName := d.Get("environment").(string)
		env, err := environments.FromName(envName)
		if err != nil {
			t.Fatalf("configuring environment %q: %v", envName, err)
		}

		authConfig := &auth.Credentials{
			Environment:                         *env,
			TenantID:                            d.Get("tenant_id").(string),
			ClientID:                            d.Get("client_id").(string),
			EnableAuthenticationUsingGitHubOIDC: true,
			GitHubOIDCTokenRequestToken:         d.Get("oidc_request_token").(string),
			GitHubOIDCTokenRequestURL:           d.Get("oidc_request_url").(string),
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d != nil && d.HasError() {
		t.Fatalf("err: %+v", d)
	}

	if errs := testCheckProvider(provider); len(errs) > 0 {
		for _, err := range errs {
			t.Error(err)
		}
	}
}

func testCheckProvider(provider *schema.Provider) (errs []error) {
	client := provider.Meta().(*clients.Client)

	if endpoint, ok := client.Account.Environment.MicrosoftGraph.Endpoint(); !ok || *endpoint == "" {
		errs = append(errs, fmt.Errorf("client.Account.Environment returned blank endpoint for Microsoft Graph"))
	}

	if endpoint, ok := client.Account.Environment.ResourceManager.Endpoint(); !ok || *endpoint == "" {
		errs = append(errs, fmt.Errorf("client.Account.Environment returned blank endpoint for Microsoft Graph"))
	}

	if client.Account.ClientId == "" {
		errs = append(errs, fmt.Errorf("client.Account.ClientId was empty"))
	}

	if client.Account.ObjectId == "" {
		errs = append(errs, fmt.Errorf("client.Account.ObjectId was empty"))
	}

	if client.Account.SubscriptionId == "" {
		errs = append(errs, fmt.Errorf("client.Account.SubscriptionId was empty"))
	}

	if client.Account.TenantId == "" {
		errs = append(errs, fmt.Errorf("client.Account.TenantId was empty"))
	}

	return //nolint:nakedret
}
