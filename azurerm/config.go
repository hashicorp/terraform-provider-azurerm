package azurerm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/go-azure-helpers/sender"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

// ArmClient contains the handles to all the specific Azure Resource Manager
// resource classes' respective clients.
type ArmClient struct {
	// inherit the fields from the parent, so that we should be able to set/access these at either level
	clients.Client
}

type armClientBuilder struct {
	authConfig                  *authentication.Config
	skipProviderRegistration    bool
	terraformVersion            string
	partnerId                   string
	disableCorrelationRequestID bool
	disableTerraformPartnerID   bool
}

// getArmClient is a helper method which returns a fully instantiated
// *ArmClient based on the Config's current settings.
func getArmClient(ctx context.Context, builder armClientBuilder) (*ArmClient, error) {
	env, err := authentication.DetermineEnvironment(builder.authConfig.Environment)
	if err != nil {
		return nil, err
	}

	// client declarations:
	account, err := clients.NewResourceManagerAccount(ctx, *builder.authConfig, *env)
	if err != nil {
		return nil, fmt.Errorf("Error building account: %+v", err)
	}

	client := ArmClient{
		Client: clients.Client{
			Account: account,
		},
	}

	oauthConfig, err := builder.authConfig.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, err
	}

	// OAuthConfigForTenant returns a pointer, which can be nil.
	if oauthConfig == nil {
		return nil, fmt.Errorf("Unable to configure OAuthConfig for tenant %s", builder.authConfig.TenantID)
	}

	sender := sender.BuildSender("AzureRM")

	// Resource Manager endpoints
	endpoint := env.ResourceManagerEndpoint
	auth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, env.TokenAudience)
	if err != nil {
		return nil, err
	}

	// Graph Endpoints
	graphEndpoint := env.GraphEndpoint
	graphAuth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, graphEndpoint)
	if err != nil {
		return nil, err
	}

	// Storage Endpoints
	storageAuth, err := builder.authConfig.GetAuthorizationToken(sender, oauthConfig, env.ResourceIdentifiers.Storage)
	if err != nil {
		return nil, err
	}

	// Key Vault Endpoints
	keyVaultAuth := builder.authConfig.BearerAuthorizerCallback(sender, oauthConfig)

	o := &common.ClientOptions{
		SubscriptionId:              builder.authConfig.SubscriptionID,
		TenantID:                    builder.authConfig.TenantID,
		PartnerId:                   builder.partnerId,
		TerraformVersion:            builder.terraformVersion,
		GraphAuthorizer:             graphAuth,
		GraphEndpoint:               graphEndpoint,
		KeyVaultAuthorizer:          keyVaultAuth,
		ResourceManagerAuthorizer:   auth,
		ResourceManagerEndpoint:     endpoint,
		StorageAuthorizer:           storageAuth,
		PollingDuration:             180 * time.Minute,
		SkipProviderReg:             builder.skipProviderRegistration,
		DisableCorrelationRequestID: builder.disableCorrelationRequestID,
		DisableTerraformPartnerID:   builder.disableTerraformPartnerID,
		Environment:                 *env,
	}

	if err := client.Client.Build(o); err != nil {
		return nil, fmt.Errorf("Error building Client: %+v", err)
	}

	return &client, nil
}
