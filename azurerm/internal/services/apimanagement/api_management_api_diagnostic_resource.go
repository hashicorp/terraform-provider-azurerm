package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementApiDiagnostic() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementApiDiagnosticCreateUpdate,
		Read:   resourceArmApiManagementApiDiagnosticRead,
		Update: resourceArmApiManagementApiDiagnosticCreateUpdate,
		Delete: resourceArmApiManagementApiDiagnosticDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ApiManagementApiDiagnosticID(id)
			return err
		}),

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
					"azuremonitor",
				}, false),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"api_name": azure.SchemaApiManagementApiName(),

			"api_management_logger_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementLoggerID,
			},
		},
	}
}

func resourceArmApiManagementApiDiagnosticCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId := d.Get("identifier").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, diagnosticId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %s", diagnosticId, resourceGroup, serviceName, apiName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_diagnostic", *existing.ID)
		}
	}

	parameters := apimanagement.DiagnosticContract{
		DiagnosticContractProperties: &apimanagement.DiagnosticContractProperties{
			LoggerID: utils.String(d.Get("api_management_logger_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, diagnosticId, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId, resourceGroup, serviceName, apiName, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, diagnosticId)
	if err != nil {
		return fmt.Errorf("retrieving Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId, resourceGroup, serviceName, apiName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("reading ID for Diagnostic %q (Resource Group %q / API Management Service %q / API %q): ID is empty", diagnosticId, resourceGroup, serviceName, apiName)
	}
	d.SetId(*resp.ID)

	return resourceArmApiManagementApiDiagnosticRead(d, meta)
}

func resourceArmApiManagementApiDiagnosticRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := parse.ApiManagementApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Diagnostic %q (Resource Group %q / API Management Service %q / API %q) was not found - removing from state!", diagnosticId.Name, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId.Name, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, err)
	}

	d.Set("api_name", diagnosticId.ApiName)
	d.Set("identifier", resp.Name)
	d.Set("resource_group_name", diagnosticId.ResourceGroup)
	d.Set("api_management_name", diagnosticId.ServiceName)
	if props := resp.DiagnosticContractProperties; props != nil {
		d.Set("api_management_logger_id", props.LoggerID)
	}

	return nil
}

func resourceArmApiManagementApiDiagnosticDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiDiagnosticClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	diagnosticId, err := parse.ApiManagementApiDiagnosticID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, diagnosticId.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Diagnostic %q (Resource Group %q / API Management Service %q / API %q): %+v", diagnosticId.Name, diagnosticId.ResourceGroup, diagnosticId.ServiceName, diagnosticId.ApiName, err)
		}
	}

	return nil
}
