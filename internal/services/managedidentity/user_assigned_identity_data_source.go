// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedidentity

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmUserAssignedIdentity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmUserAssignedIdentityRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"principal_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmUserAssignedIdentityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedIdentity.V20230131.ManagedIdentities
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewUserAssignedIdentityID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.UserAssignedIdentitiesGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.ClientId)
			d.Set("principal_id", props.PrincipalId)
			d.Set("tenant_id", props.TenantId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
