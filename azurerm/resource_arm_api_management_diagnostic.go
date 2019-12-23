package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementDiagnostic() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementDiagnosticCreateUpdate,
		Read:   resourceArmApiManagementDiagnosticRead,
		Update: resourceArmApiManagementDiagnosticCreateUpdate,
		Delete: resourceArmApiManagementDiagnosticDelete,
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
			"identifier": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"applicationinsights",
				}, false),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceArmApiManagementDiagnosticCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId := d.Get("identifier").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	enabled := d.Get("enabled").(bool)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, diagnosticId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Diagnostic %q (API Management Service %q / Resource Group %q): %s", diagnosticId, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_diagnostic", *existing.ID)
		}
	}

	parameters := apimanagement.DiagnosticContract{
		DiagnosticContractProperties: &apimanagement.DiagnosticContractProperties{
			Enabled: utils.Bool(enabled),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, diagnosticId, parameters, ""); err != nil {
		return fmt.Errorf("Error creating or updating Diagnostic %q (Resource Group %q / API Management Service %q): %+v", diagnosticId, resourceGroup, serviceName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, diagnosticId)
	if err != nil {
		return fmt.Errorf("Error retrieving Diagnostic %q (Resource Group %q / API Management Service %q): %+v", diagnosticId, resourceGroup, serviceName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Diagnostic %q (Resource Group %q / API Management Service %q)", diagnosticId, resourceGroup, serviceName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementDiagnosticRead(d, meta)
}

func resourceArmApiManagementDiagnosticRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	diagnosticId := id.Path["diagnostics"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, diagnosticId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Diagnostic %q (Resource Group %q / API Management Service %q) was not found - removing from state!", diagnosticId, resourceGroup, serviceName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request for Diagnostic %q (Resource Group %q / API Management Service %q): %+v", diagnosticId, resourceGroup, serviceName, err)
	}

	d.Set("identifier", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.DiagnosticContractProperties; props != nil {
		d.Set("enabled", props.Enabled)
	}

	return nil
}

func resourceArmApiManagementDiagnosticDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.DiagnosticClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	diagnosticId := id.Path["diagnostics"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, diagnosticId, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Diagnostic %q (Resource Group %q / API Management Service %q): %+v", diagnosticId, resourceGroup, serviceName, err)
		}
	}

	return nil
}
