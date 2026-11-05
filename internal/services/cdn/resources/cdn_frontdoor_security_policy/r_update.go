// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_security_policy

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontdoorSecurityPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// NOTE: The profile id is used to retrieve properties from the related profile that must match in this security policy
	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := securitypolicies.NewSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	resp, err := profileClient.Get(ctx, pointer.From(profileId))
	if err != nil {
		return fmt.Errorf("unable to retrieve the `sku_name` from %s: %+v", pointer.From(profileId), err)
	}

	profileModel := resp.Model

	if profileModel == nil {
		return errors.New("profileModel is 'nil'")
	}

	if profileModel.Sku.Name == nil {
		return errors.New("profileModel.Sku.Name is 'nil'")
	}

	params, err := expandCdnFrontdoorFirewallPolicyParameters(d.Get("security_policies").([]any), strings.EqualFold(pointer.FromEnum(profileModel.Sku.Name), string(profiles.SkuNameStandardAzureFrontDoor)))
	if err != nil {
		return fmt.Errorf("expanding 'security_policies': %+v", err)
	}

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: params,
		},
	}

	// Using 'Create' for update because it is a PUT operation
	if err = client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCdnFrontdoorSecurityPolicyRead(d, meta)
}
