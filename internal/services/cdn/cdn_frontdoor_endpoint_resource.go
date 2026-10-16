// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorEndpointCreate,
		Read:   resourceCdnFrontDoorEndpointRead,
		Update: resourceCdnFrontDoorEndpointUpdate,
		Delete: resourceCdnFrontDoorEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := afdendpoints.ParseAfdEndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorEndpointName,
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: afdendpoints.ValidateProfileID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"tags": commonschema.Tags(),

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontDoorEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := afdendpoints.NewAfdEndpointID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_endpoint", id.ID())
		}
	}

	enabledState := afdendpoints.EnabledStateEnabled
	if !d.Get("enabled").(bool) {
		enabledState = afdendpoints.EnabledStateDisabled
	}

	props := afdendpoints.AFDEndpoint{
		Name:     pointer.To(id.AfdEndpointName),
		Location: "global",
		Properties: &afdendpoints.AFDEndpointProperties{
			EnabledState: &enabledState,
		},

		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorEndpointRead(d, meta)
}

func resourceCdnFrontDoorEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AfdEndpointName)
	d.Set("cdn_frontdoor_profile_id", afdendpoints.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			d.Set("enabled", pointer.From(props.EnabledState) == afdendpoints.EnabledStateEnabled)
			d.Set("host_name", props.HostName)
		}
	}

	return nil
}

func resourceCdnFrontDoorEndpointUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	// TODO: Check if this can be migrated to a PUT upgrade path for consistency with the rest of the provider
	props := afdendpoints.AFDEndpointUpdateParameters{}

	if d.HasChange("enabled") {
		enabledState := afdendpoints.EnabledStateEnabled
		if !d.Get("enabled").(bool) {
			enabledState = afdendpoints.EnabledStateDisabled
		}

		props.Properties = &afdendpoints.AFDEndpointPropertiesUpdateParameters{
			EnabledState: &enabledState,
		}
	}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorEndpointRead(d, meta)
}

func resourceCdnFrontDoorEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointsClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afdendpoints.ParseAfdEndpointID(d.Id())
	if err != nil {
		return err
	}

	// Before deleting the endpoint, remove it from any security policy associations
	// This is necessary because Azure will reject deletion of an endpoint that is
	// still associated with a security policy
	if err := removeEndpointFromSecurityPolicies(ctx, meta, id); err != nil {
		return fmt.Errorf("removing endpoint from security policies before deletion: %+v", err)
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func removeEndpointFromSecurityPolicies(ctx context.Context, meta any, endpointId *afdendpoints.AfdEndpointId) error {
	securityPoliciesClient := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient

	profileId := securitypolicies.NewProfileID(endpointId.SubscriptionId, endpointId.ResourceGroupName, endpointId.ProfileName)

	// List all security policies for the profile
	resp, err := securityPoliciesClient.ListByProfileComplete(ctx, profileId)
	if err != nil {
		return fmt.Errorf("listing security policies for profile %s: %+v", profileId, err)
	}

	endpointIdStr := endpointId.ID()

	for _, policy := range resp.Items {
		if policy.Properties == nil || policy.Properties.Parameters == nil {
			continue
		}

		// Check if this is a WAF policy
		if policy.Properties.Parameters.SecurityPolicyPropertiesParameters().Type != securitypolicies.SecurityPolicyTypeWebApplicationFirewall {
			continue
		}

		wafParams, ok := policy.Properties.Parameters.(securitypolicies.SecurityPolicyWebApplicationFirewallParameters)
		if !ok || wafParams.Associations == nil {
			continue
		}

		// Check if this endpoint is referenced in any association
		endpointFound := false
		for _, assoc := range *wafParams.Associations {
			if assoc.Domains == nil {
				continue
			}
			for _, domain := range *assoc.Domains {
				if domain.Id != nil && strings.EqualFold(*domain.Id, endpointIdStr) {
					endpointFound = true
					break
				}
			}
			if endpointFound {
				break
			}
		}

		if !endpointFound {
			continue
		}

		// Remove the endpoint from the security policy associations
		newAssociations := make([]securitypolicies.SecurityPolicyWebApplicationFirewallAssociation, 0)
		for _, assoc := range *wafParams.Associations {
			if assoc.Domains == nil {
				newAssociations = append(newAssociations, assoc)
				continue
			}

			newDomains := make([]securitypolicies.ActivatedResourceReference, 0)
			for _, domain := range *assoc.Domains {
				if domain.Id != nil && strings.EqualFold(*domain.Id, endpointIdStr) {
					// Skip this endpoint - we're removing it
					continue
				}
				newDomains = append(newDomains, domain)
			}

			// Only include the association if it still has domains
			if len(newDomains) > 0 {
				newAssociations = append(newAssociations, securitypolicies.SecurityPolicyWebApplicationFirewallAssociation{
					Domains:         &newDomains,
					PatternsToMatch: assoc.PatternsToMatch,
				})
			}
		}

		// Update the security policy with the endpoint removed
		// Use case-insensitive parsing because Azure may return IDs with different casing
		policyId, err := securitypolicies.ParseSecurityPolicyIDInsensitively(*policy.Id)
		if err != nil {
			return fmt.Errorf("parsing security policy ID %s: %+v", *policy.Id, err)
		}

		updatedParams := securitypolicies.SecurityPolicyWebApplicationFirewallParameters{
			Associations: &newAssociations,
			WafPolicy:    wafParams.WafPolicy,
		}

		updatedPolicy := securitypolicies.SecurityPolicy{
			Properties: &securitypolicies.SecurityPolicyProperties{
				Parameters: updatedParams,
			},
		}

		if err := securityPoliciesClient.CreateThenPoll(ctx, *policyId, updatedPolicy); err != nil {
			return fmt.Errorf("updating security policy %s to remove endpoint association: %+v", *policyId, err)
		}
	}

	return nil
}
