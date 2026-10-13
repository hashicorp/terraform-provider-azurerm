// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_global_schema

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGlobalSchemaRead(d *pluginsdk.ResourceData, meta any) error {
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

			var value any
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
