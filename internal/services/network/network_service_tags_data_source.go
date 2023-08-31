// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/servicetags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceNetworkServiceTags() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetworkServiceTagsRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"location": commonschema.Location(),

			"service": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location_filter": commonschema.LocationOptional(),

			"address_prefixes": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"ipv4_cidrs": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"ipv6_cidrs": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceNetworkServiceTagsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceTags
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	locationId := servicetags.NewLocationID(subscriptionId, location.Normalize(d.Get("location").(string)))
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
				addressPrefixes := make([]string, 0)
				if props.AddressPrefixes != nil {
					addressPrefixes = *props.AddressPrefixes
				}
				err = d.Set("address_prefixes", addressPrefixes)
				if err != nil {
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

				err = d.Set("ipv4_cidrs", ipv4Cidrs)
				if err != nil {
					return fmt.Errorf("setting `ipv4_cidrs`: %+v", err)
				}

				err = d.Set("ipv6_cidrs", ipv6Cidrs)
				if err != nil {
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
	return stNameComponents[0] == serviceName
}
