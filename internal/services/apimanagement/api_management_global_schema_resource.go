// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGlobalSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGlobalSchemaCreateUpdate,
		Read:   resourceApiManagementGlobalSchemaRead,
		Update: resourceApiManagementGlobalSchemaCreateUpdate,
		Delete: resourceApiManagementGlobalSchemaDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := schema.ParseSchemaID(id)
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

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"type": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(schema.PossibleValuesForSchemaType(), false),
			},

			"value": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsNotEmpty,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceApiManagementGlobalSchemaCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schema.NewSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("schema_id").(string))
	if d.IsNewResource() {
		existing, err := client.GlobalSchemaGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_global_schema", id.ID())
		}
	}

	payload := schema.GlobalSchemaContract{
		Properties: &schema.GlobalSchemaContractProperties{
			Description: pointer.To(d.Get("description").(string)),
			SchemaType:  schema.SchemaType(d.Get("type").(string)),
		},
	}

	// value for type=xml, document for type=json
	value := d.Get("value")
	if d.Get("type").(string) == string(schema.SchemaTypeJson) {
		var document interface{}
		if err := json.Unmarshal([]byte(value.(string)), &document); err != nil {
			return fmt.Errorf(" error preparing value data to send %s: %s", id, err)
		}
		payload.Properties.Document = &document
	}
	if d.Get("type").(string) == string(schema.SchemaTypeXml) {
		payload.Properties.Value = &value
	}

	if err := client.GlobalSchemaCreateOrUpdateThenPoll(ctx, id, payload, schema.DefaultGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementGlobalSchemaRead(d, meta)
}

func resourceApiManagementGlobalSchemaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schema.ParseSchemaID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GlobalSchemaGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %s", *id, err)
	}

	d.Set("schema_id", id.SchemaId)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("type", props.SchemaType)

			var value interface{}
			// value for type=xml, document for type=json
			if props.SchemaType == schema.SchemaTypeJson && props.Document != nil {
				var document []byte
				if document, err = json.Marshal(props.Document); err != nil {
					return fmt.Errorf(" reading the schema document %s: %s", *id, err)
				}
				value = string(document)
			}
			if props.SchemaType == schema.SchemaTypeXml && props.Value != nil {
				value = *props.Value
			}
			d.Set("value", value)
		}
	}

	return nil
}

func resourceApiManagementGlobalSchemaDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := schema.ParseSchemaID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.GlobalSchemaDelete(ctx, *id, schema.DefaultGlobalSchemaDeleteOperationOptions()); err != nil {
		return fmt.Errorf("deleting %s: %s", *id, err)
	}

	return nil
}
