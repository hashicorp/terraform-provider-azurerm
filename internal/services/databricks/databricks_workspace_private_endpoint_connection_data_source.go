package databricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/sdk/2021-04-01-preview/workspaces"
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
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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
		result := make(map[string]interface{})

		if name := v.Name; name != nil {
			result["name"] = *name
		}

		if id := v.Id; id != nil {
			result["workspace_private_endpoint_id"] = *id
		}

		connState := v.Properties.PrivateLinkServiceConnectionState
		if description := connState.Description; description != nil {
			result["description"] = *description
		}
		if status := connState.Status; status != "" {
			result["status"] = status
		}
		if actionReq := connState.ActionRequired; actionReq != nil {
			result["action_required"] = *actionReq
		}

		results = append(results, result)
	}

	return results
}
