// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	cdnClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Gets information about an existing Front Door (standard/premium) Custom Domain.",
		Read:        dataSourceCdnFrontDoorCustomDomainRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Front Door Custom Domain.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorCustomDomainName,
			},

			"profile_name": {
				Description:  "The name of the Front Door Profile which the Front Door Custom Domain is bound to.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cdn_frontdoor_profile_id": {
				Description: "The ID of the Front Door Profile which the Front Door Custom Domain is bound to.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"dns_zone_id": {
				Description: "The ID of the Azure DNS Zone which should be used for this Front Door Custom Domain.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"host_name": {
				Description: "The host name of the domain.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tls": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_type": {
							Description: "The SSL certificate type.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"minimum_version": {
							Description: "The TLS protocol version that will be used for HTTPS connections.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"cdn_frontdoor_secret_id": {
							Description: "The Resource ID of the Front Door Secret.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"cipher_suite": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"type": {
										Description: "The TLS policy type used for this Front Door Custom Domain.",
										Type:        pluginsdk.TypeString,
										Computed:    true,
									},

									"custom_ciphers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"tls12": {
													Description: "A set of TLS 1.2 cipher suites.",
													Type:        pluginsdk.TypeSet,
													Computed:    true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
													},
												},

												"tls13": {
													Description: "A set of TLS 1.3 cipher suites.",
													Type:        pluginsdk.TypeSet,
													Computed:    true,
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
				Description: "The date and time that the token expires.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"validation_token": {
				Description: "The challenge used for DNS TXT record or file based validation.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta any) error {
	client := cdnClient.From(meta).AFDCustomDomainsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := afddomains.NewCustomDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))

	resp, err := client.AFDCustomDomainsGet(ctx, id)
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

			if err := d.Set("dns_zone_id", flattenAfdDNSZoneResourceReference(props.AzureDnsZone)); err != nil {
				return fmt.Errorf("setting `dns_zone_id`: %+v", err)
			}

			if err := d.Set("tls", flattenAfdDomainHttpsParameters(props.TlsSettings, true)); err != nil {
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
