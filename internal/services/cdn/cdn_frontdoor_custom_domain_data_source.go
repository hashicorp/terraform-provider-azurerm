// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCdnFrontDoorCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCdnFrontDoorCustomDomainRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorCustomDomainName,
			},

			"profile_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cdn_frontdoor_profile_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tls": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"minimum_tls_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"cdn_frontdoor_secret_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validation_token": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewFrontDoorCustomDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.CustomDomainName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("profile_name", id.ProfileName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDDomainProperties; props != nil {
		d.Set("host_name", props.HostName)

		dnsZoneId, err := flattenDNSZoneResourceReference(props.AzureDNSZone)
		if err != nil {
			return fmt.Errorf("flattening `dns_zone_id`: %+v", err)
		}

		if err := d.Set("dns_zone_id", dnsZoneId); err != nil {
			return fmt.Errorf("setting `dns_zone_id`: %+v", err)
		}

		tls, err := flattenCustomDomainAFDDomainHttpsParameters(props.TLSSettings)
		if err != nil {
			return fmt.Errorf("flattening `tls`: %+v", err)
		}

		if err := d.Set("tls", tls); err != nil {
			return fmt.Errorf("setting `tls`: %+v", err)
		}

		if validationProps := props.ValidationProperties; validationProps != nil {
			d.Set("expiration_date", validationProps.ExpirationDate)
			d.Set("validation_token", validationProps.ValidationToken)
		}
	}

	return nil
}
