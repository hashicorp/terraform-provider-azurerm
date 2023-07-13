// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorRuleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorRuleSetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorRuleSetName,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cdn_frontdoor_profile_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCdnFrontDoorRuleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRuleSetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorRuleSetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.RuleSetName)
	d.Set("profile_name", id.ProfileName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	return nil
}
