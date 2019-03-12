package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementApiVersionSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementApiVersionSetCreateUpdate,
		Read:   resourceArmApiManagementApiVersionSetRead,
		Update: resourceArmApiManagementApiVersionSetCreateUpdate,
		Delete: resourceArmApiManagementApiVersionSetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementChildName(),

			"resource_group_name": resourceGroupNameSchema(),

			"api_management_name": azure.SchemaApiManagementName(),

			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"versioning_schema": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.VersioningSchemeHeader),
					string(apimanagement.VersioningSchemeQuery),
					string(apimanagement.VersioningSchemeSegment),
				}, false),
			},

			"version_header_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"version_query_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmApiManagementApiVersionSetCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiVersionSetClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Api Version Set %q (Api Management Service %q / Resource Group %q): %s", name, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_version_set", *existing.ID)
		}
	}

	parameters := apimanagement.APIVersionSetContract{
		APIVersionSetContractProperties: &apimanagement.APIVersionSetContractProperties{
			DisplayName:      utils.String(d.Get("display_name").(string)),
			VersioningScheme: apimanagement.VersioningScheme(d.Get("versioning_schema").(string)),
			Description:      utils.String(d.Get("description").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, name, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Api Version Set %q (Resource Group %q / Api Management Service %q)", name, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementApiVersionSetRead(d, meta)
}

func resourceArmApiManagementApiVersionSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiVersionSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["api-version-sets"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Api Version Set %q (Resource Group %q / Api Management Service %q) was not found - removing from state!", name, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.APIVersionSetContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("versioning_schema", props.VersioningScheme)
		d.Set("version_header_name", props.VersionHeaderName)
		d.Set("version_query_name", props.VersionQueryName)
	}

	return nil
}

func resourceArmApiManagementApiVersionSetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiVersionSetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	name := id.Path["api-version-sets"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Api Version Set %q (Resource Group %q / Api Management Service %q): %+v", name, resourceGroup, serviceName, err)
		}
	}

	return nil
}
