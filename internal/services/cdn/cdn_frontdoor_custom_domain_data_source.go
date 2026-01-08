// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-04-15/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

						"cipher_suite": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"custom_ciphers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"tls12": {
													Type:     pluginsdk.TypeSet,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},

												"tls13": {
													Type:     pluginsdk.TypeSet,
													Computed: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},
											},
										},
									},
								},
							},
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
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := afdcustomdomains.NewCustomDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.CustomDomainName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("profile_name", id.ProfileName)
	d.Set("cdn_frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("host_name", props.HostName)

			dnsZoneId := flattenAfdDNSZoneResourceReference(props.AzureDnsZone)
			if err := d.Set("dns_zone_id", dnsZoneId); err != nil {
				return fmt.Errorf("setting `dns_zone_id`: %+v", err)
			}

			tls := flattenAfdDomainHttpsParameters(props.TlsSettings)
			if err := d.Set("tls", tls); err != nil {
				return fmt.Errorf("setting `tls`: %+v", err)
			}

			if validationProps := props.ValidationProperties; validationProps != nil {
				d.Set("expiration_date", validationProps.ExpirationDate)
				d.Set("validation_token", validationProps.ValidationToken)
			}
		}
	}

	return nil
}
