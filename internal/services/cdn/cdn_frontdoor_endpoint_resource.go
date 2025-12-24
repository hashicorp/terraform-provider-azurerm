// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorEndpointCreate,
		Read:   resourceCdnFrontDoorEndpointRead,
		Update: resourceCdnFrontDoorEndpointUpdate,
		Delete: resourceCdnFrontDoorEndpointDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorEndpointID(id)
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
				ValidateFunc: validate.FrontDoorProfileID,
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

func resourceCdnFrontDoorEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileRaw := d.Get("cdn_frontdoor_profile_id").(string)
	profileId, err := parse.FrontDoorProfileID(profileRaw)
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorEndpointID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_endpoint", id.ID())
	}

	props := cdn.AFDEndpoint{
		Name:     pointer.To(d.Get("name").(string)),
		Location: pointer.To(location.Normalize("global")),
		AFDEndpointProperties: &cdn.AFDEndpointProperties{
			EnabledState: expandEnabledBool(d.Get("enabled").(bool)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorEndpointRead(d, meta)
}

func resourceCdnFrontDoorEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.AfdEndpointName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDEndpointProperties; props != nil {
		d.Set("enabled", flattenEnabledBool(props.EnabledState))
		d.Set("host_name", props.HostName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCdnFrontDoorEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	props := cdn.AFDEndpointUpdateParameters{}

	if d.HasChange("enabled") {
		props.AFDEndpointPropertiesUpdateParameters = &cdn.AFDEndpointPropertiesUpdateParameters{
			EnabledState: expandEnabledBool(d.Get("enabled").(bool)),
		}
	}

	if d.HasChange("tags") {
		props.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorEndpointRead(d, meta)
}

func resourceCdnFrontDoorEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorEndpointID(d.Id())
	if err != nil {
		return err
	}

	// Before deleting the endpoint, remove it from any security policy associations
	// This is necessary because Azure will reject deletion of an endpoint that is
	// still associated with a security policy
	if err := removeEndpointFromSecurityPolicies(ctx, meta, id); err != nil {
		return fmt.Errorf("removing endpoint from security policies before deletion: %+v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func removeEndpointFromSecurityPolicies(ctx context.Context, meta interface{}, endpointId *parse.FrontDoorEndpointId) error {
	securityPoliciesClient := meta.(*clients.Client).Cdn.FrontDoorSecurityPoliciesClient

	profileId := securitypolicies.NewProfileID(endpointId.SubscriptionId, endpointId.ResourceGroup, endpointId.ProfileName)

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
