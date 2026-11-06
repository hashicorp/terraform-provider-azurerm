// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_origin_group

import (
	"context"
	"fmt"
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

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(cdnFrontDoorOriginGroupAuthenticationCustomizeDiff),
		),

		Schema: cdnFrontdoorOriginGroupSchema(),
	}
}

func cdnFrontDoorOriginGroupAuthenticationCustomizeDiff(_ context.Context, diff *pluginsdk.ResourceDiff, _ any) error {
	authentication := diff.Get("authentication").([]any)
	if len(authentication) == 0 || authentication[0] == nil {
		return nil
	}

	v := authentication[0].(map[string]any)
	authenticationType := v["type"].(string)
	userAssignedIdentityId := v["user_assigned_identity_id"].(string)

	// Skip the check when `user_assigned_identity_id` is `(known after apply)`, e.g. the identity is created in the same apply
	userAssignedIdentityIdKnown := true
	if rawConfig := diff.GetRawConfig(); !rawConfig.IsNull() {
		authenticationConfig := rawConfig.GetAttr("authentication")
		if authenticationConfig.IsKnown() && !authenticationConfig.IsNull() && authenticationConfig.LengthInt() > 0 {
			authenticationBlock := authenticationConfig.AsValueSlice()[0]
			if !authenticationBlock.IsNull() {
				userAssignedIdentityIdKnown = authenticationBlock.GetAttr("user_assigned_identity_id").IsKnown()
			}
		}
	}

	if authenticationType == string(afdorigingroups.OriginAuthenticationTypeUserAssignedIdentity) && userAssignedIdentityId == "" && userAssignedIdentityIdKnown {
		return fmt.Errorf("`authentication.0.user_assigned_identity_id` must be specified when `authentication.0.type` is `%s`", afdorigingroups.OriginAuthenticationTypeUserAssignedIdentity)
	}

	if authenticationType == string(afdorigingroups.OriginAuthenticationTypeSystemAssignedIdentity) && userAssignedIdentityId != "" {
		return fmt.Errorf("`authentication.0.user_assigned_identity_id` cannot be specified when `authentication.0.type` is `%s`", afdorigingroups.OriginAuthenticationTypeSystemAssignedIdentity)
	}

	return nil
}
