package frontdoor

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2019-11-01/frontdoor"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"time"
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
			"front_door_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"frontend_endpoint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorBackendPoolRoutingRuleName,
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
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func resourceArmFrontDoorCustomHttpsConfigurationCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	frontDoorName := d.Get("front_door_name").(string)
	frontendEndpointName := d.Get("frontend_endpoint_name").(string)

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, frontendEndpointName)
	if err != nil {
		return fmt.Errorf("reading Front Door Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	customHttpsProvisioningEnabled := d.Get("custom_https_provisioning_enabled").(bool)
	customHttpsConfigurationNew := d.Get("custom_https_configuration").([]interface{})
	err = resourceArmFrontDoorFrontendEndpointCustomHttpsConfigurationUpdate(d, customHttpsProvisioningEnabled, frontDoorName, frontendEndpointName, resourceGroup, resp.CustomHTTPSProvisioningState, resp.CustomHTTPSConfiguration, customHttpsConfigurationNew, meta)
	if err != nil {
		return fmt.Errorf("Unable to update Custom HTTPS configuration for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	return nil
}

func resourceArmFrontDoorCustomHttpsConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err, resourceGroup, frontDoorName, frontendEndpointName := frontDoorCustomHttpsConfigurationReadParams(d)
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

	d.Set("front_door_name", frontDoorName)
	d.Set("frontend_endpoint_name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if resp.Name != nil {
		if frontDoorFrontendEndpoint, err := flattenArmFrontDoorFrontendEndpoint(d, &resp, resourceGroup, *resp.Name, meta); frontDoorFrontendEndpoint != nil {
			if err := d.Set("frontend_endpoint", frontDoorFrontendEndpoint); err != nil {
				return fmt.Errorf("setting `frontend_endpoint`: %+v", err)
			}
		} else {
			return fmt.Errorf("flattening `frontend_endpoint`: %+v", err)
		}
	} else {
		return fmt.Errorf("flattening `frontend_endpoint`: Unable to read Frontdoor Name")
	}

	return nil
}

func resourceArmFrontDoorCustomHttpsConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsFrontendClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err, resourceGroup, frontDoorName, frontendEndpointName := frontDoorCustomHttpsConfigurationReadParams(d)
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
	err = resourceArmFrontDoorFrontendEndpointCustomHttpsConfigurationUpdate(d, false, frontDoorName, frontendEndpointName, resourceGroup, resp.CustomHTTPSProvisioningState, resp.CustomHTTPSConfiguration, customHttpsConfigurationNew, meta)
	if err != nil {
		return fmt.Errorf("unable to disable Custom HTTPS configuration for Frontend Endpoint %q (Resource Group %q): %+v", frontendEndpointName, resourceGroup, err)
	}

	return nil
}

func frontDoorCustomHttpsConfigurationReadParams(d *schema.ResourceData) (err error, resourceGroup string, frontDoorName string, frontendEndpointName string) {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err, "", "", ""
	}
	resourceGroup = id.ResourceGroup
	frontDoorName = id.Path["frontdoors"]
	// Link to issue: https://github.com/Azure/azure-sdk-for-go/issues/6762
	if frontDoorName == "" {
		frontDoorName = id.Path["Frontdoors"]
	}
	frontendEndpointName = id.Path["frontendEndpoints"]
	if frontendEndpointName == "" {
		frontDoorName = id.Path["FrontendEndpoints"]
	}

	return nil, resourceGroup, frontDoorName, frontendEndpointName
}
