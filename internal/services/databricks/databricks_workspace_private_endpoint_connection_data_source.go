package databricks

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/databricks/mgmt/2021-04-01-preview/databricks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

	id, err := parse.WorkspaceID(workspace)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("workspace_id", workspace)
	d.Set("private_endpoint_id", endpointId)

	if props := resp.WorkspaceProperties; props != nil {
		if err := d.Set("connections", flattenPrivateEndpointConnections(props.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("setting `connections`: %+v", err)
		}
	}

	return nil
}

func flattenPrivateEndpointConnections(input *[]databricks.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		result := make(map[string]interface{})

		if name := v.Name; name != nil {
			result["name"] = *name
		}

		if id := v.ID; id != nil {
			result["workspace_private_endpoint_id"] = *id
		}

		if props := v.Properties; props != nil {
			if connState := props.PrivateLinkServiceConnectionState; connState != nil {
				if description := connState.Description; description != nil {
					result["description"] = *description
				}
				if status := connState.Status; status != "" {
					result["status"] = status
				}
				if actionReq := connState.ActionRequired; actionReq != nil {
					result["action_required"] = *actionReq
				}
			}
		}

		results = append(results, result)
	}

	return results
}
