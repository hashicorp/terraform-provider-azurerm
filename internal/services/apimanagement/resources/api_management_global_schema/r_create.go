// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_global_schema

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGlobalSchemaCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GlobalSchemaClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := schema.NewSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("schema_id").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		var document any
		if err := json.Unmarshal([]byte(value.(string)), &document); err != nil {
			return fmt.Errorf(" error preparing value data to send %s: %s", id, err)
		}
		payload.Properties.Document = &document
	}
	if d.Get("type").(string) == string(schema.SchemaTypeXml) {
		payload.Properties.Value = &value
	}

	if d.IsNewResource() {
		if err := client.GlobalSchemaCreateOrUpdateCallbackThenPoll(ctx, id, payload, schema.DefaultGlobalSchemaCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %s", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.GlobalSchemaCreateOrUpdateThenPoll(ctx, id, payload, schema.DefaultGlobalSchemaCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %s", id, err)
		}
	}

	return resourceApiManagementGlobalSchemaRead(d, meta)
}
