// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiConnectionRead(d *schema.ResourceData, meta any) error {
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

// Because this API may return other primitive types for `parameter_values`
// we need to ensure each value in the map is a string to prevent panics when setting this into state.
func flattenParameterValues(input map[string]any) map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		output[k] = fmt.Sprintf("%v", v)
	}

	return output
}
