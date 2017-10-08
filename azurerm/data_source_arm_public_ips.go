package azurerm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceArmPublicIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"minimum_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"attached": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func dataSourceArmPublicIPsRead(d *schema.ResourceData, meta interface{}) error {
	publicIPClient := meta.(*ArmClient).publicIPClient

	resGroup := d.Get("resource_group_name").(string)
	minimumCount, minimumCountOk := d.GetOk("minimum_count")
	attachedOnly := d.Get("attached").(bool)
	resp, err := publicIPClient.List(resGroup)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
		}
		return fmt.Errorf("Error making Read request on Azure resource group %s: %s", resGroup, err)
	}
	var filteredIps []network.PublicIPAddress
	for _, element := range *resp.Value {
		if (element.IPConfiguration != nil) == attachedOnly {
			filteredIps = append(filteredIps, element)
		}
	}
	if minimumCountOk && len(filteredIps) < minimumCount.(int) {
		return fmt.Errorf("Not enough unassigned public IP addresses in resource group %s", resGroup)
	}
	var results []map[string]string
	for _, element := range filteredIps {
		m := make(map[string]string)
		m["public_ip_address_id"] = *element.ID
		m["name"] = *element.Name
		if element.PublicIPAddressPropertiesFormat.DNSSettings != nil {
			if element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != nil && *element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn != "" {
				m["fqdn"] = *element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn
			} else {
				m["fqdn"] = ""
			}
			if element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != nil && *element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel != "" {
				m["domain_name_label"] = *element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel
			} else {
				m["domain_name_label"] = ""
			}
		}
		if element.PublicIPAddressPropertiesFormat.IPAddress != nil && *element.PublicIPAddressPropertiesFormat.IPAddress != "" {
			m["ip_address"] = *element.PublicIPAddressPropertiesFormat.IPAddress
		} else {
			m["ip_address"] = ""
		}

		results = append(results, m)
	}

	d.SetId(time.Now().UTC().String())
	d.Set("public_ips", results)

	return nil
}
