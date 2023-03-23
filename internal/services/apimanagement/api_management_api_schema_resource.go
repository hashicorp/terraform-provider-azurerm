package apimanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
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

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"content_type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"value": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					if d.Get("content_type") == "application/vnd.ms-azure-apim.swagger.definitions+json" || d.Get("content_type") == "application/vnd.oai.openapi.components+json" {
						return pluginsdk.SuppressJsonDiff(k, old, new, d)
					}
					return old == new
				},
				ExactlyOneOf: []string{"value", "definitions", "components"},
			},

			"components": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ExactlyOneOf:     []string{"value", "definitions", "components"},
			},

			"definitions": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
				ExactlyOneOf:     []string{"value", "definitions", "components"},
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
	parameters := apimanagement.SchemaContract{
		SchemaContractProperties: &apimanagement.SchemaContractProperties{
			ContentType:              &contentType,
			SchemaDocumentProperties: &apimanagement.SchemaDocumentProperties{},
		},
	}

	if v, ok := d.GetOk("value"); ok {
		parameters.SchemaContractProperties.SchemaDocumentProperties.Value = utils.String(v.(string))
	}

	if v, ok := d.GetOk("components"); ok {
		parameters.SchemaContractProperties.SchemaDocumentProperties.Components = v.(string)
	}

	if v, ok := d.GetOk("definitions"); ok {
		parameters.SchemaContractProperties.SchemaDocumentProperties.Definitions = v.(string)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.SchemaName, parameters, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
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
			if documentProperties.Value != nil {
				d.Set("value", documentProperties.Value)
			}

			if properties.Components != nil {
				value, err := convert2Str(properties.Components)
				if err != nil {
					return err
				}
				d.Set("components", value)
			}

			if properties.Definitions != nil {
				value, err := convert2Str(properties.Definitions)
				if err != nil {
					return err
				}
				d.Set("definitions", value)
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

func convert2Str(rawVal interface{}) (string, error) {
	value := ""
	if val, ok := rawVal.(string); ok {
		value = val
	} else {
		val, err := json.Marshal(rawVal)
		if err != nil {
			return "", fmt.Errorf("failed to marshal to json: %+v", err)
		}
		value = string(val)
	}
	return value, nil
}
