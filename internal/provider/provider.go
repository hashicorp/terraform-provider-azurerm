// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func AzureProvider() *schema.Provider {
	return azureProvider(false)
}

func TestAzureProvider() *schema.Provider {
	return azureProvider(true)
}

func ValidatePartnerID(i interface{}, k string) ([]string, []error) {
	// ValidatePartnerID checks if partner_id is any of the following:
	//  * a valid UUID - will add "pid-" prefix to the ID if it is not already present
	//  * a valid UUID prefixed with "pid-"
	//  * a valid UUID prefixed with "pid-" and suffixed with "-partnercenter"

	debugLog := func(f string, v ...interface{}) {
		if os.Getenv("TF_LOG") == "" {
			return
		}

		if os.Getenv("TF_ACC") != "" {
			return
		}

		log.Printf(f, v...)
	}

	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if v == "" {
		return nil, nil
	}

	// Check for pid=<guid>-partnercenter format
	if strings.HasPrefix(v, "pid-") && strings.HasSuffix(v, "-partnercenter") {
		g := strings.TrimPrefix(v, "pid-")
		g = strings.TrimSuffix(g, "-partnercenter")

		if _, err := validation.IsUUID(g, ""); err != nil {
			return nil, []error{fmt.Errorf("expected %q to contain a valid UUID", v)}
		}

		debugLog("[DEBUG] %q partner_id matches pid-<GUID>-partnercenter...", v)
		return nil, nil
	}

	// Check for pid=<guid> (without the -partnercenter suffix)
	if strings.HasPrefix(v, "pid-") && !strings.HasSuffix(v, "-partnercenter") {
		g := strings.TrimPrefix(v, "pid-")

		if _, err := validation.IsUUID(g, ""); err != nil {
			return nil, []error{fmt.Errorf("expected %q to be a valid UUID", k)}
		}

		debugLog("[DEBUG] %q partner_id matches pid-<GUID>...", v)
		return nil, nil
	}

	// Check for straight UUID
	if _, err := validation.IsUUID(v, ""); err != nil {
		return nil, []error{fmt.Errorf("expected %q to be a valid UUID", k)}
	} else {
		debugLog("[DEBUG] %q partner_id is an un-prefixed UUID...", v)
		return nil, nil
	}
}

func azureProvider(supportLegacyTestSuite bool) *schema.Provider {
	// avoids this showing up in test output
	debugLog := func(f string, v ...interface{}) {
		if os.Getenv("TF_LOG") == "" {
			return
		}

		if os.Getenv("TF_ACC") != "" {
			return
		}

		log.Printf(f, v...)
	}

	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)

	// first handle the typed services
	for _, service := range SupportedTypedServices() {
		debugLog("[DEBUG] Registering Data Sources for %q..", service.Name())
		for _, ds := range service.DataSources() {
			key := ds.ResourceType()
			if existing := dataSources[key]; existing != nil {
				panic(fmt.Sprintf("An existing Data Source exists for %q", key))
			}

			wrapper := sdk.NewDataSourceWrapper(ds)
			dataSource, err := wrapper.DataSource()
			if err != nil {
				panic(fmt.Errorf("creating Wrapper for Data Source %q: %+v", key, err))
			}

			dataSources[key] = dataSource
		}

		debugLog("[DEBUG] Registering Resources for %q..", service.Name())
		for _, r := range service.Resources() {
			key := r.ResourceType()
			if existing := resources[key]; existing != nil {
				panic(fmt.Sprintf("An existing Resource exists for %q", key))
			}

			wrapper := sdk.NewResourceWrapper(r)
			resource, err := wrapper.Resource()
			if err != nil {
				panic(fmt.Errorf("creating Wrapper for Resource %q: %+v", key, err))
			}
			resources[key] = resource
		}
	}

	// then handle the untyped services
	for _, service := range SupportedUntypedServices() {
		debugLog("[DEBUG] Registering Data Sources for %q..", service.Name())
		for k, v := range service.SupportedDataSources() {
			if existing := dataSources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Data Source exists for %q", k))
			}

			dataSources[k] = v
		}

		debugLog("[DEBUG] Registering Resources for %q..", service.Name())
		for k, v := range service.SupportedResources() {
			if existing := resources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Resource exists for %q", k))
			}

			resources[k] = v
		}
	}

	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
				Description: "The Subscription ID which should be used.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used.",
			},

			"client_id_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID_FILE_PATH", ""),
				Description: "The path to a file containing the Client ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used.",
			},

			"auxiliary_tenant_ids": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 3,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment, and china. Defaults to public.",
			},

			"metadata_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			// Client Certificate specific fields
			"client_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", ""),
				Description: "Base64 encoded PKCS#12 certificate bundle to use when authenticating as a Service Principal using a Client Certificate",
			},

			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"client_secret_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET_FILE_PATH", ""),
				Description: "The path to a file containing the Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			// OIDC specifc fields
			"oidc_request_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, ""),
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, ""),
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN", ""),
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", ""),
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"use_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_OIDC", false),
				Description: "Allow OpenID Connect to be used for authentication",
			},

			// Managed Service Identity specific fields
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
				Description: "Allow Managed Service Identity to be used for Authentication.",
			},
			"msi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
				Description: "The path to a custom endpoint for Managed Service Identity - in most circumstances this should be detected automatically. ",
			},

			// Azure CLI specific fields
			"use_cli": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_CLI", true),
				Description: "Allow Azure CLI to be used for Authentication.",
			},

			// Managed Tracking GUID for User-agent
			"partner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(ValidatePartnerID, validation.StringIsEmpty),
				DefaultFunc:  schema.EnvDefaultFunc("ARM_PARTNER_ID", ""),
				Description:  "A GUID/UUID that is registered with Microsoft to facilitate partner resource usage attribution.",
			},

			"disable_correlation_request_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_DISABLE_CORRELATION_REQUEST_ID", false),
				Description: "This will disable the x-ms-correlation-request-id header.",
			},

			"disable_terraform_partner_id": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_DISABLE_TERRAFORM_PARTNER_ID", false),
				Description: "This will disable the Terraform Partner ID which is used if a custom `partner_id` isn't specified.",
			},

			"features": schemaFeatures(supportLegacyTestSuite),

			// Advanced feature flags
			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", false),
				Description: "Should the AzureRM Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
			},

			"storage_use_azuread": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_STORAGE_USE_AZUREAD", false),
				Description: "Should the AzureRM Provider use AzureAD to access the Storage Data Plane API's?",
			},
		},

		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var auxTenants []string
		if v, ok := d.Get("auxiliary_tenant_ids").([]interface{}); ok && len(v) > 0 {
			auxTenants = *utils.ExpandStringSlice(v)
		} else if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
			auxTenants = strings.Split(v, ";")
		}

		if len(auxTenants) > 3 {
			return nil, diag.Errorf("the provider only supports 3 auxiliary tenant IDs")
		}

		var clientCertificateData []byte
		if encodedCert := d.Get("client_certificate").(string); encodedCert != "" {
			var err error
			clientCertificateData, err = decodeCertificate(encodedCert)
			if err != nil {
				return nil, diag.FromErr(err)
			}
		}

		oidcToken, err := getOidcToken(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		clientSecret, err := getClientSecret(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		clientId, err := getClientId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		var (
			env *environments.Environment

			envName      = d.Get("environment").(string)
			metadataHost = d.Get("metadata_host").(string)
		)

		if metadataHost != "" {
			if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost), envName); err != nil {
				return nil, diag.FromErr(err)
			}
		} else if env, err = environments.FromName(envName); err != nil {
			return nil, diag.FromErr(err)
		}

		var (
			enableAzureCli        = d.Get("use_cli").(bool)
			enableManagedIdentity = d.Get("use_msi").(bool)
			enableOidc            = d.Get("use_oidc").(bool)
		)

		authConfig := &auth.Credentials{
			Environment:        *env,
			ClientID:           *clientId,
			TenantID:           d.Get("tenant_id").(string),
			AuxiliaryTenantIDs: auxTenants,

			ClientCertificateData:     clientCertificateData,
			ClientCertificatePath:     d.Get("client_certificate_path").(string),
			ClientCertificatePassword: d.Get("client_certificate_password").(string),
			ClientSecret:              *clientSecret,

			OIDCAssertionToken:          *oidcToken,
			GitHubOIDCTokenRequestURL:   d.Get("oidc_request_url").(string),
			GitHubOIDCTokenRequestToken: d.Get("oidc_request_token").(string),

			CustomManagedIdentityEndpoint: d.Get("msi_endpoint").(string),

			EnableAuthenticatingUsingClientCertificate: true,
			EnableAuthenticatingUsingClientSecret:      true,
			EnableAuthenticatingUsingAzureCLI:          enableAzureCli,
			EnableAuthenticatingUsingManagedIdentity:   enableManagedIdentity,
			EnableAuthenticationUsingOIDC:              enableOidc,
			EnableAuthenticationUsingGitHubOIDC:        enableOidc,
		}

		return buildClient(ctx, p, d, authConfig)
	}
}

func buildClient(ctx context.Context, p *schema.Provider, d *schema.ResourceData, authConfig *auth.Credentials) (*clients.Client, diag.Diagnostics) {
	skipProviderRegistration := d.Get("skip_provider_registration").(bool)

	clientBuilder := clients.ClientBuilder{
		AuthConfig:                  authConfig,
		DisableCorrelationRequestID: d.Get("disable_correlation_request_id").(bool),
		DisableTerraformPartnerID:   d.Get("disable_terraform_partner_id").(bool),
		Features:                    expandFeatures(d.Get("features").([]interface{})),
		MetadataHost:                d.Get("metadata_host").(string),
		PartnerID:                   d.Get("partner_id").(string),
		SkipProviderRegistration:    skipProviderRegistration,
		StorageUseAzureAD:           d.Get("storage_use_azuread").(bool),
		SubscriptionID:              d.Get("subscription_id").(string),
		TerraformVersion:            p.TerraformVersion,

		// this field is intentionally not exposed in the provider block, since it's only used for
		// platform level tracing
		CustomCorrelationRequestID: os.Getenv("ARM_CORRELATION_REQUEST_ID"),
	}

	//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golanci-lint
	stopCtx, ok := schema.StopContext(ctx) //nolint:staticcheck
	if !ok {
		stopCtx = ctx
	}

	client, err := clients.Build(stopCtx, clientBuilder)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.StopContext = stopCtx

	if !skipProviderRegistration {
		subscriptionId := commonids.NewSubscriptionID(client.Account.SubscriptionId)
		requiredResourceProviders := resourceproviders.Required()
		ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()

		if err := resourceproviders.EnsureRegistered(ctx2, client.Resource.ResourceProvidersClient, subscriptionId, requiredResourceProviders); err != nil {
			return nil, diag.Errorf(resourceProviderRegistrationErrorFmt, err)
		}
	}

	return client, nil
}

func decodeCertificate(clientCertificate string) ([]byte, error) {
	var pfx []byte
	if clientCertificate != "" {
		out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
		n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
		if err != nil {
			return pfx, fmt.Errorf("could not decode client certificate data: %v", err)
		}
		pfx = out[:n]
	}
	return pfx, nil
}

func getOidcToken(d *schema.ResourceData) (*string, error) {
	idToken := strings.TrimSpace(d.Get("oidc_token").(string))

	if path := d.Get("oidc_token_file_path").(string); path != "" {
		fileTokenRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match")
		}

		idToken = fileToken
	}

	return &idToken, nil
}

func getClientId(d *schema.ResourceData) (*string, error) {
	clientId := strings.TrimSpace(d.Get("client_id").(string))

	if path := d.Get("client_id_file_path").(string); path != "" {
		fileClientIdRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client ID from file %q: %v", path, err)
		}

		fileClientId := strings.TrimSpace(string(fileClientIdRaw))

		if clientId != "" && clientId != fileClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and supplied Client ID file contents - please either remove one or ensure they match")
		}

		clientId = fileClientId
	}

	return &clientId, nil
}

func getClientSecret(d *schema.ResourceData) (*string, error) {
	clientSecret := strings.TrimSpace(d.Get("client_secret").(string))

	if path := d.Get("client_secret_file_path").(string); path != "" {
		fileSecretRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client Secret from file %q: %v", path, err)
		}

		fileSecret := strings.TrimSpace(string(fileSecretRaw))

		if clientSecret != "" && clientSecret != fileSecret {
			return nil, fmt.Errorf("mismatch between supplied Client Secret and supplied Client Secret file contents - please either remove one or ensure they match")
		}

		clientSecret = fileSecret
	}

	return &clientSecret, nil
}

const resourceProviderRegistrationErrorFmt = `Error ensuring Resource Providers are registered.

Terraform automatically attempts to register the Resource Providers it supports to
ensure it's able to provision resources.

If you don't have permission to register Resource Providers you may wish to use the
"skip_provider_registration" flag in the Provider block to disable this functionality.

Please note that if you opt out of Resource Provider Registration and Terraform tries
to provision a resource from a Resource Provider which is unregistered, then the errors
may appear misleading - for example:

> API version 2019-XX-XX was not found for Microsoft.Foo

Could indicate either that the Resource Provider "Microsoft.Foo" requires registration,
but this could also indicate that this Azure Region doesn't support this API version.

More information on the "skip_provider_registration" flag can be found here:
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#skip_provider_registration

Original Error: %s`
