package frontdoor

import (
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func schemaCustomHttpsConfiguration() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"certificate_source": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(frontdoor.CertificateSourceFrontDoor),
			ValidateFunc: validation.StringInSlice([]string{
				string(frontdoor.CertificateSourceAzureKeyVault),
				string(frontdoor.CertificateSourceFrontDoor),
			}, false),
		},
		"minimum_tls_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provisioning_substate": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// NOTE: None of these attributes are valid if
		//       certificate_source is set to FrontDoor
		"azure_key_vault_certificate_secret_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"azure_key_vault_certificate_secret_version": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"azure_key_vault_certificate_vault_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

type flattenedCustomHttpsConfiguration struct {
	CustomHTTPSConfiguration       []interface{}
	CustomHTTPSProvisioningEnabled bool
}

func flattenCustomHttpsConfiguration(properties *frontdoor.FrontendEndpointProperties) flattenedCustomHttpsConfiguration {
	result := flattenedCustomHttpsConfiguration{
		CustomHTTPSConfiguration:       make([]interface{}, 0),
		CustomHTTPSProvisioningEnabled: false,
	}
	if properties == nil {
		return result
	}

	if config := properties.CustomHTTPSConfiguration; config != nil {
		certificateSource := string(frontdoor.CertificateSourceFrontDoor)

		keyVaultCertificateVaultId := ""
		keyVaultCertificateSecretName := ""
		keyVaultCertificateSecretVersion := ""
		if config.CertificateSource == frontdoor.CertificateSourceAzureKeyVault {
			if vault := config.KeyVaultCertificateSourceParameters; vault != nil {
				certificateSource = string(frontdoor.CertificateSourceAzureKeyVault)

				if vault.Vault != nil && vault.Vault.ID != nil {
					keyVaultCertificateVaultId = *vault.Vault.ID
				}

				if vault.SecretName != nil {
					keyVaultCertificateSecretName = *vault.SecretName
				}

				if vault.SecretVersion != nil {
					keyVaultCertificateSecretVersion = *vault.SecretVersion
				}
			}
		}

		provisioningState := ""
		provisioningSubstate := ""
		if properties.CustomHTTPSProvisioningState != "" {
			provisioningState = string(properties.CustomHTTPSProvisioningState)
			if properties.CustomHTTPSProvisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled || properties.CustomHTTPSProvisioningState == frontdoor.CustomHTTPSProvisioningStateEnabling {
				result.CustomHTTPSProvisioningEnabled = true

				if properties.CustomHTTPSProvisioningSubstate != "" {
					provisioningSubstate = string(properties.CustomHTTPSProvisioningSubstate)
				}
			}

			result.CustomHTTPSConfiguration = append(result.CustomHTTPSConfiguration, map[string]interface{}{
				"azure_key_vault_certificate_vault_id":       keyVaultCertificateVaultId,
				"azure_key_vault_certificate_secret_name":    keyVaultCertificateSecretName,
				"azure_key_vault_certificate_secret_version": keyVaultCertificateSecretVersion,
				"certificate_source":                         certificateSource,
				"minimum_tls_version":                        string(config.MinimumTLSVersion),
				"provisioning_state":                         provisioningState,
				"provisioning_substate":                      provisioningSubstate,
			})
		}
	}

	return result
}
