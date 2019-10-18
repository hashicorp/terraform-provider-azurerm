package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementApiPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementAPIPolicyCreateUpdate,
		Read:   resourceArmApiManagementAPIPolicyRead,
		Update: resourceArmApiManagementAPIPolicyCreateUpdate,
		Delete: resourceArmApiManagementAPIPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"api_name": azure.SchemaApiManagementChildName(),

			"xml_content": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"xml_link"},
				DiffSuppressFunc: suppress.XmlDiff,
			},

			"xml_link": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
			},
		},
	}
}

func resourceArmApiManagementAPIPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing API Policy (API Management Service %q / API %q / Resource Group %q): %s", serviceName, apiName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_policy", *existing.ID)
		}
	}

	parameters := apimanagement.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlContent != "" {
		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			ContentFormat: apimanagement.XML,
			PolicyContent: utils.String(xmlContent),
		}
	}

	if xmlLink != "" {
		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			ContentFormat: apimanagement.XMLLink,
			PolicyContent: utils.String(xmlLink),
		}
	}

	if parameters.PolicyContractProperties == nil {
		return fmt.Errorf("Either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName)
	if err != nil {
		return fmt.Errorf("Error retrieving API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementAPIPolicyRead(d, meta)
}

func resourceArmApiManagementAPIPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Policy (Resource Group %q / API Management Service %q / API %q) was not found - removing from state!", resourceGroup, serviceName, apiName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("api_name", apiName)

	if properties := resp.PolicyContractProperties; properties != nil {
		// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
		// as such there is no way to set `xml_link` and we'll let Terraform handle it
		d.Set("xml_content", properties.PolicyContent)
	}

	return nil
}

func resourceArmApiManagementAPIPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apiName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
		}
	}

	return nil
}
