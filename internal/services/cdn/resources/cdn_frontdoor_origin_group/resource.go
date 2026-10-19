// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:cdn_frontdoor_profile_id,resource_group_name:cdn_frontdoor_profile_id,profile_name:cdn_frontdoor_profile_id"

const azureCdnFrontDoorOriginGroupResourceName = "azurerm_cdn_frontdoor_origin_group"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginGroupCreate,
		Read:   resourceCdnFrontDoorOriginGroupRead,
		Update: resourceCdnFrontDoorOriginGroupUpdate,
		Delete: resourceCdnFrontDoorOriginGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&afdorigingroups.OriginGroupId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&afdorigingroups.OriginGroupId{}),

		Schema: cdnFrontdoorOriginGroupSchema(),
	}
}
