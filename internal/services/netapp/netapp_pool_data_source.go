// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNetAppPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetAppPoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PoolName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"service_level": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"size_in_tb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"encryption_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"cool_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceNetAppPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := capacitypools.NewCapacityPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))
	resp, err := client.PoolsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.CapacityPoolName)
	d.Set("account_name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		d.Set("service_level", string(props.ServiceLevel))
		d.Set("size_in_tb", props.Size/1099511627776)
		d.Set("encryption_type", string(pointer.From(props.EncryptionType)))
		d.Set("cool_access_enabled", pointer.From(props.CoolAccess))
	}

	return nil
}
