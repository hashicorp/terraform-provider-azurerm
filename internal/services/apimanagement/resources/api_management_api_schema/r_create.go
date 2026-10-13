// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_schema

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiSchemaCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ApiSchemasClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := apischema.NewApiSchemaID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("api_name").(string), d.Get("schema_id").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		var value any
		if err := json.Unmarshal([]byte(v.(string)), &value); err != nil {
			return fmt.Errorf("failed to unmarshal components %v: %+v", v.(string), err)
		}

		parameters.Properties.Document.Components = pointer.To(value)
	}

	if v, ok := d.GetOk("definitions"); ok {
		var value any
		if err := json.Unmarshal([]byte(v.(string)), &value); err != nil {
			return fmt.Errorf("failed to unmarshal definitions %v: %+v", v.(string), err)
		}

		parameters.Properties.Document.Definitions = pointer.To(value)
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, apischema.CreateOrUpdateOperationOptions{}, sdk.SetIDCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %s", id, err)
		}
		d.SetId(id.ID())
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, apischema.CreateOrUpdateOperationOptions{}); err != nil {
			return fmt.Errorf("updating %s: %s", id, err)
		}
	}

	return resourceApiManagementApiSchemaRead(d, meta)
}
