package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementGroupCreateUpdate,
		Read:   resourceArmApiManagementGroupRead,
		Update: resourceArmApiManagementGroupCreateUpdate,
		Delete: resourceArmApiManagementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementChildName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"external_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(apimanagement.Custom),
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.Custom),
					string(apimanagement.External),
				}, false),
			},
		},
	}
}

func resourceArmApiManagementGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	externalID := d.Get("external_id").(string)
	groupType := d.Get("type").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Group %q (API Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_group", *existing.ID)
		}
	}

	parameters := apimanagement.GroupCreateParameters{
		GroupCreateParametersProperties: &apimanagement.GroupCreateParametersProperties{
			DisplayName: utils.String(displayName),
			Description: utils.String(description),
			ExternalID:  utils.String(externalID),
			Type:        apimanagement.GroupType(groupType),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating Group %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Group %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Group %q (Resource Group %q / API Management Service %q)", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementGroupRead(d, meta)
}

func resourceArmApiManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["groups"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Group %q (Resource Group %q / API Management Service %q) was not found - removing from state!", name, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for Group %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if properties := resp.GroupContractProperties; properties != nil {
		d.Set("display_name", properties.DisplayName)
		d.Set("description", properties.Description)
		d.Set("external_id", properties.ExternalID)
		d.Set("type", string(properties.Type))
	}

	return nil
}

func resourceArmApiManagementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.GroupClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["groups"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Group %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}
