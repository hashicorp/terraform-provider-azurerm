// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/resources/cdn_frontdoor_origin/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:cdn_frontdoor_origin_group_id,resource_group_name:cdn_frontdoor_origin_group_id,profile_name:cdn_frontdoor_origin_group_id,origin_group_name:cdn_frontdoor_origin_group_id"

const azureCdnFrontDoorOriginResourceName = "azurerm_cdn_frontdoor_origin"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorOriginCreate,
		Read:   resourceCdnFrontDoorOriginRead,
		Update: resourceCdnFrontDoorOriginUpdate,
		Delete: resourceCdnFrontDoorOriginDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&afdorigins.OriginGroupOriginId{}),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&afdorigins.OriginGroupOriginId{}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnFrontDoorOriginV0ToV1{},
		}),

		Schema: cdnFrontdoorOriginSchema(),
	}
}
