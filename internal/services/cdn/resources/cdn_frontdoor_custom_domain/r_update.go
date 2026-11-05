// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := afddomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	props := afddomains.AFDDomainUpdateParameters{
		Properties: &afddomains.AFDDomainUpdatePropertiesParameters{},
	}

	if d.HasChange("dns_zone_id") {
		if dnsZone := d.Get("dns_zone_id").(string); dnsZone != "" {
			props.Properties.AzureDnsZone = expandAfdResourceReference(dnsZone)
		}
	}

	if d.HasChange("tls") {
		tlsSettings := d.Get("tls").([]any)
		tls, err := expandAfdDomainTlsParameters(d, tlsSettings)
		if err != nil {
			return err
		}
		props.Properties.TlsSettings = tls
	}

	// The AFDx service currently rejects custom-domain updates until its internal
	// validation and backend synchronization flow completes. The service team
	// confirmed clients must wait for approval here, otherwise Update can return
	// 409 Conflict while the sync job is still in progress.
	approvalPollerType := custompollers.NewFrontDoorCustomDomainWaitForApprovedPoller(client, *id)
	approvalPoller := pollers.NewPoller(approvalPollerType, 30*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := approvalPoller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be approved: %+v", *id, err)
	}

	if err := client.AFDCustomDomainsUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}
