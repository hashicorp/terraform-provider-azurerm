// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceContainerRegistry() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceContainerRegistryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceContainerRegistrySchema(),
	}
}

func dataSourceContainerRegistryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient_v2023_06_01_preview.Registries
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := registries.NewRegistryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.RegistryName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if model.Sku.Tier != nil {
			d.Set("sku", string(*model.Sku.Tier))
		}

		if props := model.Properties; props != nil {
			d.Set("admin_enabled", props.AdminUserEnabled)
			d.Set("login_server", props.LoginServer)
			d.Set("data_endpoint_enabled", props.DataEndpointEnabled)

			if *props.AdminUserEnabled {
				credsResp, err := client.ListCredentials(ctx, id)
				if err != nil {
					return fmt.Errorf("retrieving credentials for %s: %s", id, err)
				}

				if credsModel := credsResp.Model; credsModel != nil {
					d.Set("admin_username", credsModel.Username)
					for _, v := range *credsModel.Passwords {
						d.Set("admin_password", v.Value)
						break
					}
				} else {
					d.Set("admin_username", "")
					d.Set("admin_password", "")
				}
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func dataSourceContainerRegistrySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"location": commonschema.LocationComputed(),

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"admin_password": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}
