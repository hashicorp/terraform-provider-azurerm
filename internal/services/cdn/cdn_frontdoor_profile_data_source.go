// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceCdnFrontDoorProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorProfileRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"response_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCdnFrontDoorProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	profileId := profiles.ProfileId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroup,
		ProfileName:       id.ProfileName,
	}

	resp, err := client.Get(ctx, profileId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	model := resp.Model

	if model == nil {
		return fmt.Errorf("model is 'nil'")
	}

	if model.Properties == nil {
		return fmt.Errorf("model.Properties is 'nil'")
	}

	d.SetId(id.ID())
	d.Set("name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("response_timeout_seconds", int(pointer.From(model.Properties.OriginResponseTimeoutSeconds)))

	// whilst this is returned in the API as FrontDoorID other resources refer to
	// this as the Resource GUID, so we will for consistency
	d.Set("resource_guid", pointer.From(model.Properties.FrontDoorId))

	identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
	if err == nil {
		d.Set("identity", identity)
	}

	skuName := ""
	if model.Sku.Name != nil {
		skuName = string(pointer.From(model.Sku.Name))
	}

	d.Set("sku_name", skuName)
	d.Set("tags", flattenNewFrontDoorTags(model.Tags))

	return nil
}
