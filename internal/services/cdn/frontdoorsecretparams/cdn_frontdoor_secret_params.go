package cdnfrontdoorsecretparams

import (
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
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

func ExpandCdnFrontdoorUrlSigningKeyParameters(input []interface{}) (*track1.BasicSecretParameters, error) {
	// Placeholder
	return nil, nil
}

func ExpandCdnFrontdoorManagedCertificateParameters(input []interface{}) (*track1.BasicSecretParameters, error) {
	// Placeholder
	return nil, nil
}

func ExpandCdnFrontdoorCustomerCertificateParameters(input []interface{}) (*track1.BasicSecretParameters, error) {
	m := InitializeCdnFrontdoorSecretMappings()
	item := input[0].(map[string]interface{})

	customerCertificate := &track1.CustomerCertificateParameters{
		Type: m.CustomerCertificate.TypeName,
		SecretSource: &track1.ResourceReference{
			ID: utils.String(item["secret_source_id"].(string)),
		},
		SecretVersion:           utils.String(item["secret_version"].(string)),
		UseLatestVersion:        utils.Bool(item["use_latest"].(bool)),
		SubjectAlternativeNames: utils.ExpandStringSlice(item["subject_alternative_names"].([]interface{})),
	}

	if secretParameter := track1.BasicSecretParameters(customerCertificate); secretParameter != nil {
		return &secretParameter, nil
	}

	return nil, nil
}

func ExpandCdnFrontdoorAzureFirstPartyManagedCertificateyParameters(input []interface{}) (*track1.BasicSecretParameters, error) {
	// Placeholder
	return nil, nil
}

func FlattenCdnFrontdoorUrlSigningKeyParameters(input track1.BasicSecretParameters) (map[string]interface{}, error) {
	// Placeholder
	return nil, nil
}

func FlattenCdnFrontdoorManagedCertificateParameters(input track1.BasicSecretParameters) (map[string]interface{}, error) {
	// Placeholder
	return nil, nil
}

// func flattenCdnFrontdoorCustomerCertificateParameters(input *track1.BasicSecretParameters) (map[string]interface{}, error) {
// 	// TODO: Get this working once I can successfully create a secret
// 	return nil, nil
// }

func FlattenCdnFrontdoorAzureFirstPartyManagedCertificateyParameters(input track1.BasicSecretParameters) (map[string]interface{}, error) {
	// Placeholder
	return nil, nil
}
