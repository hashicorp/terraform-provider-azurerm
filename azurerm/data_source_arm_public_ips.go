package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPublicIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"attached": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"domain_name_label": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPublicIPsRead(d *schema.ResourceData, meta interface{}) error {
	publicIPClient := meta.(*ArmClient).publicIPClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	attachedOnly := d.Get("attached").(bool)
	resp, err := publicIPClient.List(ctx, resGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response().Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure resource group %q: %v", resGroup, err)
	}
	var filteredIps []network.PublicIPAddress
	for _, element := range resp.Values() {
		if (element.IPConfiguration != nil) == attachedOnly {
			filteredIps = append(filteredIps, element)
		}
	}
	var results []map[string]string
	for _, element := range filteredIps {
		m := make(map[string]string)
		if element.ID != nil && *element.ID != "" {
			m["id"] = *element.ID
		}
		if element.Name != nil && *element.Name != "" {
			m["name"] = *element.Name
		}
		if element.PublicIPAddressPropertiesFormat.DNSSettings != nil {
			if element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != nil && *element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != "" {
				m["fqdn"] = *element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn
			}
			if element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != nil && *element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != "" {
				m["domain_name_label"] = *element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel
			}
		}
		if element.PublicIPAddressPropertiesFormat.IPAddress != nil && *element.PublicIPAddressPropertiesFormat.IPAddress != "" {
			m["ip_address"] = *element.PublicIPAddressPropertiesFormat.IPAddress
		}

		results = append(results, m)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("public_ips", results)

	return nil
}
