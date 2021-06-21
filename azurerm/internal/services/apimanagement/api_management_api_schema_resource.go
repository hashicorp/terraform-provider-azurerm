package apimanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementApiSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiSchemaCreateUpdate,
		Read:   resourceApiManagementApiSchemaRead,
		Update: resourceApiManagementApiSchemaCreateUpdate,
		Delete: resourceApiManagementApiSchemaDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"schema_id": schemaz.SchemaApiManagementChildName(),

			"api_name": schemaz.SchemaApiManagementApiName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"content_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					if d.Get("content_type") == "application/vnd.ms-azure-apim.swagger.definitions+json" || d.Get("content_type") == "application/vnd.oai.openapi.components+json" {
						return pluginsdk.SuppressJsonDiff(k, old, new, d)
					}
					return old == new
				},
			},
		},
	}
}

func resourceApiManagementApiSchemaCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiSchemasClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	schemaID := d.Get("schema_id").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	apiName := d.Get("api_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_api_schema", *existing.ID)
		}
	}

	contentType := d.Get("content_type").(string)
	value := d.Get("value").(string)
	parameters := apimanagement.SchemaContract{
		SchemaContractProperties: &apimanagement.SchemaContractProperties{
			ContentType: &contentType,
			SchemaDocumentProperties: &apimanagement.SchemaDocumentProperties{
				Value: &value,
			},
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, apiName, schemaID, parameters, ""); err != nil {
		return fmt.Errorf("creating or updating API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}

	//lintignore:R006
	err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return pluginsdk.RetryableError(fmt.Errorf("Expected schema %q (API Management Service %q / API %q / Resource Group %q) to be created but was in non existent state, retrying", schemaID, serviceName, apiName, resourceGroup))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("Error getting schema %q (API Management Service %q / API %q / Resource Group %q): %+v", schemaID, serviceName, apiName, resourceGroup, err))
		}
		if resp.ID == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("Cannot read ID for API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err))
		}
		d.SetId(*resp.ID)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error getting schema %q (API Management Service %q / API %q / Resource Group %q): %+v", schemaID, serviceName, apiName, resourceGroup, err)
	}
	return resourceApiManagementApiSchemaRead(d, meta)
}

func resourceApiManagementApiSchemaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiSchemasClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiSchemaID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiName := id.ApiName
	schemaID := id.SchemaName

	resp, err := client.Get(ctx, resourceGroup, serviceName, apiName, schemaID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] API Schema %q (API Management Service %q / API %q / Resource Group %q) was not found - removing from state!", schemaID, serviceName, apiName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
	}

	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)
	d.Set("api_name", apiName)
	d.Set("schema_id", schemaID)

	if properties := resp.SchemaContractProperties; properties != nil {
		d.Set("content_type", properties.ContentType)
		if documentProperties := properties.SchemaDocumentProperties; documentProperties != nil {
			/*
				As per https://docs.microsoft.com/en-us/rest/api/apimanagement/2019-12-01/api-schema/get#schemacontract

				- Swagger Schema use application/vnd.ms-azure-apim.swagger.definitions+json
				- WSDL Schema use application/vnd.ms-azure-apim.xsd+xml
				- OpenApi Schema use application/vnd.oai.openapi.components+json
				- WADL Schema use application/vnd.ms-azure-apim.wadl.grammars+xml.

				Definitions used for Swagger/OpenAPI schemas only, otherwise Value is used
			*/
			switch *properties.ContentType {
			case "application/vnd.ms-azure-apim.swagger.definitions+json", "application/vnd.oai.openapi.components+json":
				if documentProperties.Definitions != nil {
					value, err := json.Marshal(documentProperties.Definitions)
					if err != nil {
						return fmt.Errorf("[FATAL] Unable to serialize schema to json. Error: %+v. Schema struct: %+v", err, documentProperties.Definitions)
					}
					d.Set("value", string(value))
				}
			case "application/vnd.ms-azure-apim.xsd+xml", "application/vnd.ms-azure-apim.wadl.grammars+xml":
				d.Set("value", documentProperties.Value)
			default:
				log.Printf("[WARN] Unknown content type %q for schema %q (API Management Service %q / API %q / Resource Group %q)", *properties.ContentType, schemaID, serviceName, apiName, resourceGroup)
				d.Set("value", documentProperties.Value)
			}
		}
	}
	return nil
}

func resourceApiManagementApiSchemaDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiSchemasClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApiSchemaID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	apiName := id.ApiName
	schemaID := id.SchemaName

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, apiName, schemaID, "", utils.Bool(false)); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting API Schema %q (API Management Service %q / API %q / Resource Group %q): %s", schemaID, serviceName, apiName, resourceGroup, err)
		}
	}

	return nil
}
