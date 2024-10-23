// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceContainerRegistryScopeMap() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceContainerRegistryScopeMapRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"container_registry_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ContainerRegistryName,
			},
			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceContainerRegistryScopeMapRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2023_06_01_preview.ScopeMaps
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scopemaps.NewScopeMapID(subscriptionId, d.Get("resource_group_name").(string), d.Get("container_registry_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.ScopeMapName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("container_registry_name", id.RegistryName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			description := ""
			if v := props.Description; v != nil {
				description = *v
			}
			d.Set("description", description)
			d.Set("actions", utils.FlattenStringSlice(&props.Actions))
		}
	}
	return nil
}
