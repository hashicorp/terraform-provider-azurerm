package connections

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := connections.ParseConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"managed_api_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: managedapis.ValidateManagedApiID,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				// @tombuildsstuff: this can't be patched in API version 2016-06-01 and there isn't Swagger for
				// API version 2018-07-01-preview, so I guess this is ForceNew for now
				//
				// > Status=400 Code="PatchApiConnectionPropertiesNotSupported"
				// > Message="The request to patch API connection 'acctestconn-220307135205093274' is not supported.
				// > None of the fields inside the properties object can be patched."
				ForceNew: true,
			},

			"parameter_values": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				// @tombuildsstuff: this can't be patched in API version 2016-06-01 and there isn't Swagger for
				// API version 2018-07-01-preview, so I guess this is ForceNew for now
				//
				// > Status=400 Code="PatchApiConnectionPropertiesNotSupported"
				// > Message="The request to patch API connection 'acctestconn-220307135205093274' is not supported.
				// > None of the fields inside the properties object can be patched."
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := connections.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_api_connection", id.ID())
	}

	managedAppId, err := managedapis.ParseManagedApiID(d.Get("managed_api_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `managed_app_id`: %+v", err)
	}
	location := location.Normalize(managedAppId.Location)
	parameterValues := expandConnectionParameterValues(d.Get("parameter_values").(map[string]interface{}))
	model := connections.ApiConnectionDefinition{
		Location: utils.String(location),
		Properties: &connections.ApiConnectionDefinitionProperties{
			Api: &connections.ApiReference{
				Id: utils.String(managedAppId.ID()),
			},
			DisplayName:     utils.String(d.Get("display_name").(string)),
			ParameterValues: parameterValues,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if v := d.Get("display_name").(string); v != "" {
		model.Properties.DisplayName = utils.String(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceConnectionRead(d, meta)
}

func resourceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)

			apiId := ""
			if props.Api != nil && props.Api.Id != nil {
				apiId = *props.Api.Id
			}
			d.Set("managed_api_id", apiId)

			parameterValues := flattenConnectionParameterValues(props.ParameterValues)
			if err := d.Set("parameter_values", parameterValues); err != nil {
				return fmt.Errorf("setting `parameter_values`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	model := connections.ApiConnectionDefinition{
		// @tombuildsstuff: this can't be patched in API version 2016-06-01 and there isn't Swagger for
		// API version 2018-07-01-preview, so I guess this is ForceNew for now. The following error is returned
		// for both CreateOrUpdate and Update:
		//
		// > Status=400 Code="PatchApiConnectionPropertiesNotSupported"
		// > Message="The request to patch API connection 'acctestconn-220307135205093274' is not supported.
		// > None of the fields inside the properties object can be patched."
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.Update(ctx, *id, model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceConnectionRead(d, meta)
}

func resourceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandConnectionParameterValues(input map[string]interface{}) *map[string]string {
	parameterValues := make(map[string]string)
	for k, v := range input {
		parameterValues[k] = v.(string)
	}
	return &parameterValues
}

func flattenConnectionParameterValues(input *map[string]string) map[string]interface{} {
	parameterValues := make(map[string]interface{})
	if input != nil {
		for k, v := range *input {
			parameterValues[k] = v
		}
	}
	return parameterValues
}
