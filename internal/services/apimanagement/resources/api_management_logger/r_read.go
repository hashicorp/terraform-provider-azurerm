// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_logger

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementLoggerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.LoggerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := logger.ParseLoggerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("resource_id", pointer.From(props.ResourceId))
			d.Set("buffered", pointer.From(props.IsBuffered))
			d.Set("description", pointer.From(props.Description))
			if err := d.Set("eventhub", flattenApiManagementLoggerEventHub(d, props)); err != nil {
				return fmt.Errorf("setting `eventhub`: %s", err)
			}
		}
	}

	return nil
}

func flattenApiManagementLoggerEventHub(d *pluginsdk.ResourceData, properties *logger.LoggerContractProperties) []any {
	result := make([]any, 0)
	if c := properties.Credentials; c != nil && (*c)["name"] != "" {
		eventHub := make(map[string]any)
		eventHub["name"] = (*c)["name"]
		if existing := d.Get("eventhub").([]any); len(existing) > 0 {
			existingEventHub := existing[0].(map[string]any)
			if conn, ok := existingEventHub["connection_string"]; ok {
				eventHub["connection_string"] = conn.(string)
			}
			if endpoint, ok := existingEventHub["endpoint_uri"]; ok {
				eventHub["endpoint_uri"] = endpoint
			}
			if clientId, ok := existingEventHub["user_assigned_identity_client_id"]; ok {
				if clientId != "SystemAssigned" {
					eventHub["user_assigned_identity_client_id"] = clientId
				}
			}
		}
		result = append(result, eventHub)
	}
	return result
}
