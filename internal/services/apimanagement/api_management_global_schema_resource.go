// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	globalSchema "github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementGlobalSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGlobalSchemaCreateUpdate,
		Read:   resourceApiManagementGlobalSchemaRead,
		Update: resourceApiManagementGlobalSchemaCreateUpdate,
		Delete: resourceApiManagementGlobalSchemaDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := globalSchema.ParseSchemaID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(globalSchema.SchemaTypeJson),
					string(globalSchema.SchemaTypeXml)}, false),
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

	id := globalSchema.NewSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("schema_id").(string))

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
	parameters := globalSchema.GlobalSchemaContract{
		Properties: &globalSchema.GlobalSchemaContractProperties{
			Description: utils.String(d.Get("description").(string)),
			SchemaType:  globalSchema.SchemaType(d.Get("type").(string)),
		},
	}

	// value for type=xml, document for type=json
	value := d.Get("value")
	if d.Get("type").(string) == string(globalSchema.SchemaTypeXml) {
		parameters.Properties.Value = &value
	} else {
		var document interface{}
		if err := json.Unmarshal([]byte(value.(string)), &document); err != nil {
			return fmt.Errorf(" error preparing value data to send %s: %s", id, err)
		}
		parameters.Properties.Document = &document
	}

	future, err := client.GlobalSchemaCreateOrUpdate(ctx, id, parameters, globalSchema.DefaultGlobalSchemaCreateOrUpdateOperationOptions())
	if err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}
	if err := future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiManagementGlobalSchemaRead(d, meta)
}

func resourceApiManagementGlobalSchemaRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := globalSchema.ParseSchemaID(d.Id())
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
	d.Set("type", resp.Model.Properties.SchemaType)
	d.Set("description", resp.Model.Properties.Description)

	if resp.Model != nil {
		if resp.Model.Properties.SchemaType == globalSchema.SchemaTypeXml {
			d.Set("value", resp.Model.Properties.Value)
		} else {
			var document []byte
			if document, err = json.Marshal(resp.Model.Properties.Document); err != nil {
				return fmt.Errorf(" reading the schema document %s: %s", *id, err)
			}
			d.Set("value", string(document))
		}
	}

	return nil
}

func resourceApiManagementGlobalSchemaDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := globalSchema.ParseSchemaID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.GlobalSchemaDelete(ctx, *id, globalSchema.DefaultGlobalSchemaDeleteOperationOptions()); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %s", *id, err)
		}
	}

	return nil
}
