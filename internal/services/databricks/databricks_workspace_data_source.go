// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package databricks

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2024-05-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDatabricksWorkspace() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDatabricksWorkspaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"workspace_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"managed_disk_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"storage_account_identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"tenant_id": {
							Type:      pluginsdk.TypeString,
							Sensitive: true,
							Computed:  true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func dataSourceDatabricksWorkspaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.WorkspacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := workspaces.NewWorkspaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.WorkspaceName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model := resp.Model; model != nil {
		if sku := resp.Model.Sku; sku != nil {
			d.Set("sku", sku.Name)
		}
		d.Set("workspace_id", model.Properties.WorkspaceId)
		if err := d.Set("storage_account_identity", flattenWorkspaceManagedIdentity(model.Properties.StorageAccountIdentity)); err != nil {
			return fmt.Errorf("setting `storage_account_identity`: %+v", err)
		}
		if err := d.Set("managed_disk_identity", flattenWorkspaceManagedIdentity(model.Properties.ManagedDiskIdentity)); err != nil {
			return fmt.Errorf("setting `managed_disk_identity`: %+v", err)
		}
		d.Set("workspace_url", model.Properties.WorkspaceUrl)
		d.Set("location", model.Location)

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}
