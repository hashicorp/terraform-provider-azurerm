package cdnfrontdoorsecretparams

import (
	"fmt"

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

func ExpandCdnFrontdoorCustomerCertificateParameters(input []interface{}) (*track1.BasicSecretParameters, error) {
	m := InitializeCdnFrontdoorSecretMappings()
	item := input[0].(map[string]interface{})

	// must create the secret_source
	kv := item["key_vault_id"].(string)
	certName := item["key_vault_certificate_name"].(string)

	kvId, err := keyVaultParse.VaultID(kv)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q, %q", "key_vault_id", kv)
	}

	secretSource := frontdoorParse.NewFrontdoorKeyVaultSecretID(kvId.SubscriptionId, kvId.ResourceGroup, kvId.Name, certName)

	customerCertificate := &track1.CustomerCertificateParameters{
		Type: m.CustomerCertificate.TypeName,
		SecretSource: &track1.ResourceReference{
			ID: utils.String(secretSource.ID()),
		},
		SecretVersion:           utils.String(item["key_vault_certificate_version"].(string)),
		UseLatestVersion:        utils.Bool(item["use_latest"].(bool)),
		SubjectAlternativeNames: utils.ExpandStringSlice(item["subject_alternative_names"].([]interface{})),
	}

	if secretParameter := track1.BasicSecretParameters(customerCertificate); secretParameter != nil {
		return &secretParameter, nil
	}

	return nil, nil
}
