package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/frontdoor/mgmt/2019-04-01/frontdoor"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmFrontDoor() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmFrontDoorCreateUpdate,
		Read:   resourceArmFrontDoorRead,
		Update: resourceArmFrontDoorCreateUpdate,
		Delete: resourceArmFrontDoorDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateFrontDoorName,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"enabled_state": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.EnabledStateEnabled),
					string(frontdoor.EnabledStateDisabled),
				}, false),
				Default: string(frontdoor.Enabled),
			},

			"enforce_certificate_name_check": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(frontdoor.EnforceCertificateNameCheckEnabledStateEnabled),
					string(frontdoor.EnforceCertificateNameCheckEnabledStateDisabled),
				}, false),
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmFrontDoorCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoorClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_front_door", *resp.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	enabledState := d.Get("enabled_state").(string)
	enforceCertificateNameCheck := d.Get("enforce_certificate_name_check").(string)
	friendlyName := d.Get("friendly_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := frontdoor.FrontDoor{
		Location: utils.String(location),
		Properties: &frontdoor.Properties{
			BackendPoolsSettings: &frontdoor.BackendPoolsSettings{
				EnforceCertificateNameCheck: frontdoor.EnforceCertificateNameCheckEnabledState(enforceCertificateNameCheck),
			},
			EnabledState: frontdoor.EnabledState(enabledState),
			FriendlyName: utils.String(friendlyName),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Front Door %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmFrontDoorRead(d, meta)
}

func resourceArmFrontDoorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoorClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontDoors"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	if properties := resp.Properties; properties != nil {
		d.Set("enabled_state", string(properties.EnabledState))
		if backendPoolsSettings := properties.BackendPoolsSettings; backendPoolsSettings != nil {
			d.Set("enforce_certificate_name_check", string(backendPoolsSettings.EnforceCertificateNameCheck))
		}
		d.Set("friendly_name", properties.FriendlyName)
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmFrontDoorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).frontdoorClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["frontDoors"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Front Door %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
