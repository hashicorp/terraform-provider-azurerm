// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_backend

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/backend"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementBackendRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.BackendClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := backend.ParseBackendID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("protocol", pointer.FromEnum(props.Protocol))
			d.Set("resource_id", pointer.From(props.ResourceId))
			d.Set("title", pointer.From(props.Title))
			d.Set("url", props.Url)
			if err := d.Set("circuit_breaker_rule", flattenApiManagementBackendCircuitBreaker(props.CircuitBreaker)); err != nil {
				return fmt.Errorf("setting `circuit_breaker_rule`: %s", err)
			}
			if err := d.Set("credentials", flattenApiManagementBackendCredentials(props.Credentials)); err != nil {
				return fmt.Errorf("setting `credentials`: %s", err)
			}
			if err := d.Set("proxy", flattenApiManagementBackendProxy(props.Proxy)); err != nil {
				return fmt.Errorf("setting `proxy`: %s", err)
			}
			if properties := props.Properties; properties != nil {
				if err := d.Set("service_fabric_cluster", flattenApiManagementBackendServiceFabricCluster(properties.ServiceFabricCluster)); err != nil {
					return fmt.Errorf("setting `service_fabric_cluster`: %s", err)
				}
			}
			if err := d.Set("tls", flattenApiManagementBackendTls(props.Tls)); err != nil {
				return fmt.Errorf("setting `tls`: %s", err)
			}
		}
	}

	return nil
}
