// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	helperTags "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataShareAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName(),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceDataShareAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if err := d.Set("identity", identity.FlattenSystemAssigned(&model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		return helperTags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
