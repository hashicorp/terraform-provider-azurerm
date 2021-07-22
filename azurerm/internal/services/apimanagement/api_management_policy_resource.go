package apimanagement

import (
	"fmt"
	"html"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementPolicyCreateUpdate,
		Read:   resourceApiManagementPolicyRead,
		Update: resourceApiManagementPolicyCreateUpdate,
		Delete: resourceApiManagementPolicyDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"xml_content": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"xml_link"},
				ExactlyOneOf:     []string{"xml_link", "xml_content"},
				DiffSuppressFunc: XmlWithDotNetInterpolationsDiffSuppress,
			},

			"xml_link": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xml_content"},
				ExactlyOneOf:  []string{"xml_link", "xml_content"},
			},
		},
	}
}

func resourceApiManagementPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apiManagementID := d.Get("api_management_id").(string)
	id, err := parse.ApiManagementID(apiManagementID)
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	/*
		Other resources would have a check for d.IsNewResource() at this location, and would error out using `tf.ImportAsExistsError` if the resource already existed.
		However, this is a sub-resource, and the API always returns a policy when queried, either a default policy or one configured by the user or by this pluginsdk.
		Instead of the usual check, the resource documentation clearly states that any existing policy will be overwritten if the resource is used.
	*/

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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Policy (Resource Group %q / API Management Service %q): %+v", resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.PolicyExportFormatXML)
	if err != nil {
		return fmt.Errorf("retrieving Policy (Resource Group %q / API Management Service %q): %+v", resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Policy (Resource Group %q / API Management Service %q): %+v", resourceGroup, serviceName, err)
	}
	d.SetId(*resp.ID)

	return resourceApiManagementPolicyRead(d, meta)
}

func resourceApiManagementPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	serviceClient := meta.(*clients.Client).ApiManagement.ServiceClient
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	serviceResp, err := serviceClient.Get(ctx, resourceGroup, serviceName)
	if err != nil {
		if utils.ResponseWasNotFound(serviceResp.Response) {
			log.Printf("API Management Service %q was not found in Resource Group %q - removing Policy from state!", serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on API Management Service %q (Resource Group %q): %+v", serviceName, resourceGroup, err)
	}

	d.Set("api_management_id", serviceResp.ID)

	resp, err := client.Get(ctx, resourceGroup, serviceName, apimanagement.PolicyExportFormatXML)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Policy (Resource Group %q / API Management Service %q) was not found - removing from state!", resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Policy (Resource Group %q / API Management Service %q): %+v", resourceGroup, serviceName, err)
	}

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

func resourceApiManagementPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.PolicyClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Policy (Resource Group %q / API Management Service %q): %+v", resourceGroup, serviceName, err)
		}
	}

	return nil
}
