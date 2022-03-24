package apimanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementApiSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiSchemaCreateUpdate,
		Read:   resourceApiManagementApiSchemaRead,
		Update: resourceApiManagementApiSchemaCreateUpdate,
		Delete: resourceApiManagementApiSchemaDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ApiSchemaID(id)
			return err
		}),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewApiSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("schema_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_api_schema", id.ID())
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

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName, parameters, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return pluginsdk.RetryableError(fmt.Errorf("expected schema %s to be created but was in non existent state, retrying", id))
			}
			return pluginsdk.NonRetryableError(fmt.Errorf("getting schema %s: %+v", id, err))
		}
		if resp.ID == nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("cannot read ID for %s: %s", id, err))
		}
		d.SetId(id.ID())
		return nil
	})
	if err != nil {
		return fmt.Errorf("getting %s: %+v", id, err)
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

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %s", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)
	d.Set("api_name", id.ApiName)
	d.Set("schema_id", id.SchemaName)

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
				log.Printf("[WARN] Unknown content type %q for %s", *properties.ContentType, *id)
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

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName, "", utils.Bool(false)); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}
