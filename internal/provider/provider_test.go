// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
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
			if (resource.Timeouts.Create == nil) != (resource.Create == nil && resource.CreateContext == nil) { //nolint:staticcheck
				t.Fatalf("Resource %q should define/not define the Create(Context) method and the Create Timeout at the same time", resourceName)
			}
			if (resource.Timeouts.Delete == nil) != (resource.Delete == nil && resource.DeleteContext == nil) { //nolint:staticcheck
				t.Fatalf("Resource %q should define/not define the Delete(Context) method and the Delete Timeout at the same time", resourceName)
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
			if (resource.Timeouts.Update == nil) != (resource.Update == nil && resource.UpdateContext == nil) { //nolint:staticcheck
				t.Fatalf("Resource %q should define/not define the Update(Context) method and the Update Timeout at the same time", resourceName)
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

func TestAccProvider_resourceProviders_legacy(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(nil)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := resourceproviders.Legacy()
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
}

// TODO: Remove this test in v5.0
func TestAccProvider_resourceProviders_deprecatedSkip(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()
	config := map[string]interface{}{
		"skip_provider_registration": "true",
	}

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(config)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := make(resourceproviders.ResourceProviders)
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
}

func TestAccProvider_resourceProviders_legacyWithAdditional(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()
	config := map[string]interface{}{
		"resource_providers_to_register": []interface{}{
			"Microsoft.ApiManagement",
			"Microsoft.ContainerService",
			"Microsoft.KeyVault",
			"Microsoft.Kubernetes",
		},
	}

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(config)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := resourceproviders.Legacy().Merge(resourceproviders.ResourceProviders{
		"Microsoft.ApiManagement":    {},
		"Microsoft.ContainerService": {},
		"Microsoft.KeyVault":         {},
		"Microsoft.Kubernetes":       {},
	})
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
}

func TestAccProvider_resourceProviders_core(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()
	config := map[string]interface{}{
		"resource_provider_registrations": "core",
	}

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(config)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := resourceproviders.Core()
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
}

func TestAccProvider_resourceProviders_coreWithAdditional(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()
	config := map[string]interface{}{
		"resource_provider_registrations": "core",
		"resource_providers_to_register": []interface{}{
			"Microsoft.ApiManagement",
			"Microsoft.KeyVault",
		},
	}

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(config)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := resourceproviders.Core().Merge(resourceproviders.ResourceProviders{
		"Microsoft.ApiManagement": {},
		"Microsoft.KeyVault":      {},
	})
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
}

func TestAccProvider_resourceProviders_explicit(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	logging.SetOutput(t)

	provider := TestAzureProvider()
	config := map[string]interface{}{
		"resource_provider_registrations": "none",
		"resource_providers_to_register": []interface{}{
			"Microsoft.Compute",
			"Microsoft.Network",
			"Microsoft.Storage",
		},
	}

	if diags := provider.Configure(ctx, terraform.NewResourceConfigRaw(config)); diags != nil && diags.HasError() {
		t.Fatalf("provider failed to configure: %v", diags)
	}

	expectedResourceProviders := resourceproviders.ResourceProviders{
		"Microsoft.Compute": {},
		"Microsoft.Network": {},
		"Microsoft.Storage": {},
	}
	registeredResourceProviders := provider.Meta().(*clients.Client).Account.RegisteredResourceProviders

	if !reflect.DeepEqual(registeredResourceProviders, expectedResourceProviders) {
		t.Fatalf("unexpected value for RegisteredResourceProviders: %#v", registeredResourceProviders)
	}
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
			AzureCliSubscriptionIDHint:        d.Get("subscription_id").(string),
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

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:               *env,
			TenantID:                  *tenantId,
			ClientID:                  *clientId,
			ClientCertificateData:     certData,
			ClientCertificatePath:     d.Get("client_certificate_path").(string),
			ClientCertificatePassword: d.Get("client_certificate_password").(string),
			EnableAuthenticatingUsingClientCertificate: true,
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

	// Support only Client Secret authentication
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

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                           *env,
			TenantID:                              *tenantId,
			ClientID:                              *clientId,
			ClientSecret:                          *clientSecret,
			EnableAuthenticatingUsingClientSecret: true,
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

	// Support only Client Secret authentication
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

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                           *env,
			TenantID:                              *tenantId,
			ClientID:                              *clientId,
			ClientSecret:                          *clientSecret,
			EnableAuthenticatingUsingClientSecret: true,
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

	// Support only OIDC authentication
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

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                   *env,
			TenantID:                      *tenantId,
			ClientID:                      *clientId,
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

	// Support only GitHub OIDC authentication
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

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                         *env,
			TenantID:                            *tenantId,
			ClientID:                            *clientId,
			OIDCTokenRequestToken:               d.Get("oidc_request_token").(string),
			OIDCTokenRequestURL:                 d.Get("oidc_request_url").(string),
			EnableAuthenticationUsingGitHubOIDC: true,
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

func TestAccProvider_adoOidcAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("SYSTEM_ACCESSTOKEN") == "" {
		t.Skip("SYSTEM_ACCESSTOKEN not set")
	}
	if os.Getenv("SYSTEM_OIDCREQUESTURI") == "" {
		t.Skip("SYSTEM_OIDCREQUESTURI not set")
	}
	if os.Getenv("ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID") == "" {
		t.Skip("ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only ADO OIDC authentication
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

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                              *env,
			TenantID:                                 *tenantId,
			ClientID:                                 *clientId,
			OIDCTokenRequestToken:                    d.Get("oidc_request_token").(string),
			OIDCTokenRequestURL:                      d.Get("oidc_request_url").(string),
			ADOPipelineServiceConnectionID:           d.Get("ado_pipeline_service_connection_id").(string),
			EnableAuthenticationUsingADOPipelineOIDC: true,
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

func TestAccProvider_aksWorkloadIdentityAuth(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set")
	}
	if os.Getenv("AZURE_CLIENT_ID") == "" {
		t.Skip("AZURE_CLIENT_ID not set")
	}
	if os.Getenv("AZURE_TENANT_ID") == "" {
		t.Skip("AZURE_TENANT_ID not set")
	}
	if os.Getenv("AZURE_FEDERATED_TOKEN_FILE") == "" {
		t.Skip("AZURE_FEDERATED_TOKEN_FILE not set")
	}

	logging.SetOutput(t)

	provider := TestAzureProvider()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Support only AKS Workload Identity authentication
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

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		authConfig := &auth.Credentials{
			Environment:                   *env,
			TenantID:                      *tenantId,
			ClientID:                      *clientId,
			OIDCAssertionToken:            *oidcToken,
			EnableAuthenticationUsingOIDC: true,
		}

		return buildClient(ctx, provider, d, authConfig)
	}

	// Ensure we enable AKS Workload Identity else the configuration will not be detected
	conf := map[string]interface{}{"use_aks_workload_identity": true}
	d := provider.Configure(ctx, terraform.NewResourceConfigRaw(conf))
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
