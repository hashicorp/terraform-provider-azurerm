package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementPolicyCreateUpdate,
		Read:   resourceArmApiManagementPolicyRead,
		Update: resourceArmApiManagementPolicyCreateUpdate,
		Delete: resourceArmApiManagementPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameSchema(),

			"api_management_name": azure.SchemaApiManagementName(),

			"xml_content": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.SuppressXmlDiff,
				ConflictsWith:    []string{"xml_link"},
			},

			"xml_link": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
			},
		},
	}
}

func resourceArmApiManagementPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementPolicyClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)

	if requireResourcesToBeImported {
		resp, err := client.Get(ctx, resourceGroup, serviceName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Global Policy (API Management Service %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_policy", *resp.ID)
		}
	}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)
	if xmlContent == "" && xmlLink == "" {
		return fmt.Errorf("Either `xml_content` or `xml_link` is required")
	}

	content := xmlContent
	format := apimanagement.XML
	if xmlLink != "" {
		content = xmlLink
		format = apimanagement.XMLLink
	}

	parameters := apimanagement.PolicyContract{
		PolicyContractProperties: &apimanagement.PolicyContractProperties{
			ContentFormat: format,
			PolicyContent: utils.String(content),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, parameters); err != nil {
		return fmt.Errorf("Error creating Global Policy (API Management Service %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		return fmt.Errorf("Error retrieving Global Policy (API Management Service %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Global Policy (API Management Service %q / Resource Group %q) ID", serviceName, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementPolicyRead(d, meta)
}

func resourceArmApiManagementPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementPolicyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]

	resp, err := client.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Api Management Policy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Global Policy (API Management Service %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	if properties := resp.PolicyContractProperties; properties != nil {
		d.Set("xml_content", "")
		d.Set("xml_link", "")
		if properties.ContentFormat == apimanagement.XML {
			d.Set("xml_content", properties.PolicyContent)
		} else if properties.ContentFormat == apimanagement.XMLLink {
			d.Set("xml_link", properties.PolicyContent)
		}
	}

	return nil
}

func resourceArmApiManagementPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementPolicyClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]

	if _, err := client.Delete(ctx, resourceGroup, serviceName, ""); err != nil {
		return fmt.Errorf("Error deleting Global Policy (API Management Service %q / Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	return nil
}
