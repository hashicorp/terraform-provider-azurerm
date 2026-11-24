// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_group

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing API Management Group.",
		Read:        dataSourceApiManagementGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildDataSourceName("The Name of the API Management Group."),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"display_name": {
				Description: "The display name of this API Management Group.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"description": {
				Description: "The description of this API Management Group.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"external_id": {
				Description: "The identifier of the external Group.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"type": {
				Description: "The type of this API Management Group, such as `custom` or `external`.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceApiManagementGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)

	id := group.NewGroupID(meta.(*clients.Client).Account.SubscriptionId, resourceGroup, serviceName, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.GroupId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("description", pointer.From(props.Description))
			d.Set("external_id", pointer.From(props.ExternalId))
			d.Set("type", pointer.FromEnum(props.Type))
		}
	}

	return nil
}
