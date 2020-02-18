package clients

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

type ClientBuilder struct {
	AuthConfig                  *authentication.Config
	DisableCorrelationRequestID bool
	DisableTerraformPartnerID   bool
	PartnerId                   string
	SkipProviderRegistration    bool
	StorageUseAzureAD           bool
	TerraformVersion            string
	Features                    features.UserFeatures
}

func Build(ctx context.Context, builder ClientBuilder) (*Client, error) {
	env, err := authentication.DetermineEnvironment(builder.AuthConfig.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	account, err := NewResourceManagerAccount(ctx, *builder.AuthConfig, *env)
	if err != nil {
		return nil, fmt.Errorf("Error building account: %+v", err)
	}

	client := Client{
		Account: account,
	}

	oauthConfig, err := builder.AuthConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", builder.AuthConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Storage Endpoints
	storageAuth, err := builder.AuthConfig.GetAuthorizationToken(sender, oauthConfig, env.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, err
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
		SkipProviderReg:             builder.SkipProviderRegistration,
		DisableCorrelationRequestID: builder.DisableCorrelationRequestID,
		DisableTerraformPartnerID:   builder.DisableTerraformPartnerID,
		Environment:                 *env,
		Features:                    builder.Features,
		StorageUseAzureAD:           builder.StorageUseAzureAD,
	}

	if err := client.Build(ctx, o); err != nil {
		return nil, fmt.Errorf("Error building Client: %+v", err)
	}

	return &client, nil
}
