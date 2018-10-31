package authentication

import (
	"fmt"
	"log"

	"github.com/Azure/go-autorest/autorest/adal"
)

type Builder struct {
	// Core
	ClientID                 string
	SubscriptionID           string
	TenantID                 string
	Environment              string
	SkipProviderRegistration bool

	// Azure CLI Parsing / CloudShell Auth
	SupportsAzureCliCloudShellParsing bool

	// Managed Service Identity Auth
	SupportsManagedServiceIdentity bool
	MsiEndpoint                    string

	// Service Principal (Client Secret) Auth
	SupportsClientSecretAuth bool
	ClientSecret             string
}

func (b Builder) Build() (*Config, error) {
	config := Config{
		ClientID:                 b.ClientID,
		SubscriptionID:           b.SubscriptionID,
		TenantID:                 b.TenantID,
		Environment:              b.Environment,
		SkipProviderRegistration: b.SkipProviderRegistration,
	}

	if b.SupportsClientSecretAuth && b.ClientSecret != "" {
		log.Printf("[DEBUG] Using Service Principal / Client Secret for Authentication")
		config.clientSecret = b.ClientSecret
		config.usingClientSecret = true
		config.AuthenticatedAsAServicePrincipal = true

		err := config.validateServicePrincipalClientSecret()
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	if b.SupportsManagedServiceIdentity {
		log.Printf("[DEBUG] Using Managed Service Identity for Authentication")

		endpoint := b.MsiEndpoint
		if endpoint == "" {
			msiEndpoint, err := adal.GetMSIVMEndpoint()
			if err != nil {
				return nil, fmt.Errorf("Could not retrieve MSI endpoint from VM settings."+
					"Ensure the VM has MSI enabled, or configure the MSI Endpoint. Error: %s", err)
			}
			endpoint = msiEndpoint
		}

		log.Printf("[DEBUG] Using MSI endpoint %q", endpoint)
		config.msiEndpoint = endpoint

		err := config.validateMsi()
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	// note: this includes CloudShell
	if b.SupportsAzureCliCloudShellParsing {
		log.Printf("[DEBUG] Using CloudShell for Authentication")

		// given CloudShell is technically parsing Azure CLI creds - we set both
		// however usingCloudShell is set inside the LoadTokensFromAzureCLI method
		config.usingAzureCliParsing = true

		// load the refreshed tokens from the Azure CLI
		err := config.LoadTokensFromAzureCLI()
		if err != nil {
			return nil, fmt.Errorf("Error loading the refreshed Azure CLI token: %+v", err)
		}

		err = config.validateAzureCliBearerAuth()
		if err != nil {
			return nil, err
		}

		return &config, nil
	}

	return nil, fmt.Errorf("No supported authentication methods were found!")
}
