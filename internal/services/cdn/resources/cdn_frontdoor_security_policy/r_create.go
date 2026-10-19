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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontdoorSecurityPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// NOTE: The profile id is used to retrieve properties from the related profile that must match in this security policy
	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	securityPolicyName := d.Get("name").(string)
	id := securitypolicies.NewSecurityPolicyID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, securityPolicyName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_security_policy", id.ID())
		}
	}

	profileClient := meta.(*clients.Client).Cdn.FrontDoorProfilesClient
	resp, err := profileClient.Get(ctx, pointer.From(profileId))
	if err != nil {
		return fmt.Errorf("unable to retrieve the 'sku_name' from the CDN FrontDoor Profile(Name: %q)': %+v", profileId.ProfileName, err)
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

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorSecurityPolicyRead(d, meta)
}

func expandCdnFrontdoorFirewallPolicyParameters(input []any, isStandardSku bool) (*securitypolicies.SecurityPolicyWebApplicationFirewallParameters, error) {
	results := securitypolicies.SecurityPolicyWebApplicationFirewallParameters{}
	if len(input) == 0 {
		return &results, nil
	}

	associations := make([]securitypolicies.SecurityPolicyWebApplicationFirewallAssociation, 0)

	// pull off only the firewall policy from the security_policies list
	policyType := input[0].(map[string]any)
	firewallPolicy := policyType["firewall"].([]any)
	v := firewallPolicy[0].(map[string]any)

	if id := v["cdn_frontdoor_firewall_policy_id"].(string); id != "" {
		results.WafPolicy = &securitypolicies.ResourceReference{
			Id: pointer.To(id),
		}
	}

	configAssociations := v["association"].([]any)

	for _, item := range configAssociations {
		v := item.(map[string]any)
		domains := expandSecurityPoliciesActivatedResourceReference(v["domain"].([]any))
		domainCount := len(*domains)

		if isStandardSku {
			if domainCount > 100 {
				return &results, fmt.Errorf("the 'Standard_AzureFrontDoor' sku is only allowed to have 100 or less domains associated with the firewall policy, got %d", domainCount)
			}
		} else {
			if domainCount > 500 {
				return &results, fmt.Errorf("the 'Premium_AzureFrontDoor' sku is only allowed to have 500 or less domains associated with the firewall policy, got %d", domainCount)
			}
		}

		association := securitypolicies.SecurityPolicyWebApplicationFirewallAssociation{
			Domains:         domains,
			PatternsToMatch: pluginsdk.ExpandStringSlice(v["patterns_to_match"].([]any)),
		}

		associations = append(associations, association)
	}

	results.Associations = &associations

	return &results, nil
}

func expandSecurityPoliciesActivatedResourceReference(input []any) *[]securitypolicies.ActivatedResourceReference {
	results := make([]securitypolicies.ActivatedResourceReference, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]any)

		if id := v["cdn_frontdoor_domain_id"].(string); id != "" {
			results = append(results, securitypolicies.ActivatedResourceReference{
				Id: pointer.To(id),
			})
		}
	}

	return &results
}
