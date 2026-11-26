// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Container Registry.",
		Read:        dataSourceContainerRegistryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceContainerRegistrySchema(),
	}
}

func dataSourceContainerRegistryRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.ContainerRegistryClient.Registries
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
			d.Set("data_endpoint_host_names", props.DataEndpointHostNames)

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
			Description:  "The name of the Container Registry.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"location": commonschema.LocationComputed(),

		"admin_enabled": {
			Description: "Whether the Administrator account is enabled for this Container Registry.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"admin_password": {
			Description: "The Password associated with the Container Registry Admin account - if the admin account is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"admin_username": {
			Description: "The Username associated with the Container Registry Admin account - if the admin account is enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"data_endpoint_enabled": {
			Description: "Whether dedicated data endpoints for this Container Registry are enabled.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"data_endpoint_host_names": {
			Description: "A set of data endpoint hostnames associated with the container registry if data endpoints are enabled.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"login_server": {
			Description: "The URL that can be used to log into the container registry.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"sku": {
			Description: "The SKU of this Container Registry, such as `Basic`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}
