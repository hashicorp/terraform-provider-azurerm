// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func customizeHttpsConfigurationCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
	if helpers.IsFrontDoorFullyRetired() {
		return fmt.Errorf("%s", helpers.FullyRetiredMessage)
	}

	// New resources are not supported, and since this is a ForceNew field, we also need to block changing the field as
	// the re-create would fail with the create error from the service API...
	if helpers.IsFrontDoorDeprecatedForCreation() && d.HasChange("frontend_endpoint_id") {
		return fmt.Errorf("%s", helpers.CreateDeprecationMessage)
	}

	if v, ok := d.GetOk("frontend_endpoint_id"); ok && v.(string) != "" {
		// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
		// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
		id, err := frontdoors.ParseFrontendEndpointIDInsensitively(v.(string))
		if err != nil {
			return err
		}

		if err := customHttpsSettings(d); err != nil {
			return fmt.Errorf("validating Front Door Custom Https Configuration for Endpoint %q (Front Door %q / Resource Group %q): %+v", id.FrontendEndpointName, id.FrontDoorName, id.ResourceGroupName, err)
		}
	}

	return nil
}

func customHttpsSettings(d *pluginsdk.ResourceDiff) error {
	frontendEndpointCustomHttpsConfig := d.Get("custom_https_configuration").([]any)
	customHttpsEnabled := d.Get("custom_https_provisioning_enabled").(bool)

	if len(frontendEndpointCustomHttpsConfig) > 0 {
		if !customHttpsEnabled {
			return errors.New(`"custom_https_provisioning_enabled" is set to "false". please remove the "custom_https_configuration" block from the configuration file`)
		}

		// Verify frontend endpoints custom https configuration is valid if defined
		if err := verifyCustomHttpsConfiguration(frontendEndpointCustomHttpsConfig); err != nil {
			return err
		}
	} else if customHttpsEnabled {
		return errors.New(`"custom_https_provisioning_enabled" is set to "true". please add a "custom_https_configuration" block to the configuration file`)
	}

	return nil
}

func verifyCustomHttpsConfiguration(frontendEndpointCustomHttpsConfig []any) error {
	if len(frontendEndpointCustomHttpsConfig) > 0 {
		customHttpsConfiguration := frontendEndpointCustomHttpsConfig[0].(map[string]any)
		certificateSource := customHttpsConfiguration["certificate_source"].(string)
		certificateVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)

		if certificateSource == string(frontdoors.FrontDoorCertificateSourceFrontDoor) {
			if azureKeyVaultCertificateHasValues(customHttpsConfiguration, true) {
				return errors.New(`a Front Door managed "custom_https_configuration" block does not support the following keys. Please remove the following keys from your configuration file: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`)
			}
		} else {
			// The latest secret version is no longer valid for key vaults
			if strings.EqualFold(certificateVersion, "latest") {
				return errors.New(`"azure_key_vault_certificate_secret_version" can not be set to "latest" please remove this attribute from the configuration file. Removing the value has the same functionality as setting it to "latest"`)
			}

			if !azureKeyVaultCertificateHasValues(customHttpsConfiguration, false) {
				if certificateVersion == "" {
					// If using latest, empty string is now equivalent to using the keyword latest
					return errors.New(`a "AzureKeyVault" managed "custom_https_configuration" block must have values in the following fields: "azure_key_vault_certificate_secret_name" and "azure_key_vault_certificate_vault_id"`)
				} else {
					// If using a specific version of the secret
					return errors.New(`a "AzureKeyVault" managed "custom_https_configuration" block must have values in the following fields: "azure_key_vault_certificate_secret_name", "azure_key_vault_certificate_secret_version", and "azure_key_vault_certificate_vault_id"`)
				}
			}
		}
	}

	return nil
}

func azureKeyVaultCertificateHasValues(customHttpsConfiguration map[string]any, isFrontDoorManaged bool) bool {
	certificateSecretName := customHttpsConfiguration["azure_key_vault_certificate_secret_name"].(string)
	certificateSecretVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)
	certificateVaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"].(string)

	if isFrontDoorManaged {
		// if any of these keys have values it is invalid
		if strings.TrimSpace(certificateSecretName) != "" || strings.TrimSpace(certificateSecretVersion) != "" || strings.TrimSpace(certificateVaultId) != "" {
			return true
		}
	} else {
		if certificateSecretVersion == "" {
			// using latest ignore certificate secret version
			if strings.TrimSpace(certificateSecretName) != "" && strings.TrimSpace(certificateVaultId) != "" {
				return true
			}
		} else {
			// not using latest make sure all keys have values
			if strings.TrimSpace(certificateSecretName) != "" && strings.TrimSpace(certificateSecretVersion) != "" && strings.TrimSpace(certificateVaultId) != "" {
				return true
			}
		}
	}

	return false
}
