// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDatabricksWorkspacePrivateEndpointConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDatabricksWorkspacePrivateEndpointConnectionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"private_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: privateendpoints.ValidatePrivateEndpointID,
			},

			"connections": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"workspace_private_endpoint_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"status": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"action_required": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDatabricksWorkspacePrivateEndpointConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	workspace := d.Get("workspace_id").(string)
	endpointId := d.Get("private_endpoint_id").(string)

	id, err := workspaces.ParseWorkspaceID(workspace)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	d.Set("workspace_id", workspace)
	d.Set("private_endpoint_id", endpointId)

	if model := resp.Model; model != nil {
		if err := d.Set("connections", flattenPrivateEndpointConnections(model.Properties.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("setting `connections`: %+v", err)
		}
	}

	return nil
}

func flattenPrivateEndpointConnections(input *[]workspaces.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		workspacePrivateEndpointId := ""
		if v.Id != nil {
			workspacePrivateEndpointId = *v.Id
		}

		connState := v.Properties.PrivateLinkServiceConnectionState
		actionRequired := ""
		if connState.ActionsRequired != nil {
			actionRequired = *connState.ActionsRequired
		}
		description := ""
		if connState.Description != nil {
			description = *connState.Description
		}
		status := ""
		if connState.Status != "" {
			status = string(connState.Status)
		}

		results = append(results, map[string]interface{}{
			"action_required":               actionRequired,
			"description":                   description,
			"name":                          name,
			"status":                        status,
			"workspace_private_endpoint_id": workspacePrivateEndpointId,
		})
	}

	return results
}
