// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontDoorCustomHTTPSConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCustomHTTPSConfigurationCreateUpdate,
		Read:   resourceFrontDoorCustomHTTPSConfigurationRead,
		Update: resourceFrontDoorCustomHTTPSConfigurationCreateUpdate,
		Delete: resourceFrontDoorCustomHTTPSConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.CustomHttpsConfigurationID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Frontdoor.FrontDoorsClient

			// validate that the passed ID is a valid custom HTTPS configuration ID
			custom, err := parse.CustomHttpsConfigurationID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing Custom HTTPS Configuration ID %q for import: %v", d.Id(), err)
			}

			// convert the passed custom HTTPS configuration ID to a frontend endpoint ID
			frontend := frontdoors.NewFrontendEndpointID(custom.SubscriptionId, custom.ResourceGroup, custom.FrontDoorName, custom.CustomHttpsConfigurationName)

			// validate that the frontend endpoint ID exists in the Frontdoor resource
			if _, err = client.FrontendEndpointsGet(ctx, frontend); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving the Custom HTTPS Configuration(ID: %q) for the frontend endpoint (ID: %q): %s", custom.ID(), frontend.ID(), err)
			}

			// set the new values for the custom HTTPS configuration resource
			d.Set("id", custom.ID())
			d.Set("frontend_endpoint_id", frontend.ID())

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"frontend_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontendEndpointID,
			},

			"custom_https_provisioning_enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"custom_https_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: schemaCustomHttpsConfiguration(),
				},
			},
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(customizeHttpsConfigurationCustomizeDiff),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CustomHttpsConfigurationV0ToV1{},
		}),
	}
}

func resourceFrontDoorCustomHTTPSConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	customHttpsConfigurationId := parse.NewCustomHttpsConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.FrontDoorName, id.FrontendEndpointName)

	resp, err := client.FrontendEndpointsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("reading %s: `properties` was nil", id)
	}
	props := *resp.Model.Properties

	input := customHttpsConfigurationUpdateInput{
		customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
		customHttpsConfigurationNew:     d.Get("custom_https_configuration").([]interface{}),
		customHttpsProvisioningEnabled:  d.Get("custom_https_provisioning_enabled").(bool),
		frontendEndpointId:              *id,
	}

	if props.CustomHTTPSProvisioningState != nil {
		input.provisioningState = *props.CustomHTTPSProvisioningState
	}

	if err := updateCustomHTTPSConfiguration(ctx, client, input); err != nil {
		return fmt.Errorf("updating Custom HTTPS configuration for %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(customHttpsConfigurationId.ID())
	}

	return resourceFrontDoorCustomHTTPSConfigurationRead(d, meta)
}

func resourceFrontDoorCustomHTTPSConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		flattenedHttpsConfig := flattenCustomHttpsConfiguration(model.Properties)
		if err := d.Set("custom_https_configuration", flattenedHttpsConfig.CustomHTTPSConfiguration); err != nil {
			return fmt.Errorf("setting `custom_https_configuration`: %+v", err)
		}
		if err := d.Set("custom_https_provisioning_enabled", flattenedHttpsConfig.CustomHTTPSProvisioningEnabled); err != nil {
			return fmt.Errorf("setting `custom_https_provisioning_enabled`: %+v", err)
		}
	}

	return nil
}

func resourceFrontDoorCustomHTTPSConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := frontdoors.ParseFrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.FrontendEndpointsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if model.Properties == nil {
			return fmt.Errorf("reading %s: `properties` was nil", id)
		}
		props := *model.Properties

		input := customHttpsConfigurationUpdateInput{
			customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
			customHttpsConfigurationNew:     make([]interface{}, 0),
			customHttpsProvisioningEnabled:  false,
			frontendEndpointId:              *id,
		}

		if props.CustomHTTPSProvisioningState != nil {
			input.provisioningState = *props.CustomHTTPSProvisioningState
		}
		if err := updateCustomHTTPSConfiguration(ctx, client, input); err != nil {
			return fmt.Errorf("disabling Custom HTTPS configuration for Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", id.FrontendEndpointName, id.FrontDoorName, id.ResourceGroupName, err)
		}
	}

	return nil
}

type customHttpsConfigurationUpdateInput struct {
	customHttpsConfigurationCurrent *frontdoors.CustomHTTPSConfiguration
	customHttpsConfigurationNew     []interface{}
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
			customHTTPSConfiguration := input.customHttpsConfigurationNew[0].(map[string]interface{})
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

func makeCustomHTTPSConfiguration(customHttpsConfiguration map[string]interface{}, minTLSVersion frontdoors.MinimumTLSVersion) frontdoors.CustomHTTPSConfiguration {
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
				Id: utils.String(vaultId),
			},
			SecretName:    utils.String(vaultSecret),
			SecretVersion: utils.String(vaultVersion),
		}
	} else {
		customHTTPSConfigurationUpdate.CertificateSource = frontdoors.FrontDoorCertificateSourceFrontDoor
		certificateType := frontdoors.FrontDoorCertificateTypeDedicated
		customHTTPSConfigurationUpdate.FrontDoorCertificateSourceParameters = &frontdoors.FrontDoorCertificateSourceParameters{
			CertificateType: &certificateType,
		}
	}

	return customHTTPSConfigurationUpdate
}
