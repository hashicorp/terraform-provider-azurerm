package authentication

import (
	"fmt"

	"log"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure/cli"
)

type Builder struct {
	// Core
	ClientID                 string
	SubscriptionID           string
	TenantID                 string
	Environment              string
	SkipProviderRegistration bool

	// Service Principal (Client Secret) Auth
	SupportsClientSecretAuth bool
	ClientSecret             string

	// Bearer Auth
	AccessToken  *adal.Token
	IsCloudShell bool
	UseMsi       bool
	MsiEndpoint  string
}

func (b Builder) Build() (*Config, error) {
	config := Config{
		ClientID:                 b.ClientID,
		SubscriptionID:           b.SubscriptionID,
		TenantID:                 b.TenantID,
		Environment:              b.Environment,
		SkipProviderRegistration: b.SkipProviderRegistration,

		// Bearer Auth
		AccessToken:  b.AccessToken,
		IsCloudShell: b.IsCloudShell,
		MsiEndpoint:  b.MsiEndpoint,
		UseMsi:       b.UseMsi,
	}

	if b.SupportsClientSecretAuth && b.ClientSecret != "" {
		log.Printf("[DEBUG] Using Service Principal / Client Secret Authentication")
		config.clientSecret = b.ClientSecret
		config.usingClientSecret = true
		config.AuthenticatedAsAServicePrincipal = true

		err := config.validateServicePrincipal()
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	// TODO: remove this in favour of individual calls above
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Config is the configuration structure used to instantiate a
// new Azure management client.
type Config struct {
	// Core
	ClientID                 string
	SubscriptionID           string
	TenantID                 string
	Environment              string
	SkipProviderRegistration bool

	// Service Principal (Client Secret) Auth
	clientSecret string

	// Bearer Auth
	AccessToken  *adal.Token
	IsCloudShell bool
	UseMsi       bool
	MsiEndpoint  string

	// internal-only feature flags
	usingClientSecret bool

	// temporarily public feature flags
	AuthenticatedAsAServicePrincipal bool
}

// LoadTokensFromAzureCLI loads the access tokens and subscription/tenant ID's from the
// Azure CLI metadata if it's not provided
// NOTE: this'll become an internal-only method in the near future
func (c *Config) LoadTokensFromAzureCLI() error {
	profilePath, err := cli.ProfilePath()
	if err != nil {
		return fmt.Errorf("Error loading the Profile Path from the Azure CLI: %+v", err)
	}

	profile, err := cli.LoadProfile(profilePath)
	if err != nil {
		return fmt.Errorf("Azure CLI Authorization Profile was not found. Please ensure the Azure CLI is installed and then log-in with `az login`.")
	}

	cliProfile := AzureCLIProfile{
		Profile: profile,
	}

	// find the Subscription ID if it's not specified
	if c.SubscriptionID == "" {
		// we want to expose a more friendly error to the user, but this is useful for debug purposes
		err := c.populateSubscriptionFromCLIProfile(cliProfile)
		if err != nil {
			log.Printf("Error Populating the Subscription from the CLI Profile: %s", err)
		}
	}

	// find the Tenant ID for that subscription if they're not specified
	if c.TenantID == "" {
		err := c.populateTenantFromCLIProfile(cliProfile)
		if err != nil {
			// we want to expose a more friendly error to the user, but this is useful for debug purposes
			log.Printf("Error Populating the Tenant from the CLI Profile: %s", err)
		}
	}

	foundToken := false
	if c.TenantID != "" {
		// pull out the ClientID and the AccessToken from the Azure Access Token
		tokensPath, err := cli.AccessTokensPath()
		if err != nil {
			return fmt.Errorf("Error loading the Tokens Path from the Azure CLI: %+v", err)
		}

		tokens, err := cli.LoadTokens(tokensPath)
		if err != nil {
			return fmt.Errorf("Azure CLI Authorization Tokens were not found. Please ensure the Azure CLI is installed and then log-in with `az login`.")
		}

		validToken, _ := findValidAccessTokenForTenant(tokens, c.TenantID)
		if validToken != nil {
			foundToken, err = c.populateFromAccessToken(validToken)
			if err != nil {
				return err
			}
		}
	}

	if !foundToken {
		return fmt.Errorf("No valid (unexpired) Azure CLI Auth Tokens found. Please run `az login`.")
	}

	// always pull the Environment from the CLI
	err = c.populateEnvironmentFromCLIProfile(cliProfile)
	if err != nil {
		// we want to expose a more friendly error to the user, but this is useful for debug purposes
		log.Printf("Error Populating the Environment from the CLI Profile: %s", err)
	}

	return nil
}

func (c *Config) populateSubscriptionFromCLIProfile(cliProfile AzureCLIProfile) error {
	subscriptionId, err := cliProfile.FindDefaultSubscriptionId()
	if err != nil {
		return err
	}

	c.SubscriptionID = subscriptionId
	return nil
}

func (c *Config) populateTenantFromCLIProfile(cliProfile AzureCLIProfile) error {
	subscription, err := cliProfile.FindSubscription(c.SubscriptionID)
	if err != nil {
		return err
	}

	if c.TenantID == "" {
		c.TenantID = subscription.TenantID
	}

	return nil
}

func (c *Config) populateEnvironmentFromCLIProfile(cliProfile AzureCLIProfile) error {
	subscription, err := cliProfile.FindSubscription(c.SubscriptionID)
	if err != nil {
		return err
	}

	c.Environment = normalizeEnvironmentName(subscription.EnvironmentName)

	return nil
}

func (c *Config) populateFromAccessToken(token *AccessToken) (bool, error) {
	if token == nil {
		return false, fmt.Errorf("No valid access token was found to populate from")
	}

	c.ClientID = token.ClientID
	c.AccessToken = token.AccessToken
	c.IsCloudShell = token.IsCloudShell

	return true, nil
}
