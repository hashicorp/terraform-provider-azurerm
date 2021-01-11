package apimanagement

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementAPIPolicyCreateUpdate,
		Read:   resourceApiManagementAPIPolicyRead,
		Update: resourceApiManagementAPIPolicyCreateUpdate,
		Delete: resourceApiManagementAPIPolicyDelete,
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

			"api_name": azure.SchemaApiManagementApiName(),

			"xml_content": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"xml_link"},
				DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
			},

			"xml_link": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
			},
		},
	}
}

func resourceApiManagementAPIPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, apimanagement.PolicyExportFormatXML)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing API Policy (API Management Service %q / API %q / Resource Group %q): %s", serviceName, apiName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_policy", *existing.ID)
		}
	}

	parameters := apimanagement.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlLink != "" {
		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			Format: apimanagement.RawxmlLink,
			Value:  utils.String(xmlLink),
		}
	} else if xmlContent != "" {
		// this is intentionally an else-if since `xml_content` is computed

		// clear out any existing value for xml_link
		if !d.IsNewResource() {
			d.Set("xml_link", "")
		}

		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			Format: apimanagement.Rawxml,
			Value:  utils.String(xmlContent),
		}
	}

	if parameters.PolicyContractProperties == nil {
		return fmt.Errorf("Either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, apimanagement.PolicyExportFormatXML)
	if err != nil {
		return fmt.Errorf("retrieving API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementAPIPolicyRead(d, meta)
}

func resourceApiManagementAPIPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	apiName := id.Path["apis"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, apimanagement.PolicyExportFormatXML)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Policy (Resource Group %q / API Management Service %q / API %q) was not found - removing from state!", resourceGroup, serviceName, apiName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("api_name", apiName)

	if properties := resp.PolicyContractProperties; properties != nil {
		policyContent := ""
		if pc := properties.Value; pc != nil {
			policyContent = html.UnescapeString(*pc)
		}

		// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
		// as such there is no way to set `xml_link` and we'll let Terraform handle it
		d.Set("xml_content", policyContent)
	}

	return nil
}

func resourceApiManagementAPIPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
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
			return fmt.Errorf("deleting API Policy (Resource Group %q / API Management Service %q / API %q): %+v", resourceGroup, serviceName, apiName, err)
		}
	}

	return nil
}
