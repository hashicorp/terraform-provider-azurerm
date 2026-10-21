// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_service_endpoint_storage_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetServiceEndpointStoragePolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ServiceEndpointPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			if err := d.Set("definition", flattenServiceEndpointPolicyDefinitions(props.ServiceEndpointPolicyDefinitions)); err != nil {
				return fmt.Errorf("setting `definition`: %v", err)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenServiceEndpointPolicyDefinitions(input *[]serviceendpointpolicies.ServiceEndpointPolicyDefinition) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0)
	for _, e := range *input {
		var (
			description     = ""
			service         = ""
			serviceResource = []any{}
		)
		if b := e.Properties; b != nil {
			if b.Description != nil {
				description = *b.Description
			}
			serviceResource = pluginsdk.FlattenSlice(b.ServiceResources)
			if b.Service != nil {
				service = *b.Service
			}
		}

		output = append(output, map[string]any{
			"name":              pointer.From(e.Name),
			"description":       description,
			"service_resources": serviceResource,
			"service":           service,
		})
	}

	return output
}
