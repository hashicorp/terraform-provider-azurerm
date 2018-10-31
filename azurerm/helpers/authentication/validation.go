package authentication

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
)

// Validate ensures that the current set of credentials provided
// are valid for the selected authentication type (e.g. Client Secret, Azure CLI, MSI etc.)
func (c *Config) Validate() error {
	// Azure CLI / CloudShell
	log.Printf("[DEBUG] No Client Secret specified - loading credentials from Azure CLI")
	if err := c.LoadTokensFromAzureCLI(); err != nil {
		return err
	}

	if err := c.validateAzureCliBearerAuth(); err != nil {
		return fmt.Errorf("Please specify either a Service Principal, or log in with the Azure CLI (using `az login`)")
	}

	return nil
}

func (c *Config) validateAzureCliBearerAuth() error {
	var err *multierror.Error

	if c.AccessToken == nil {
		err = multierror.Append(err, fmt.Errorf("Access Token was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID was not found in your Azure CLI Credentials.\n\nPlease login to the Azure CLI again via `az login`"))
	}

	return err.ErrorOrNil()
}

func (c *Config) validateServicePrincipal() error {
	var err *multierror.Error

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID must be configured for the AzureRM provider"))
	}
	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID must be configured for the AzureRM provider"))
	}
	if c.clientSecret == "" {
		err = multierror.Append(err, fmt.Errorf("Client Secret must be configured for the AzureRM provider"))
	}
	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID must be configured for the AzureRM provider"))
	}
	if c.Environment == "" {
		err = multierror.Append(err, fmt.Errorf("Environment must be configured for the AzureRM provider"))
	}

	return err.ErrorOrNil()
}

func (c *Config) validateMsi() error {
	var err *multierror.Error

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID must be configured for the AzureRM provider"))
	}
	if c.TenantID == "" {
		err = multierror.Append(err, fmt.Errorf("Tenant ID must be configured for the AzureRM provider"))
	}
	if c.Environment == "" {
		err = multierror.Append(err, fmt.Errorf("Environment must be configured for the AzureRM provider"))
	}
	if c.msiEndpoint == "" {
		err = multierror.Append(err, fmt.Errorf("MSI endpoint must be configured for the AzureRM provider"))
	}

	return err.ErrorOrNil()
}
