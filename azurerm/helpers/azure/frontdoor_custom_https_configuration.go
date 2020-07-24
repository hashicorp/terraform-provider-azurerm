package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func SchemaFrontdoorCustomHttpsConfiguration() map[string]*schema.Schema {
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

func FlattenArmFrontDoorCustomHttpsConfiguration(input *frontdoor.FrontendEndpoint, output map[string]interface{}, resourceGroup string) error {
	if input == nil {
		return fmt.Errorf("cannot read Front Door Frontend Endpoint (Resource Group %q): endpoint is empty", resourceGroup)
	}

	customHttpsConfiguration := make([]interface{}, 0)
	chc := make(map[string]interface{})

	if properties := input.FrontendEndpointProperties; properties != nil {
		if properties.CustomHTTPSConfiguration != nil {
			customHTTPSConfiguration := properties.CustomHTTPSConfiguration
			if customHTTPSConfiguration.CertificateSource == frontdoor.CertificateSourceAzureKeyVault {
				if kvcsp := customHTTPSConfiguration.KeyVaultCertificateSourceParameters; kvcsp != nil {
					chc["certificate_source"] = string(frontdoor.CertificateSourceAzureKeyVault)
					chc["azure_key_vault_certificate_vault_id"] = *kvcsp.Vault.ID
					chc["azure_key_vault_certificate_secret_name"] = *kvcsp.SecretName
					chc["azure_key_vault_certificate_secret_version"] = *kvcsp.SecretVersion
				}
			} else {
				chc["certificate_source"] = string(frontdoor.CertificateSourceFrontDoor)
			}

			chc["minimum_tls_version"] = string(customHTTPSConfiguration.MinimumTLSVersion)

			if provisioningState := properties.CustomHTTPSProvisioningState; provisioningState != "" {
				chc["provisioning_state"] = provisioningState
				if provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled || provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabling {
					output["custom_https_provisioning_enabled"] = true

					if provisioningSubstate := properties.CustomHTTPSProvisioningSubstate; provisioningSubstate != "" {
						chc["provisioning_substate"] = provisioningSubstate
					}
				} else {
					output["custom_https_provisioning_enabled"] = false
				}

				customHttpsConfiguration = append(customHttpsConfiguration, chc)
				output["custom_https_configuration"] = customHttpsConfiguration
			}
		}
	}

	return nil
}
