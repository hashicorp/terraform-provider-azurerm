package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
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
			"location": azure.SchemaLocation(),

			"service": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location_filter": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        azure.NormalizeLocation,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

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
	client := meta.(*clients.Client).Network.ServiceTagsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location"))
	res, err := client.List(ctx, location)
	if err != nil {
		return fmt.Errorf("listing network service tags: %+v", err)
	}

	if res.Values == nil {
		return fmt.Errorf("unexpected nil value for service tag information")
	}

	service := d.Get("service").(string)
	locationFilter := azure.NormalizeLocation(d.Get("location_filter"))

	for _, sti := range *res.Values {
		if sti.Name == nil || !isServiceTagOf(*sti.Name, service) {
			continue
		}

		if props := sti.Properties; props != nil {
			if props.Region == nil {
				continue
			}

			if azure.NormalizeLocation(*props.Region) == locationFilter {
				addressPrefixes := make([]string, 0)
				if props.AddressPrefixes != nil {
					addressPrefixes = *props.AddressPrefixes
				}
				err = d.Set("address_prefixes", addressPrefixes)
				if err != nil {
					return fmt.Errorf("setting `address_prefixes`: %+v", err)
				}

				var IPv4 []string
				var IPv6 []string

				for _, prefix := range addressPrefixes {
					ip, ipNet, err := net.ParseCIDR(prefix)
					if err != nil {
						return err
					}

					if ip.To4() != nil {
						IPv4 = append(IPv4, ipNet.String())
					} else {
						IPv6 = append(IPv6, ipNet.String())
					}
				}

				err = d.Set("ipv4_cidrs", IPv4)
				if err != nil {
					return fmt.Errorf("setting `ipv4_cidrs`: %+v", err)
				}

				err = d.Set("ipv6_cidrs", IPv6)
				if err != nil {
					return fmt.Errorf("setting `ipv6_cidrs`: %+v", err)
				}

				if sti.ID == nil {
					return fmt.Errorf("unexcepted nil ID for service tag")
				}

				d.SetId(*sti.ID)
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
