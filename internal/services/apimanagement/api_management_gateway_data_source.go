// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceApiManagementGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildDataSourceName(),

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location_data": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"city": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"district": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"region": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApiManagementGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gateway.NewGatewayID(apimId.SubscriptionId, apimId.ResourceGroupName, apimId.ServiceName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	if resp.Model != nil {
		_, err = gateway.ParseGatewayID(*resp.Model.Id)
		if err != nil {
			return fmt.Errorf("parsing Gateway ID %q", *resp.Model.Id)
		}
	}
	d.SetId(id.ID())

	d.Set("api_management_id", apimId.ID())

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("location_data", flattenApiManagementGatewayLocationData(props.LocationData))
		}
	}

	return nil
}
