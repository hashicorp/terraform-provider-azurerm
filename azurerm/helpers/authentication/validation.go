package authentication

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (c *Config) ValidateBearerAuth() error {
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

func (c *Config) ValidateServicePrincipal() error {
	var err *multierror.Error

	if c.SubscriptionID == "" {
		err = multierror.Append(err, fmt.Errorf("Subscription ID must be configured for the AzureRM provider"))
	}
	if c.ClientID == "" {
		err = multierror.Append(err, fmt.Errorf("Client ID must be configured for the AzureRM provider"))
	}
	if c.ClientSecret == "" {
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

func (c *Config) ValidateMsi() error {
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
	if c.MsiEndpoint == "" {
		err = multierror.Append(err, fmt.Errorf("MSI endpoint must be configured for the AzureRM provider"))
	}

	return err.ErrorOrNil()
}
