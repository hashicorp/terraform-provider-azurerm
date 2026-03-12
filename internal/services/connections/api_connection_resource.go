// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package connections

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
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
)

func resourceApiConnection() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceApiConnectionCreate,
		Read:   resourceApiConnectionRead,
		Update: resourceApiConnectionUpdate,
		Delete: resourceApiConnectionDelete,

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
				// Note: O+C because Azure sets a default when `display_name` is not defined but the value depends on which managed API is provided.
				// For example:
				//   - Managed API `servicebus` defaults to `Service Bus`
				//   - Managed API `sftpwithssh` defaults to `SFTP - SSH`
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"parameter_values": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}

	return resource
}

func resourceApiConnectionCreate(d *schema.ResourceData, meta interface{}) error {
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
	location := location.Normalize(managedAppId.LocationName)
	model := connections.ApiConnectionDefinition{
		Location: pointer.To(location),
		Properties: &connections.ApiConnectionDefinitionProperties{
			Api: &connections.ApiReference{
				Id: pointer.To(managedAppId.ID()),
			},
			DisplayName:     pointer.To(d.Get("display_name").(string)),
			ParameterValues: pointer.To(d.Get("parameter_values").(map[string]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if v := d.Get("display_name").(string); v != "" {
		model.Properties.DisplayName = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiConnectionRead(d, meta)
}

func resourceApiConnectionRead(d *schema.ResourceData, meta interface{}) error {
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

			// In version 2016-06-01 the API doesn't return `ParameterValues`.
			// The non-secret parameters are returned in `NonSecretParameterValues` instead.
			if err := d.Set("parameter_values", flattenParameterValues(pointer.From(props.NonSecretParameterValues))); err != nil {
				return fmt.Errorf("setting `parameter_values`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceApiConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("display_name") {
		props.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	// The GET operation returns `NonSecretParameterValues` but we're making updates through `ParameterValues`
	// so we remove `NonSecretParameterValues` from the request to avoid conflicting parameters.
	// this is fixed in later (preview) versions of the API but these don't have an API spec available.
	props.NonSecretParameterValues = nil
	if d.HasChange("parameter_values") {
		props.ParameterValues = pointer.To(d.Get("parameter_values").(map[string]interface{}))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceApiConnectionRead(d, meta)
}

func resourceApiConnectionDelete(d *schema.ResourceData, meta interface{}) error {
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

// Because this API may return other primitive types for `parameter_values`
// we need to ensure each value in the map is a string to prevent panics when setting this into state.
func flattenParameterValues(input map[string]interface{}) map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		output[k] = fmt.Sprintf("%v", v)
	}

	return output
}
