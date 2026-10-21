// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_custom_https_configuration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorCustomHTTPSConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.FrontendEndpointsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Front Door Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("frontend_endpoint_id", id.ID())

	if model := resp.Model; model != nil {
		flattenedHttpsConfig := helpers.FlattenCustomHttpsConfiguration(model.Properties)
		if err := d.Set("custom_https_configuration", flattenedHttpsConfig.CustomHTTPSConfiguration); err != nil {
			return fmt.Errorf("setting `custom_https_configuration`: %+v", err)
		}
		if err := d.Set("custom_https_provisioning_enabled", flattenedHttpsConfig.CustomHTTPSProvisioningEnabled); err != nil {
			return fmt.Errorf("setting `custom_https_provisioning_enabled`: %+v", err)
		}
	}

	return nil
}

type customHttpsConfigurationUpdateInput struct {
	customHttpsConfigurationCurrent *frontdoors.CustomHTTPSConfiguration
	customHttpsConfigurationNew     []any
	customHttpsProvisioningEnabled  bool
	frontendEndpointId              frontdoors.FrontendEndpointId
	provisioningState               frontdoors.CustomHTTPSProvisioningState
}

func updateCustomHTTPSConfiguration(ctx context.Context, client *frontdoors.FrontDoorsClient, input customHttpsConfigurationUpdateInput) error {
	// Locking to prevent parallel changes causing issues
	frontendEndpointResourceId := input.frontendEndpointId.ID()
	locks.ByID(frontendEndpointResourceId)
	defer locks.UnlockByID(frontendEndpointResourceId)

	if input.provisioningState == "" {
		return nil
	}

	// Check to see if we are going to change the CustomHTTPSProvisioningState, if so check to
	// see if its current state is configurable, if not return an error...
	if input.customHttpsProvisioningEnabled != NormalizeCustomHTTPSProvisioningStateToBool(input.provisioningState) {
		if err := isFrontDoorFrontendEndpointConfigurable(input.provisioningState, input.customHttpsProvisioningEnabled, input.frontendEndpointId); err != nil {
			return err
		}
	}

	if input.customHttpsProvisioningEnabled {
		if len(input.customHttpsConfigurationNew) > 0 && input.customHttpsConfigurationNew[0] != nil {
			customHTTPSConfiguration := input.customHttpsConfigurationNew[0].(map[string]any)
			minTLSVersion := frontdoors.MinimumTLSVersionOnePointTwo // Default to TLS 1.2
			if httpsConfig := input.customHttpsConfigurationCurrent; httpsConfig != nil {
				minTLSVersion = httpsConfig.MinimumTlsVersion
			}
			customHTTPSConfigurationUpdate := makeCustomHTTPSConfiguration(customHTTPSConfiguration, minTLSVersion)
			if input.provisioningState == frontdoors.CustomHTTPSProvisioningStateDisabled || customHTTPSConfigurationUpdate != *input.customHttpsConfigurationCurrent {
				// Enable Custom Domain HTTPS for the Frontend Endpoint
				if err := resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, client, input.frontendEndpointId, true, customHTTPSConfigurationUpdate); err != nil {
					return fmt.Errorf("unable to enable/update Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", input.frontendEndpointId.FrontendEndpointName, input.frontendEndpointId.ResourceGroupName, err)
				}
			}
		}
	} else if !input.customHttpsProvisioningEnabled && input.provisioningState == frontdoors.CustomHTTPSProvisioningStateEnabled {
		// Disable Custom Domain HTTPS for the Frontend Endpoint
		if err := resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, client, input.frontendEndpointId, false, frontdoors.CustomHTTPSConfiguration{}); err != nil {
			return fmt.Errorf("unable to disable Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", input.frontendEndpointId.FrontendEndpointName, input.frontendEndpointId.ResourceGroupName, err)
		}
	}

	return nil
}

func resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx context.Context, client *frontdoors.FrontDoorsClient, id frontdoors.FrontendEndpointId, enableCustomHTTPSProvisioning bool, customHTTPSConfiguration frontdoors.CustomHTTPSConfiguration) error {
	if enableCustomHTTPSProvisioning {
		if err := client.FrontendEndpointsEnableHTTPSThenPoll(ctx, id, customHTTPSConfiguration); err != nil {
			return fmt.Errorf("enabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		return nil
	}

	if err := client.FrontendEndpointsDisableHTTPSThenPoll(ctx, id); err != nil {
		return fmt.Errorf("disabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
	}

	return nil
}

func makeCustomHTTPSConfiguration(customHttpsConfiguration map[string]any, minTLSVersion frontdoors.MinimumTLSVersion) frontdoors.CustomHTTPSConfiguration {
	// https://github.com/Azure/azure-sdk-for-go/issues/6882

	customHTTPSConfigurationUpdate := frontdoors.CustomHTTPSConfiguration{
		ProtocolType:      frontdoors.FrontDoorTlsProtocolType("ServerNameIndication"),
		MinimumTlsVersion: minTLSVersion,
	}

	if customHttpsConfiguration["certificate_source"].(string) == "AzureKeyVault" {
		vaultSecret := customHttpsConfiguration["azure_key_vault_certificate_secret_name"].(string)
		vaultVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)
		vaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"].(string)

		customHTTPSConfigurationUpdate.CertificateSource = frontdoors.FrontDoorCertificateSourceAzureKeyVault
		customHTTPSConfigurationUpdate.KeyVaultCertificateSourceParameters = &frontdoors.KeyVaultCertificateSourceParameters{
			Vault: &frontdoors.KeyVaultCertificateSourceParametersVault{
				Id: pointer.To(vaultId),
			},
			SecretName:    pointer.To(vaultSecret),
			SecretVersion: pointer.To(vaultVersion),
		}
	} else {
		customHTTPSConfigurationUpdate.CertificateSource = frontdoors.FrontDoorCertificateSourceFrontDoor
		customHTTPSConfigurationUpdate.FrontDoorCertificateSourceParameters = &frontdoors.FrontDoorCertificateSourceParameters{
			CertificateType: pointer.To(frontdoors.FrontDoorCertificateTypeDedicated),
		}
	}

	return customHTTPSConfigurationUpdate
}
