// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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

		logEntry("[DEBUG] %q partner_id matches pid-<GUID>-partnercenter...", v)
		return nil, nil
	}

	// Check for pid=<guid> (without the -partnercenter suffix)
	if strings.HasPrefix(v, "pid-") && !strings.HasSuffix(v, "-partnercenter") {
		g := strings.TrimPrefix(v, "pid-")

		if _, err := validation.IsUUID(g, ""); err != nil {
			return nil, []error{fmt.Errorf("expected %q to be a valid UUID", k)}
		}

		logEntry("[DEBUG] %q partner_id matches pid-<GUID>...", v)
		return nil, nil
	}

	// Check for straight UUID
	if _, err := validation.IsUUID(v, ""); err != nil {
		return nil, []error{fmt.Errorf("expected %q to be a valid UUID", k)}
	} else {
		logEntry("[DEBUG] %q partner_id is an un-prefixed UUID...", v)
		return nil, nil
	}
}

func azureProvider(supportLegacyTestSuite bool) *schema.Provider {
	dataSources := make(map[string]*schema.Resource)
	resources := make(map[string]*schema.Resource)

	// first handle the typed services
	for _, service := range SupportedTypedServices() {
		logEntry("[DEBUG] Registering Data Sources for %q..", service.Name())
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

		logEntry("[DEBUG] Registering Resources for %q..", service.Name())
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
		logEntry("[DEBUG] Registering Data Sources for %q..", service.Name())
		for k, v := range service.SupportedDataSources() {
			if existing := dataSources[k]; existing != nil {
				panic(fmt.Sprintf("An existing Data Source exists for %q", k))
			}

			dataSources[k] = v
		}

		logEntry("[DEBUG] Registering Resources for %q..", service.Name())
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
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", nil),
				Description: "The Subscription ID which should be used.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
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
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
				Description: "The Cloud Environment which should be used. Possible values are public, usgovernment, and china. Defaults to public. Not used and should not be specified when `metadata_host` is specified.",
			},

			"metadata_host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", nil),
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
				Description: "The Client ID which should be used.",
			},

			"client_id_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID_FILE_PATH", nil),
				Description: "The path to a file containing the Client ID which should be used.",
			},

			// Client Certificate specific fields
			"client_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", nil),
				Description: "Base64 encoded PKCS#12 certificate bundle to use when authenticating as a Service Principal using a Client Certificate",
			},

			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", nil),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate.",
			},

			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", nil),
				Description: "The password associated with the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
				Description: "The Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			"client_secret_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET_FILE_PATH", nil),
				Description: "The path to a file containing the Client Secret which should be used. For use When authenticating as a Service Principal using a Client Secret.",
			},

			// OIDC specific fields
			"use_oidc": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_OIDC", false),
				Description: "Allow OpenID Connect to be used for authentication",
			},

			"oidc_request_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, nil),
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, nil),
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN", nil),
				Description: "The OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},

			"oidc_token_file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", nil),
				Description: "The path to a file containing an OIDC ID token for use when authenticating as a Service Principal using OpenID Connect.",
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
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", nil),
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

			// Azure AKS Workload Identity fields
			"use_aks_workload_identity": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_AKS_WORKLOAD_IDENTITY", false),
				Description: "Allow Azure AKS Workload Identity to be used for Authentication.",
			},

			// Managed Tracking GUID for User-agent
			"partner_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.Any(ValidatePartnerID, validation.StringIsEmpty),
				DefaultFunc:  schema.EnvDefaultFunc("ARM_PARTNER_ID", nil),
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
			"resource_provider_registrations": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_RESOURCE_PROVIDER_REGISTRATIONS", resourceproviders.ProviderRegistrationsLegacy),
				Description: "The set of Resource Providers which should be automatically registered for the subscription.",
				ValidateFunc: validation.StringInSlice([]string{
					resourceproviders.ProviderRegistrationsCore,
					resourceproviders.ProviderRegistrationsExtended,
					resourceproviders.ProviderRegistrationsAll,
					resourceproviders.ProviderRegistrationsNone,
					resourceproviders.ProviderRegistrationsLegacy,
				}, false),
			},

			"resource_providers_to_register": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of Resource Providers to explicitly register for the subscription, in addition to those specified by the `resource_provider_registrations` property.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: resourceproviders.EnhancedValidate,
				},
			},

			// TODO: Remove `skip_provider_registration` in v5.0
			"skip_provider_registration": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SKIP_PROVIDER_REGISTRATION", nil),
				Description: "Should the AzureRM Provider skip registering all of the Resource Providers that it supports, if they're not already registered?",
				Deprecated:  features.DeprecatedInFourPointOh("This property is deprecated and will be removed in v5.0 of the AzureRM provider. Please use the `resource_provider_registrations` property instead."),
			},

			"storage_use_azuread": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_STORAGE_USE_AZUREAD", false),
				Description: "Should the AzureRM Provider use Azure AD Authentication when accessing the Storage Data Plane APIs?",
			},
		},

		DataSourcesMap: dataSources,
		ResourcesMap:   resources,
	}

	if !features.FourPointOhBeta() {
		p.Schema["subscription_id"].Required = false
		p.Schema["subscription_id"].Optional = true

		delete(p.Schema, "resource_provider_registrations")
		delete(p.Schema, "resource_providers_to_register")
	}

	p.ConfigureContextFunc = providerConfigure(p)

	return p
}

// providerConfigure is used to configure the cloud environment and authentication.
// To configure behavioral aspects of the provider, use the buildClient function instead.
// This separation allows us to robustly test different authentication scenarios.
func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		subscriptionId := d.Get("subscription_id").(string)
		if subscriptionId == "" {
			return nil, diag.FromErr(fmt.Errorf("`subscription_id` is a required provider property when performing a plan/apply operation"))
		}

		var auxTenants []string
		if v, ok := d.Get("auxiliary_tenant_ids").([]interface{}); ok && len(v) > 0 {
			auxTenants = *utils.ExpandStringSlice(v)
		} else if v := os.Getenv("ARM_AUXILIARY_TENANT_IDS"); v != "" {
			auxTenants = strings.Split(v, ";")
		}

		if len(auxTenants) > 3 {
			return nil, diag.Errorf("the provider only supports up to 3 auxiliary tenant IDs")
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

		tenantId, err := getTenantId(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		var (
			env *environments.Environment

			envName      = d.Get("environment").(string)
			metadataHost = d.Get("metadata_host").(string)
		)

		if metadataHost != "" {
			logEntry("[DEBUG] Configuring cloud environment from Metadata Service at %q", metadataHost)
			if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost)); err != nil {
				return nil, diag.FromErr(err)
			}
		} else {
			logEntry("[DEBUG] Configuring built-in cloud environment by name: %q", envName)
			if env, err = environments.FromName(envName); err != nil {
				return nil, diag.FromErr(err)
			}
		}

		var (
			enableAzureCli        = d.Get("use_cli").(bool)
			enableManagedIdentity = d.Get("use_msi").(bool)
			enableOidc            = d.Get("use_oidc").(bool) || d.Get("use_aks_workload_identity").(bool)
		)

		authConfig := &auth.Credentials{
			Environment:        *env,
			ClientID:           *clientId,
			TenantID:           *tenantId,
			AuxiliaryTenantIDs: auxTenants,

			ClientCertificateData:     clientCertificateData,
			ClientCertificatePath:     d.Get("client_certificate_path").(string),
			ClientCertificatePassword: d.Get("client_certificate_password").(string),
			ClientSecret:              *clientSecret,

			OIDCAssertionToken:          *oidcToken,
			GitHubOIDCTokenRequestURL:   d.Get("oidc_request_url").(string),
			GitHubOIDCTokenRequestToken: d.Get("oidc_request_token").(string),

			CustomManagedIdentityEndpoint: d.Get("msi_endpoint").(string),

			AzureCliSubscriptionIDHint: subscriptionId,

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

// buildClient is used to configure behavioral aspects of the provider. To configure the
// cloud environment and authentication-related settings, use the providerConfigure function.
func buildClient(ctx context.Context, p *schema.Provider, d *schema.ResourceData, authConfig *auth.Credentials) (*clients.Client, diag.Diagnostics) {
	// TODO: This hardcoded default is for v3.x, where `resource_provider_registrations` is not defined. Remove this hardcoded default in v4.0
	providerRegistrations := resourceproviders.ProviderRegistrationsLegacy
	if features.FourPointOhBeta() {
		providerRegistrations = d.Get("resource_provider_registrations").(string)
	}

	// TODO: Remove in v5.0
	if d.Get("skip_provider_registration").(bool) {
		if providerRegistrations != resourceproviders.ProviderRegistrationsLegacy {
			return nil, diag.Errorf("provider property `skip_provider_registration` cannot be set at the same time as `resource_provider_registrations`, please remove `skip_provider_registration` from your configuration or unset the `ARM_SKIP_PROVIDER_REGISTRATION` environment variable")
		}
		providerRegistrations = resourceproviders.ProviderRegistrationsNone
	}

	requiredResourceProviders, err := resourceproviders.GetResourceProvidersSet(providerRegistrations)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if features.FourPointOhBeta() {
		additionalProvidersToRegister := make(resourceproviders.ResourceProviders)
		for _, rp := range d.Get("resource_providers_to_register").([]interface{}) {
			additionalProvidersToRegister.Add(rp.(string))
		}
		requiredResourceProviders.Merge(additionalProvidersToRegister)
	}

	clientBuilder := clients.ClientBuilder{
		AuthConfig:                  authConfig,
		DisableCorrelationRequestID: d.Get("disable_correlation_request_id").(bool),
		DisableTerraformPartnerID:   d.Get("disable_terraform_partner_id").(bool),
		Features:                    expandFeatures(d.Get("features").([]interface{})),
		MetadataHost:                d.Get("metadata_host").(string),
		PartnerID:                   d.Get("partner_id").(string),
		RegisteredResourceProviders: requiredResourceProviders,
		StorageUseAzureAD:           d.Get("storage_use_azuread").(bool),
		SubscriptionID:              d.Get("subscription_id").(string),
		TerraformVersion:            p.TerraformVersion,

		// this field is intentionally not exposed in the provider block, since it's only used for
		// platform level tracing
		CustomCorrelationRequestID: os.Getenv("ARM_CORRELATION_REQUEST_ID"),
	}

	//lint:ignore SA1019 SDKv2 migration - staticcheck's own linter directives are currently being ignored under golangci-lint
	stopCtx, ok := schema.StopContext(ctx) //nolint:staticcheck
	if !ok {
		stopCtx = ctx
	}

	client, err := clients.Build(stopCtx, clientBuilder)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	client.StopContext = stopCtx

	subscriptionId := commonids.NewSubscriptionID(client.Account.SubscriptionId)

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	if err = resourceproviders.EnsureRegistered(ctx2, client.Resource.ResourceProvidersClient, subscriptionId, requiredResourceProviders); err != nil {
		return nil, diag.FromErr(err)

	}

	return client, nil
}
