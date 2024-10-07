// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePublicIPs() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePublicIPsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourcePublicIPSchema(),
	}
}

func dataSourcePublicIPSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

		"name_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"attachment_status": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Attached",
				"Unattached",
			}, false),
		},

		"allocation_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(publicipaddresses.IPAllocationMethodDynamic),
				string(publicipaddresses.IPAllocationMethodStatic),
			}, false),
		},

		"public_ips": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"domain_name_label": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourcePublicIPsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPAddresses
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupId := commonids.NewResourceGroupID(subscriptionId, d.Get("resource_group_name").(string))

	resp, err := client.List(ctx, resourceGroupId)
	if err != nil {
		return fmt.Errorf("listing Public IP Addresses in %s: %v", resourceGroupId, err)
	}

	prefix := d.Get("name_prefix").(string)
	attachmentStatus, attachmentStatusOk := d.GetOk("attachment_status")
	allocationType := d.Get("allocation_type").(string)

	filteredIPAddresses := make([]publicipaddresses.PublicIPAddress, 0)

	if model := resp.Model; model != nil {
		for _, address := range *model {
			if props := address.Properties; props != nil {
				nicIsAttached := props.IPConfiguration != nil || props.NatGateway != nil

				if prefix != "" {
					if !strings.HasPrefix(*address.Name, prefix) {
						continue
					}
				}

				if attachmentStatusOk && attachmentStatus.(string) == "Attached" && !nicIsAttached {
					continue
				}
				if attachmentStatusOk && attachmentStatus.(string) == "Unattached" && nicIsAttached {
					continue
				}

				if allocationType != "" {
					allocation := publicipaddresses.IPAllocationMethod(allocationType)
					if props.PublicIPAllocationMethod != nil && *props.PublicIPAllocationMethod != allocation {
						continue
					}
				}
			}

			filteredIPAddresses = append(filteredIPAddresses, address)
		}
	}

	id := fmt.Sprintf("networkPublicIPs/resourceGroup/%s/namePrefix=%s;attachmentStatus=%s;allocationType=%s", resourceGroupId.ResourceGroupName, prefix, attachmentStatus, allocationType)
	d.SetId(base64.StdEncoding.EncodeToString([]byte(id)))

	results := flattenDataSourcePublicIPs(filteredIPAddresses)
	if err := d.Set("public_ips", results); err != nil {
		return fmt.Errorf("setting `public_ips`: %+v", err)
	}

	return nil
}

func flattenDataSourcePublicIPs(input []publicipaddresses.PublicIPAddress) []interface{} {
	results := make([]interface{}, 0)

	for _, element := range input {
		flattenedIPAddress := flattenDataSourcePublicIP(element)
		results = append(results, flattenedIPAddress)
	}

	return results
}

func flattenDataSourcePublicIP(input publicipaddresses.PublicIPAddress) map[string]string {
	id := ""
	if input.Id != nil {
		id = *input.Id
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	domainNameLabel := ""
	fqdn := ""
	ipAddress := ""
	if props := input.Properties; props != nil {
		if dns := props.DnsSettings; dns != nil {
			if dns.Fqdn != nil {
				fqdn = *dns.Fqdn
			}

			if dns.DomainNameLabel != nil {
				domainNameLabel = *dns.DomainNameLabel
			}
		}

		if props.IPAddress != nil {
			ipAddress = *props.IPAddress
		}
	}

	return map[string]string{
		"id":                id,
		"name":              name,
		"domain_name_label": domainNameLabel,
		"fqdn":              fqdn,
		"ip_address":        ipAddress,
	}
}
