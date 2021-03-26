package clients

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceproviders"
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
The AzureRM Provider supports the different Azure Public Clouds - including China, Germany,
Public and US Government - however it does not support Azure Stack due to differences in
API and feature availability.

Terraform instead offers a separate "azurestack" provider which supports the functionality
and API's available in Azure Stack via Azure Stack Profiles.
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

	env, err := authentication.AzureEnvironmentByNameFromEndpoint(ctx, builder.AuthConfig.MetadataHost, builder.AuthConfig.Environment)
	if err != nil {
		return nil, fmt.Errorf("unable to find environment %q from endpoint %q: %+v", builder.AuthConfig.Environment, builder.AuthConfig.MetadataHost, err)
	}

	// client declarations:
	account, err := NewResourceManagerAccount(ctx, *builder.AuthConfig, *env, builder.SkipProviderRegistration)
	if err != nil {
		return nil, fmt.Errorf("Error building account: %+v", err)
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

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, fmt.Errorf("unable to get authorization token for resource manager: %+v", err)
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to get authorization token for graph endpoints: %+v", err)
	}

	// Storage Endpoints
	storageAuth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, fmt.Errorf("unable to get authorization token for storage endpoints: %+v", err)
	}

	// Synapse Endpoints
	var synapseAuth autorest.Authorizer = nil
	if env.ResourceIdentifiers.Synapse != azure.NotAvailable {
		synapseAuth, err = builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.ResourceIdentifiers.Synapse)
		if err != nil {
			return nil, fmt.Errorf("unable to get authorization token for synapse endpoints: %+v", err)
		}
	} else {
		log.Printf("[DEBUG] Skipping building the Synapse Authorizer since this is not supported in the current Azure Environment")
	}

	// Key Vault Endpoints
	keyVaultAuth := builder.AuthConfig.BearerAuthorizerCallback(sender, oauthConfig)

	o := &common.ClientOptions{
		SubscriptionId:              builder.AuthConfig.SubscriptionID,
		TenantID:                    builder.AuthConfig.TenantID,
		PartnerId:                   builder.PartnerId,
		TerraformVersion:            builder.TerraformVersion,
		GraphAuthorizer:             graphAuth,
		GraphEndpoint:               graphEndpoint,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     endpoint,
		StorageAuthorizer:           storageAuth,
		SynapseAuthorizer:           synapseAuth,
		SkipProviderReg:             builder.SkipProviderRegistration,
		DisableCorrelationRequestID: builder.DisableCorrelationRequestID,
		CustomCorrelationRequestID:  builder.CustomCorrelationRequestID,
		DisableTerraformPartnerID:   builder.DisableTerraformPartnerID,
		Environment:                 *env,
		Features:                    builder.Features,
		StorageUseAzureAD:           builder.StorageUseAzureAD,
	}

	if err := client.Build(ctx, o); err != nil {
		return nil, fmt.Errorf("error building Client: %+v", err)
	}

	if features.EnhancedValidationEnabled() {
		location.CacheSupportedLocations(ctx, env)
		resourceproviders.CacheSupportedProviders(ctx, client.Resource.ProvidersClient)
	}

	return &client, nil
}
