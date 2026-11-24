// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNetworkServiceTags() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about Service Tags.",
		Read:        dataSourceNetworkServiceTagsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"service": {
				Description: "The type of the service for which address prefixes will be fetched. Available service tags can be found here: [Available service tags](https://learn.microsoft.com/azure/virtual-network/service-tags-overview#available-service-tags).",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"name": {
				Description: "The name of this Service Tags block.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"location_filter": commonschema.LocationOptional(),

			"address_prefixes": {
				Description: "List of address prefixes for the service type (and optionally a specific region).",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"ipv4_cidrs": {
				Description: "List of IPv4 addresses for the service type (and optionally a specific region).",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"ipv6_cidrs": {
				Description: "List of IPv6 addresses for the service type (and optionally a specific region).",
				Type:        pluginsdk.TypeList,
				Computed:    true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceNetworkServiceTagsRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	locationId := virtualnetworks.NewLocationID(subscriptionId, location.Normalize(d.Get("location").(string)))
	resp, err := client.ServiceTagsList(ctx, locationId)
	if err != nil {
		return fmt.Errorf("listing network service tags: %+v", err)
	}

	if resp.Model == nil || resp.Model.Values == nil {
		return fmt.Errorf("listing network service tags: `model.Values` was nil")
	}

	service := d.Get("service").(string)
	locationFilter := location.Normalize(d.Get("location_filter").(string))
	for _, value := range *resp.Model.Values {
		if value.Name == nil || !isServiceTagOf(*value.Name, service) {
			continue
		}

		if props := value.Properties; props != nil {
			if props.Region == nil {
				continue
			}

			if location.NormalizeNilable(props.Region) == locationFilter {
				d.Set("name", value.Name)

				addressPrefixes := make([]string, 0)
				if props.AddressPrefixes != nil {
					addressPrefixes = *props.AddressPrefixes
				}
				if err = d.Set("address_prefixes", addressPrefixes); err != nil {
					return fmt.Errorf("setting `address_prefixes`: %+v", err)
				}

				var ipv4Cidrs []string
				var ipv6Cidrs []string

				for _, prefix := range addressPrefixes {
					ip, ipNet, err := net.ParseCIDR(prefix)
					if err != nil {
						return err
					}

					if ip.To4() != nil {
						ipv4Cidrs = append(ipv4Cidrs, ipNet.String())
					} else {
						ipv6Cidrs = append(ipv6Cidrs, ipNet.String())
					}
				}

				if err = d.Set("ipv4_cidrs", ipv4Cidrs); err != nil {
					return fmt.Errorf("setting `ipv4_cidrs`: %+v", err)
				}

				if err = d.Set("ipv6_cidrs", ipv6Cidrs); err != nil {
					return fmt.Errorf("setting `ipv6_cidrs`: %+v", err)
				}

				d.SetId(fmt.Sprintf("%s-%s", locationId.LocationName, service))

				return nil
			}
		}
	}

	errSuffix := "globally"
	if locationFilter != "" {
		errSuffix = "for region " + locationFilter
	}
	return fmt.Errorf("specified service tag `%s` not found %s", service, errSuffix)
}

// isServiceTagOf is used to check whether a service tag name belongs to the service of name `serviceName`.
// Service tag name has format as below:
// - (regional) serviceName.locationName
// - (all) serviceName
func isServiceTagOf(stName, serviceName string) bool {
	stNameComponents := strings.Split(stName, ".")
	if len(stNameComponents) != 1 && len(stNameComponents) != 2 {
		return false
	}

	// Adds support for "AzureFrontDoor.Frontend" service tag
	serviceNameComponents := strings.Split(serviceName, ".")
	if len(serviceNameComponents) != 1 && len(serviceNameComponents) != 2 {
		return false
	}

	if len(serviceNameComponents) == 1 {
		return stNameComponents[0] == serviceNameComponents[0]
	}
	return stNameComponents[0] == serviceNameComponents[0] && stNameComponents[1] == serviceNameComponents[1]
}
