// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceEventGridDomainUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).EventGrid.Domains
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := domains.ParseDomainID(d.Id())
	if err != nil {
		return err
	}

	payload := domains.DomainUpdateParameters{Properties: &domains.DomainUpdateParameterProperties{}}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := domains.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = domains.PublicNetworkAccessEnabled
		}

		payload.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}

	if d.HasChange("local_auth_enabled") {
		payload.Properties.DisableLocalAuth = pointer.To(!d.Get("local_auth_enabled").(bool))
	}

	if d.HasChange("auto_create_topic_with_first_subscription") {
		payload.Properties.AutoCreateTopicWithFirstSubscription = pointer.To(d.Get("auto_create_topic_with_first_subscription").(bool))
	}

	if d.HasChange("auto_delete_topic_with_last_subscription") {
		payload.Properties.AutoDeleteTopicWithLastSubscription = pointer.To(d.Get("auto_delete_topic_with_last_subscription").(bool))
	}

	if d.HasChange("data_residency_boundary") {
		payload.Properties.DataResidencyBoundary = pointer.ToEnum[domains.DataResidencyBoundary](d.Get("data_residency_boundary").(string))
	}

	if d.HasChange("minimum_tls_version") {
		payload.Properties.MinimumTlsVersionAllowed = pointer.ToEnum[domains.TlsVersion](d.Get("minimum_tls_version").(string))
	}

	if d.HasChange("inbound_ip_rule") {
		inboundIpRule := d.Get("inbound_ip_rule").([]any)

		if len(inboundIpRule) == 0 {
			payload.Properties.InboundIPRules = pointer.To([]domains.InboundIPRule{})
		} else {
			payload.Properties.InboundIPRules = expandDomainInboundIPRules(inboundIpRule)
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %s", id, err)
	}

	return resourceEventGridDomainRead(d, meta)
}
