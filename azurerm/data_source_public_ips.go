package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmPublicIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"attached": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"allocation_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Dynamic),
					string(network.Static),
				}, false),
			},

			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPublicIPsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PublicIPsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)

	log.Printf("[DEBUG] Reading Public IP's in Resource Group %q", resourceGroup)
	resp, err := client.List(ctx, resourceGroup)
	if err != nil {
		return fmt.Errorf("Error listing Public IP Addresses in the Resource Group %q: %v", resourceGroup, err)
	}

	filteredIPAddresses := make([]network.PublicIPAddress, 0)
	for _, element := range resp.Values() {
		nicIsAttached := element.IPConfiguration != nil
		shouldInclude := true

		if v, ok := d.GetOkExists("name_prefix"); ok {
			if prefix := v.(string); prefix != "" {
				if !strings.HasPrefix(*element.Name, prefix) {
					shouldInclude = false
				}
			}
		}

		if v, ok := d.GetOkExists("attached"); ok {
			attachedOnly := v.(bool)

			if attachedOnly != nicIsAttached {
				shouldInclude = false
			}
		}

		if v, ok := d.GetOkExists("allocation_type"); ok {
			if allocationType := v.(string); allocationType != "" {
				allocation := network.IPAllocationMethod(allocationType)
				if element.PublicIPAllocationMethod != allocation {
					shouldInclude = false
				}
			}
		}

		if shouldInclude {
			filteredIPAddresses = append(filteredIPAddresses, element)
		}
	}

	d.SetId(time.Now().UTC().String())

	results := flattenDataSourcePublicIPs(filteredIPAddresses)
	if err := d.Set("public_ips", results); err != nil {
		return fmt.Errorf("Error setting `public_ips`: %+v", err)
	}

	return nil
}

func flattenDataSourcePublicIPs(input []network.PublicIPAddress) []interface{} {
	results := make([]interface{}, 0)

	for _, element := range input {
		flattenedIPAddress := flattenDataSourcePublicIP(element)
		results = append(results, flattenedIPAddress)
	}

	return results
}

func flattenDataSourcePublicIP(input network.PublicIPAddress) map[string]string {
	output := make(map[string]string)

	if input.ID != nil {
		output["id"] = *input.ID
	}

	if input.Name != nil {
		output["name"] = *input.Name
	}

	if props := input.PublicIPAddressPropertiesFormat; props != nil {
		if dns := props.DNSSettings; dns != nil {
			if fqdn := dns.Fqdn; fqdn != nil {
				output["fqdn"] = *fqdn
			}

			if label := dns.DomainNameLabel; label != nil {
				output["domain_name_label"] = *label
			}
		}

		if ip := props.IPAddress; ip != nil {
			output["ip_address"] = *ip
		}
	}

	return output
}
