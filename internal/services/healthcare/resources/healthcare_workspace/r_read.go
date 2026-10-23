// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisWorkspaceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		if props := m.Properties; props != nil {
			d.Set("private_endpoint_connection", flattenWorkspacePrivateEndpoint(props.PrivateEndpointConnections))
		}

		if err := tags.FlattenAndSet(d, m.Tags); err != nil {
			return err
		}
	}

	return nil
}

func flattenWorkspacePrivateEndpoint(input *[]workspaces.PrivateEndpointConnection) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]any{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}

		if endpoint.Id != nil {
			result["id"] = *endpoint.Id
		}
	}
	return results
}
