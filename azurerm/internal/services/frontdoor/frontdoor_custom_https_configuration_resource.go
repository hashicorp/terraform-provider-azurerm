package frontdoor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceFrontDoorCustomHttpsConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorCustomHttpsConfigurationCreateUpdate,
		Read:   resourceFrontDoorCustomHttpsConfigurationRead,
		Update: resourceFrontDoorCustomHttpsConfigurationCreateUpdate,
		Delete: resourceFrontDoorCustomHttpsConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.CustomHttpsConfigurationID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient

			// validate that the passed ID is a valid custom HTTPS configuration ID
			custom, err := parse.CustomHttpsConfigurationID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing Custom HTTPS Configuration ID %q for import: %v", d.Id(), err)
			}

			// convert the passed custom HTTPS configuration ID to a frontend endpoint ID
			frontend := parse.NewFrontendEndpointID(custom.SubscriptionId, custom.ResourceGroup, custom.FrontDoorName, custom.CustomHttpsConfigurationName)

			// validate that the frontend endpoint ID exists in the Frontdoor resource
			if _, err = client.Get(ctx, custom.ResourceGroup, custom.FrontDoorName, custom.CustomHttpsConfigurationName); err != nil {
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

func resourceFrontDoorCustomHttpsConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontendEndpointId, err := parse.FrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	customHttpsConfigurationId := parse.NewCustomHttpsConfigurationID(frontendEndpointId.SubscriptionId, frontendEndpointId.ResourceGroup, frontendEndpointId.FrontDoorName, frontendEndpointId.Name)

	resp, err := client.Get(ctx, frontendEndpointId.ResourceGroup, frontendEndpointId.FrontDoorName, frontendEndpointId.Name)
	if err != nil {
		return fmt.Errorf("reading Endpoint %q (Front Door %q / Resource Group %q): %+v", frontendEndpointId.Name, frontendEndpointId.FrontDoorName, frontendEndpointId.ResourceGroup, err)
	}

	if resp.FrontendEndpointProperties == nil {
		return fmt.Errorf("reading Endpoint %q (Front Door %q / Resource Group %q): `properties` was nil", frontendEndpointId.Name, frontendEndpointId.FrontDoorName, frontendEndpointId.ResourceGroup)
	}
	props := *resp.FrontendEndpointProperties

	input := customHttpsConfigurationUpdateInput{
		customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
		customHttpsConfigurationNew:     d.Get("custom_https_configuration").([]interface{}),
		customHttpsProvisioningEnabled:  d.Get("custom_https_provisioning_enabled").(bool),
		frontendEndpointId:              *frontendEndpointId,
		provisioningState:               props.CustomHTTPSProvisioningState,
	}

	if err := updateCustomHttpsConfiguration(ctx, client, input); err != nil {
		return fmt.Errorf("updating Custom HTTPS configuration for Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", frontendEndpointId.Name, frontendEndpointId.FrontDoorName, frontendEndpointId.ResourceGroup, err)
	}

	if d.IsNewResource() {
		d.SetId(customHttpsConfigurationId.ID())
	}

	return resourceFrontDoorCustomHttpsConfigurationRead(d, meta)
}

func resourceFrontDoorCustomHttpsConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Front Door Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("frontend_endpoint_id", id.ID())

	flattenedHttpsConfig := flattenCustomHttpsConfiguration(resp.FrontendEndpointProperties)
	if err := d.Set("custom_https_configuration", flattenedHttpsConfig.CustomHTTPSConfiguration); err != nil {
		return fmt.Errorf("setting `custom_https_configuration`: %+v", err)
	}
	if err := d.Set("custom_https_provisioning_enabled", flattenedHttpsConfig.CustomHTTPSProvisioningEnabled); err != nil {
		return fmt.Errorf("setting `custom_https_provisioning_enabled`: %+v", err)
	}

	return nil
}

func resourceFrontDoorCustomHttpsConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontendEndpointIDInsensitively(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("reading Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", id.Name, id.FrontDoorName, id.ResourceGroup, err)
	}

	if resp.FrontendEndpointProperties == nil {
		return fmt.Errorf("reading Frontend Endpoint %q (Front Door %q / Resource Group %q): `properties` was nil", id.Name, id.FrontDoorName, id.ResourceGroup)
	}
	props := *resp.FrontendEndpointProperties

	input := customHttpsConfigurationUpdateInput{
		customHttpsConfigurationCurrent: props.CustomHTTPSConfiguration,
		customHttpsConfigurationNew:     make([]interface{}, 0),
		customHttpsProvisioningEnabled:  false,
		frontendEndpointId:              *id,
		provisioningState:               props.CustomHTTPSProvisioningState,
	}
	if err := updateCustomHttpsConfiguration(ctx, client, input); err != nil {
		return fmt.Errorf("disabling Custom HTTPS configuration for Frontend Endpoint %q (Front Door %q / Resource Group %q): %+v", id.Name, id.FrontDoorName, id.ResourceGroup, err)
	}

	return nil
}

type customHttpsConfigurationUpdateInput struct {
	customHttpsConfigurationCurrent *frontdoor.CustomHTTPSConfiguration
	customHttpsConfigurationNew     []interface{}
	customHttpsProvisioningEnabled  bool
	frontendEndpointId              parse.FrontendEndpointId
	provisioningState               frontdoor.CustomHTTPSProvisioningState
}

func updateCustomHttpsConfiguration(ctx context.Context, client *frontdoor.FrontendEndpointsClient, input customHttpsConfigurationUpdateInput) error {
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
			minTLSVersion := frontdoor.OneFullStopTwo // Default to TLS 1.2
			if httpsConfig := input.customHttpsConfigurationCurrent; httpsConfig != nil {
				minTLSVersion = httpsConfig.MinimumTLSVersion
			}
			customHTTPSConfigurationUpdate := makeCustomHttpsConfiguration(customHTTPSConfiguration, minTLSVersion)
			if input.provisioningState == frontdoor.CustomHTTPSProvisioningStateDisabled || customHTTPSConfigurationUpdate != *input.customHttpsConfigurationCurrent {
				// Enable Custom Domain HTTPS for the Frontend Endpoint
				if err := resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, client, input.frontendEndpointId, true, customHTTPSConfigurationUpdate); err != nil {
					return fmt.Errorf("unable to enable/update Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", input.frontendEndpointId.Name, input.frontendEndpointId.ResourceGroup, err)
				}
			}
		}
	} else if !input.customHttpsProvisioningEnabled && input.provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled {
		// Disable Custom Domain HTTPS for the Frontend Endpoint
		if err := resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, client, input.frontendEndpointId, false, frontdoor.CustomHTTPSConfiguration{}); err != nil {
			return fmt.Errorf("unable to disable Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", input.frontendEndpointId.Name, input.frontendEndpointId.ResourceGroup, err)
		}
	}

	return nil
}

func resourceFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx context.Context, client *frontdoor.FrontendEndpointsClient, id parse.FrontendEndpointId, enableCustomHttpsProvisioning bool, customHTTPSConfiguration frontdoor.CustomHTTPSConfiguration) error {
	if enableCustomHttpsProvisioning {
		future, err := client.EnableHTTPS(ctx, id.ResourceGroup, id.FrontDoorName, id.Name, customHTTPSConfiguration)
		if err != nil {
			return fmt.Errorf("enabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to enable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}

		return nil
	}

	future, err := client.DisableHTTPS(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
	if err != nil {
		return fmt.Errorf("disabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		// If the endpoint does not exist but this is not a new resource, the custom https
		// configuration which previously existed was deleted with the endpoint, so reflect
		// that in state.
		resp, err := client.Get(ctx, id.ResourceGroup, id.FrontDoorName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
		}
		return fmt.Errorf("waiting to disable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
	}

	return nil
}

func makeCustomHttpsConfiguration(customHttpsConfiguration map[string]interface{}, minTLSVersion frontdoor.MinimumTLSVersion) frontdoor.CustomHTTPSConfiguration {
	// https://github.com/Azure/azure-sdk-for-go/issues/6882
	defaultProtocolType := "ServerNameIndication"

	customHTTPSConfigurationUpdate := frontdoor.CustomHTTPSConfiguration{
		ProtocolType:      &defaultProtocolType,
		MinimumTLSVersion: minTLSVersion,
	}

	if customHttpsConfiguration["certificate_source"].(string) == "AzureKeyVault" {
		vaultSecret := customHttpsConfiguration["azure_key_vault_certificate_secret_name"].(string)
		vaultVersion := customHttpsConfiguration["azure_key_vault_certificate_secret_version"].(string)
		vaultId := customHttpsConfiguration["azure_key_vault_certificate_vault_id"].(string)

		customHTTPSConfigurationUpdate.CertificateSource = frontdoor.CertificateSourceAzureKeyVault
		customHTTPSConfigurationUpdate.KeyVaultCertificateSourceParameters = &frontdoor.KeyVaultCertificateSourceParameters{
			Vault: &frontdoor.KeyVaultCertificateSourceParametersVault{
				ID: utils.String(vaultId),
			},
			SecretName:    utils.String(vaultSecret),
			SecretVersion: utils.String(vaultVersion),
		}
	} else {
		customHTTPSConfigurationUpdate.CertificateSource = frontdoor.CertificateSourceFrontDoor
		customHTTPSConfigurationUpdate.CertificateSourceParameters = &frontdoor.CertificateSourceParameters{
			CertificateType: frontdoor.Dedicated,
		}
	}

	return customHTTPSConfigurationUpdate
}
