package frontdoor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-01-01/frontdoor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFrontDoorCustomHttpsConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorCustomHttpsConfigurationCreateUpdate,
		Read:   resourceArmFrontDoorCustomHttpsConfigurationRead,
		Update: resourceArmFrontDoorCustomHttpsConfigurationCreateUpdate,
		Delete: resourceArmFrontDoorCustomHttpsConfigurationDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"frontend_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"custom_https_provisioning_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"custom_https_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: SchemaFrontdoorCustomHttpsConfiguration(),
				},
			},
		},

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if err := validate.FrontdoorCustomHttpsSettings(d); err != nil {
				return fmt.Errorf("creating Front Door Custom Https Configuration for endpoint %q (Resource Group %q): %+v", d.Get("frontend_endpoint_id").(string), d.Get("resource_group_name").(string), err)
			}

			return nil
		},
	}
}

func resourceArmFrontDoorCustomHttpsConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err, resourceGroup, frontDoorName, frontendEndpointName := frontDoorFrontendEndpointReadProps(d)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
	if err != nil {
		return fmt.Errorf("reading Front Door Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	// This is because azure doesn't have an 'id' for a custom https configuration
	// In order to compensate for this and allow importing of this resource we are artificially
	// creating an identity for a custom https configuration object
	var resourceId string
	if id := resp.ID; id != nil && *id != "" {
		resourceId = fmt.Sprintf("%s/customHttpsConfiguration/%s", *id, frontendEndpointName)
	} else {
		return fmt.Errorf("unable to retrieve Front Door Endpoint %q (Resource Group %q) ID", frontendEndpointName, resourceGroup)
	}

	customHttpsProvisioningEnabled := d.Get("custom_https_provisioning_enabled").(bool)
	customHttpsConfigurationNew := d.Get("custom_https_configuration").([]interface{})
	err = resourceArmFrontDoorFrontendEndpointCustomHttpsConfigurationUpdate(ctx, *resp.ID, customHttpsProvisioningEnabled, frontDoorName, frontendEndpointName, resourceGroup, resp.CustomHTTPSProvisioningState, resp.CustomHTTPSConfiguration, customHttpsConfigurationNew, meta)
	if err != nil {
		return fmt.Errorf("unable to update Custom HTTPS configuration for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	resp, err = client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
	if err != nil {
		return fmt.Errorf("retreving Front Door Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Front Door Endpoint %q (Resource Group %q) ID", frontendEndpointName, resourceGroup)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	return nil
}

func resourceArmFrontDoorCustomHttpsConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err, resourceGroup, frontDoorName, frontendEndpointName := frontDoorCustomHttpsConfigurationReadProps(d)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Front Door Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	d.Set("frontend_endpoint_id", resp.ID)

	if resp.Name != nil {

		flattenedHttpsConfig := FlattenArmFrontDoorCustomHttpsConfiguration(*resp.FrontendEndpointProperties)
		if err := d.Set("custom_https_configuration", flattenedHttpsConfig.CustomHTTPSConfiguration); err != nil {
			return fmt.Errorf("setting `custom_https_configuration`: %+v", err)
		}
		if err := d.Set("custom_https_provisioning_enabled", flattenedHttpsConfig.CustomHTTPSProvisioningEnabled); err != nil {
			return fmt.Errorf("setting `custom_https_provisioning_enabled`: %+v", err)
		}
	} else {
		return fmt.Errorf("flattening `frontend_endpoint` `custom_https_configuration`: Unable to read Frontend Endpoint Name")
	}

	return nil
}

func resourceArmFrontDoorCustomHttpsConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err, resourceGroup, frontDoorName, frontendEndpointName := frontDoorCustomHttpsConfigurationReadProps(d)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("reading Front Door Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	customHttpsConfigurationNew := make([]interface{}, 0)
	err = resourceArmFrontDoorFrontendEndpointCustomHttpsConfigurationUpdate(ctx, *resp.ID, false, frontDoorName, frontendEndpointName, resourceGroup, resp.CustomHTTPSProvisioningState, resp.CustomHTTPSConfiguration, customHttpsConfigurationNew, meta)
	if err != nil {
		return fmt.Errorf("unable to disable Custom HTTPS configuration for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	return nil
}

func frontDoorFrontendEndpointReadProps(d *schema.ResourceData) (err error, resourceGroup string, frontDoorName string, frontendEndpointName string) {
	id, err := azure.ParseAzureResourceID(d.Get("frontend_endpoint_id").(string))
	if err != nil {
		return err, "", "", ""
	}
	return frontDoorReadPropsFromId(id)
}

func frontDoorCustomHttpsConfigurationReadProps(d *schema.ResourceData) (err error, resourceGroup string, frontDoorName string, frontendEndpointName string) {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err, "", "", ""
	}
	return frontDoorReadPropsFromId(id)
}

func frontDoorReadPropsFromId(id *azure.ResourceID) (err error, resourceGroup string, frontDoorName string, frontendEndpointName string) {
	resourceGroup = id.ResourceGroup
	frontDoorName = id.Path["frontdoors"]
	// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
	if frontDoorName == "" {
		frontDoorName = id.Path["Frontdoors"]
	}
	frontendEndpointName = id.Path["frontendendpoints"]
	if frontendEndpointName == "" {
		frontDoorName = id.Path["FrontendEndpoints"]
	}

	return nil, resourceGroup, frontDoorName, frontendEndpointName
}

func resourceArmFrontDoorFrontendEndpointCustomHttpsConfigurationUpdate(ctx context.Context, frontendEndpointId string, customHttpsProvisioningEnabled bool, frontDoorName string, frontendEndpointName string, resourceGroup string, provisioningState frontdoor.CustomHTTPSProvisioningState, customHTTPSConfigurationCurrent *frontdoor.CustomHTTPSConfiguration, customHttpsConfigurationNew []interface{}, meta interface{}) error {
	// Locking to prevent parallel changes causing issues
	locks.ByID(frontendEndpointId)
	defer locks.UnlockByID(frontendEndpointId)

	if provisioningState != "" {
		// Check to see if we are going to change the CustomHTTPSProvisioningState, if so check to
		// see if its current state is configurable, if not return an error...
		if customHttpsProvisioningEnabled != NormalizeCustomHTTPSProvisioningStateToBool(provisioningState) {
			if err := IsFrontDoorFrontendEndpointConfigurable(provisioningState, customHttpsProvisioningEnabled, frontendEndpointName, resourceGroup); err != nil {
				return err
			}
		}

		if customHttpsProvisioningEnabled {
			// Build a custom Https configuration based off the config file to send to the enable call
			// NOTE: I do not need to check to see if this exists since I already do that in the validation code
			customHTTPSConfiguration := customHttpsConfigurationNew[0].(map[string]interface{})
			minTLSVersion := frontdoor.OneFullStopTwo // Default to TLS 1.2
			if httpsConfig := customHTTPSConfigurationCurrent; httpsConfig != nil {
				minTLSVersion = httpsConfig.MinimumTLSVersion
			}
			customHTTPSConfigurationUpdate := makeCustomHttpsConfiguration(customHTTPSConfiguration, minTLSVersion)
			if provisioningState == frontdoor.CustomHTTPSProvisioningStateDisabled || customHTTPSConfigurationUpdate != *customHTTPSConfigurationCurrent {
				// Enable Custom Domain HTTPS for the Frontend Endpoint
				if err := resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, true, frontDoorName, frontendEndpointName, resourceGroup, customHTTPSConfigurationUpdate, meta); err != nil {
					return fmt.Errorf("unable to enable/update Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
				}
			}
		} else if !customHttpsProvisioningEnabled && provisioningState == frontdoor.CustomHTTPSProvisioningStateEnabled {
			// Disable Custom Domain HTTPS for the Frontend Endpoint
			if err := resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx, false, frontDoorName, frontendEndpointName, resourceGroup, frontdoor.CustomHTTPSConfiguration{}, meta); err != nil {
				return fmt.Errorf("unable to disable Custom Domain HTTPS for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
			}
		}
	}

	return nil
}

func resourceArmFrontDoorFrontendEndpointEnableHttpsProvisioning(ctx context.Context, enableCustomHttpsProvisioning bool, frontDoorName string, frontendEndpointName string, resourceGroup string, customHTTPSConfiguration frontdoor.CustomHTTPSConfiguration, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient

	if enableCustomHttpsProvisioning {
		future, err := client.EnableHTTPS(ctx, resourceGroup, frontDoorName, frontendEndpointName, customHTTPSConfiguration)

		if err != nil {
			return fmt.Errorf("enabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting to enable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
	} else {
		future, err := client.DisableHTTPS(ctx, resourceGroup, frontDoorName, frontendEndpointName)

		if err != nil {
			return fmt.Errorf("disabling Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			// If the endpoint does not exist but this is not a new resource, the custom https
			// configuration which previously existed was deleted with the endpoint, so reflect
			// that in state.
			resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return nil
				}
			}
			return fmt.Errorf("waiting to disable Custom Domain HTTPS for Frontend Endpoint: %+v", err)
		}
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
