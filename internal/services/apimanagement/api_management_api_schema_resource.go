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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementApiSchemaCreateUpdate,
		Read:   resourceApiManagementApiSchemaRead,
		Update: resourceApiManagementApiSchemaCreateUpdate,
		Delete: resourceApiManagementApiSchemaDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := apischema.ParseApiSchemaID(id)
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

	id := apischema.NewApiSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("schema_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_api_schema", id.ID())
		}
	}

	parameters := apischema.SchemaContract{
		Properties: &apischema.SchemaContractProperties{
			ContentType: d.Get("content_type").(string),
			Document:    apischema.SchemaDocumentProperties{},
		},
	}

	if v, ok := d.GetOk("value"); ok {
		parameters.Properties.Document.Value = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("components"); ok {
		var value interface{}
		if err := json.Unmarshal([]byte(v.(string)), &value); err != nil {
			return fmt.Errorf("failed to unmarshal components %v: %+v", v.(string), err)
		}

		parameters.Properties.Document.Components = pointer.To(value)
	}

	if v, ok := d.GetOk("definitions"); ok {
		var value interface{}
		if err := json.Unmarshal([]byte(v.(string)), &value); err != nil {
			return fmt.Errorf("failed to unmarshal definitions %v: %+v", v.(string), err)
		}

		parameters.Properties.Document.Definitions = pointer.To(value)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, apischema.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementApiSchemaRead(d, meta)
}

func resourceApiManagementApiSchemaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ApiSchemasClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := apischema.ParseApiSchemaID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %s", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)
	d.Set("api_name", id.ApiId)
	d.Set("schema_id", id.SchemaId)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("content_type", props.ContentType)
			documentProperties := props.Document
			if documentProperties.Value != nil {
				d.Set("value", pointer.From(documentProperties.Value))
			}

			if documentProperties.Components != nil {
				value, err := convert2Str(pointer.From(documentProperties.Components))
				if err != nil {
					return err
				}
				d.Set("components", value)
			}

			if documentProperties.Definitions != nil {
				value, err := convert2Str(documentProperties.Definitions)
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

	id, err := apischema.ParseApiSchemaID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, apischema.DeleteOperationOptions{Force: pointer.To(false)}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
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
