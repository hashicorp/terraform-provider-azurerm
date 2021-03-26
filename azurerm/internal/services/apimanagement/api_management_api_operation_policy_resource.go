package apimanagement

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiOperationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementAPIOperationPolicyCreateUpdate,
		Read:   resourceApiManagementAPIOperationPolicyRead,
		Update: resourceApiManagementAPIOperationPolicyCreateUpdate,
		Delete: resourceApiManagementAPIOperationPolicyDelete,
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

			"api_management_name": schemaz.SchemaApiManagementName(),

			"api_name": schemaz.SchemaApiManagementApiName(),

			"operation_id": schemaz.SchemaApiManagementChildName(),

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

func resourceApiManagementAPIOperationPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)
	operationID := d.Get("operation_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, operationID, apimanagement.PolicyExportFormatXML)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing API Operation Policy (API Management Service %q / API %q / Operation %q / Resource Group %q): %s", serviceName, apiName, operationID, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_operation_policy", *existing.ID)
		}
	}

	parameters := apimanagement.PolicyContract{}

	xmlContent := d.Get("xml_content").(string)
	xmlLink := d.Get("xml_link").(string)

	if xmlContent != "" {
		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			Format: apimanagement.Rawxml,
			Value:  utils.String(xmlContent),
		}
	}

	if xmlLink != "" {
		parameters.PolicyContractProperties = &apimanagement.PolicyContractProperties{
			Format: apimanagement.RawxmlLink,
			Value:  utils.String(xmlLink),
		}
	}

	if parameters.PolicyContractProperties == nil {
		return fmt.Errorf("Either `xml_content` or `xml_link` must be set")
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, operationID, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q): %+v", resourceGroup, serviceName, apiName, operationID, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, operationID, apimanagement.PolicyExportFormatXML)
	if err != nil {
		return fmt.Errorf("retrieving API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q): %+v", resourceGroup, serviceName, apiName, operationID, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q): %+v", resourceGroup, serviceName, apiName, operationID, err)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementAPIOperationPolicyRead(d, meta)
}

func resourceApiManagementAPIOperationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiOperationPolicyID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiName := id.ApiName
	operationName := id.OperationName

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, operationName, apimanagement.PolicyExportFormatXML)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q) was not found - removing from state!", resourceGroup, serviceName, apiName, operationName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q): %+v", resourceGroup, serviceName, apiName, operationName, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("api_name", apiName)
	d.Set("operation_id", operationName)

	if properties := resp.PolicyContractProperties; properties != nil {
		// when you submit an `xml_link` to the API, the API downloads this link and stores it as `xml_content`
		// as such there is no way to set `xml_link` and we'll let Terraform handle it
		d.Set("xml_content", html.UnescapeString(*properties.Value))
	}

	return nil
}

func resourceApiManagementAPIOperationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiOperationPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiOperationPolicyID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiName := id.ApiName
	operationName := id.OperationName

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apiName, operationName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting API Operation Policy (Resource Group %q / API Management Service %q / API %q / Operation %q): %+v", resourceGroup, serviceName, apiName, operationName, err)
		}
	}

	return nil
}
