// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnFrontDoorCustomDomainCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := afddomains.NewCustomDomainID(subscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.AFDCustomDomainsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
		}
	}

	dnsZone := d.Get("dns_zone_id").(string)
	tls := d.Get("tls").([]any)

	props := afddomains.AFDDomain{
		Properties: &afddomains.AFDDomainProperties{
			HostName: d.Get("host_name").(string),
		},
	}

	if dnsZone != "" {
		props.Properties.AzureDnsZone = expandAfdResourceReference(dnsZone)
	}

	tlsSettings, err := expandAfdDomainTlsParameters(d, tls)
	if err != nil {
		return err
	}

	props.Properties.TlsSettings = tlsSettings

	if err := client.AFDCustomDomainsCreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}
