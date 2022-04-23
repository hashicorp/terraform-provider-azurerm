package cdnfrontdoorsecretparams

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	frontdoorParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorSecretParameters struct {
	TypeName   track1.TypeBasicSecretParameters
	ConfigName string
}

type CdnFrontdoorSecretMappings struct {
	UrlSigningKey                     CdnFrontdoorSecretParameters
	ManagedCertificate                CdnFrontdoorSecretParameters
	CustomerCertificate               CdnFrontdoorSecretParameters
	AzureFirstPartyManagedCertificate CdnFrontdoorSecretParameters
}

func InitializeCdnFrontdoorSecretMappings() *CdnFrontdoorSecretMappings {
	m := new(CdnFrontdoorSecretMappings)

	m.UrlSigningKey = CdnFrontdoorSecretParameters{
		TypeName:   track1.TypeBasicSecretParametersTypeURLSigningKey,
		ConfigName: "url_signing_key",
	}

	m.ManagedCertificate = CdnFrontdoorSecretParameters{
		TypeName:   track1.TypeBasicSecretParametersTypeManagedCertificate,
		ConfigName: "managed_certificate",
	}

	m.CustomerCertificate = CdnFrontdoorSecretParameters{
		TypeName:   track1.TypeBasicSecretParametersTypeCustomerCertificate,
		ConfigName: "customer_certificate",
	}

	m.AzureFirstPartyManagedCertificate = CdnFrontdoorSecretParameters{
		TypeName:   track1.TypeBasicSecretParametersTypeAzureFirstPartyManagedCertificate,
		ConfigName: "azure_first_party_managed_certificate",
	}

	return m
}

func ExpandCdnFrontdoorCustomerCertificateParameters(ctx context.Context, input []interface{}, clients *clients.Client) (*track1.BasicSecretParameters, error) {
	m := InitializeCdnFrontdoorSecretMappings()
	item := input[0].(map[string]interface{})

	// New Direction: Parse the certificate id (e.g. URL) and derive the rest of the information from there...
	certificateBaseURL := item["key_vault_certificate_id"].(string)
	certificateId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(certificateBaseURL)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Certificate Resource ID from the Key Vault Certificate Base URL %q: %s", certificateBaseURL, err)
	}

	var useLatest bool
	if certificateId.Version == "" {
		useLatest = true
	}

	keyVaultBaseId, err := clients.KeyVault.KeyVaultIDFromBaseUrl(ctx, clients.Resource, certificateId.KeyVaultBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("retrieving the Key Vault Resource ID from the Key Vault Base URL %q: %s", certificateId.KeyVaultBaseUrl, err)
	}

	if keyVaultBaseId == nil {
		return nil, fmt.Errorf("unexpected %q Key Vault Resource ID retrieved from the Key Vault Base URL %q", "nil", certificateId.KeyVaultBaseUrl)
	}

	keyVaultId, err := keyVaultParse.VaultID(*keyVaultBaseId)
	if err != nil {
		return nil, err
	}

	secretSource := frontdoorParse.NewFrontdoorKeyVaultSecretID(keyVaultId.SubscriptionId, keyVaultId.ResourceGroup, keyVaultId.Name, certificateId.Name)

	customerCertificate := &track1.CustomerCertificateParameters{
		Type: m.CustomerCertificate.TypeName,
		SecretSource: &track1.ResourceReference{
			ID: utils.String(secretSource.ID()),
		},
		UseLatestVersion: utils.Bool(useLatest),
	}

	if !useLatest {
		customerCertificate.SecretVersion = utils.String(certificateId.Version)
	}

	if secretParameter := track1.BasicSecretParameters(customerCertificate); secretParameter != nil {
		return &secretParameter, nil
	}

	return nil, fmt.Errorf("unexpected %q Customer Certificate received from the Key Vault Certificate Base URL %q", "nil", certificateId.KeyVaultBaseUrl)
}
