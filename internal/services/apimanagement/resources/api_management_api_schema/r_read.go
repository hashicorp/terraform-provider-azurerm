// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_api_schema

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apischema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementApiSchemaRead(d *pluginsdk.ResourceData, meta any) error {
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

func convert2Str(rawVal any) (string, error) {
	var value string
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
