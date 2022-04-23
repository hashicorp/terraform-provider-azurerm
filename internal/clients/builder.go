package clients

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
	"github.com/manicminer/hamilton/environments"
)

type ClientBuilder struct {
	AuthConfig                  *authentication.Config
	DisableCorrelationRequestID bool
	CustomCorrelationRequestID  string
	DisableTerraformPartnerID   bool
	PartnerId                   string
	SkipProviderRegistration    bool
	StorageUseAzureAD           bool
	TerraformVersion            string
	Features                    features.UserFeatures
}

const azureStackEnvironmentError = `
The AzureRM Provider supports the different Azure Public Clouds - including China, Public,
and US Government - however it does not support Azure Stack due to differences in API and
feature availability.

Terraform instead offers a separate "azurestack" provider which supports the functionality
and APIs available in Azure Stack via Azure Stack Profiles.
`

func Build(ctx context.Context, builder ClientBuilder) (*Client, error) {
	// point folks towards the separate Azure Stack Provider when using Azure Stack
	if strings.EqualFold(builder.AuthConfig.Environment, "AZURESTACKCLOUD") {
		return nil, fmt.Errorf(azureStackEnvironmentError)
	}

	isAzureStack, err := authentication.IsEnvironmentAzureStack(ctx, builder.AuthConfig.MetadataHost, builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to determine if environment is Azure Stack: %+v", err)
	}
	if isAzureStack {
		return nil, fmt.Errorf(azureStackEnvironmentError)
	}

	// Autorest environment configuration
	env, err := authentication.AzureEnvironmentByNameFromEndpoint(ctx, builder.AuthConfig.MetadataHost, builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", builder.AuthConfig.Environment, builder.AuthConfig.MetadataHost, err)
	}

	// Hamilton environment configuration
	environment, err := environments.EnvironmentFromString(builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", builder.AuthConfig.Environment, builder.AuthConfig.MetadataHost, err)
	}

	// client declarations:
	account, err := NewResourceManagerAccount(ctx, *builder.AuthConfig, *env, builder.SkipProviderRegistration)
	if err != nil {
		return nil, fmt.Errorf("building account: %+v", err)
	}

	client := Client{
		Account: account,
	}

	oauthConfig, err := builder.AuthConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building OAuth Config: %+v", err)
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("unable to configure OAuthConfig for tenant %s", builder.AuthConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	var auth, storageAuth, synapseAuth, batchManagementAuth autorest.Authorizer
	var keyVaultAuth *autorest.BearerAuthorizerCallback
	var tokenFunc common.EndpointTokenFunc

	auth, err = builder.AuthConfig.GetMSALToken(ctx, environment.ResourceManager, sender, oauthConfig, string(environment.ResourceManager.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to get MSAL authorization token for resource manager API: %+v", err)
	}

	storageAuth, err = builder.AuthConfig.GetMSALToken(ctx, environment.Storage, sender, oauthConfig, string(environment.Storage.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to get MSAL authorization token for storage API: %+v", err)
	}

	if environment.Synapse.IsAvailable() {
		synapseAuth, err = builder.AuthConfig.GetMSALToken(ctx, environment.Synapse, sender, oauthConfig, string(environment.Synapse.Endpoint))
		if err != nil {
			return nil, fmt.Errorf("unable to get MSAL authorization token for synapse API: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Skipping building the Synapse MSAL Authorizer since this is not supported in the current Azure Environment")
	}

	batchManagementAuth, err = builder.AuthConfig.GetMSALToken(ctx, environment.BatchManagement, sender, oauthConfig, string(environment.BatchManagement.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to get MSAL authorization token for batch management API: %+v", err)
	}

	keyVaultAuth = builder.AuthConfig.MSALBearerAuthorizerCallback(ctx, environment.KeyVault, sender, oauthConfig, string(environment.KeyVault.Endpoint))

	// Helper for obtaining endpoint-specific tokens
	tokenFunc = func(endpoint string) (autorest.Authorizer, error) {
		api := environments.Api{Endpoint: environments.ApiEndpoint(endpoint)}
		authorizer, err := builder.AuthConfig.GetMSALToken(ctx, api, sender, oauthConfig, endpoint)
		if err != nil {
			return nil, fmt.Errorf("getting MSAL authorization token for endpoint %s: %+v", endpoint, err)
		}
		return authorizer, nil
	}

	o := &common.ClientOptions{
		SubscriptionId:              builder.AuthConfig.SubscriptionID,
		TenantID:                    builder.AuthConfig.TenantID,
		PartnerId:                   builder.PartnerId,
		TerraformVersion:            builder.TerraformVersion,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     env.ResourceManagerEndpoint,
		StorageAuthorizer:           storageAuth,
		SynapseAuthorizer:           synapseAuth,
		BatchManagementAuthorizer:   batchManagementAuth,
		SkipProviderReg:             builder.SkipProviderRegistration,
		DisableCorrelationRequestID: builder.DisableCorrelationRequestID,
		CustomCorrelationRequestID:  builder.CustomCorrelationRequestID,
		DisableTerraformPartnerID:   builder.DisableTerraformPartnerID,
		Environment:                 *env,
		Features:                    builder.Features,
		StorageUseAzureAD:           builder.StorageUseAzureAD,
		TokenFunc:                   tokenFunc,
	}

	if err := client.Build(ctx, o); err != nil {
		return nil, fmt.Errorf("building Client: %+v", err)
	}

	if features.EnhancedValidationEnabled() {
		location.CacheSupportedLocations(ctx, env.ResourceManagerEndpoint)
		resourceproviders.CacheSupportedProviders(ctx, client.Resource.ProvidersClient)
	}

	return &client, nil
}
