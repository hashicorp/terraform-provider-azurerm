// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_service_endpoint_storage_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetServiceEndpointStoragePolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("definition") {
		payload.Properties = &serviceendpointpolicies.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]any)),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}
