package azurerm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/Azure/azure-sdk-for-go/arm/network"
)

func dataSourceArmPublicIPs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPublicIPsRead,

		Schema: map[string]*schema.Schema{
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"ids": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"fqdns": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"domain_name_labels": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	var ids, names, fqdns, ip_addresses, domain_name_labels []string
	for _, element := range filteredIps {
		ids = append(ids, *element.ID)
		names = append(names, *element.Name)
		if attachedOnly {
			fqdns = append(fqdns, *element.PublicIPAddressPropertiesFormat.DNSSettings.Fqdn)
			ip_addresses = append(ip_addresses, *element.PublicIPAddressPropertiesFormat.IPAddress)
			domain_name_labels = append(domain_name_labels, *element.PublicIPAddressPropertiesFormat.DNSSettings.DomainNameLabel)
		} else {
			fqdns = append(fqdns, "")
			ip_addresses = append(ip_addresses, "")
			domain_name_labels = append(domain_name_labels, "")
		}
	}

	if minimumCountOk && len(ids) < minimumCount.(int) {
		return fmt.Errorf("Not enough unassigned public IP addresses in resource group %s", resGroup)
	}
	d.SetId(time.Now().UTC().String())
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("fqdns", fqdns)
	d.Set("ip_addresses", ip_addresses)
	d.Set("domain_name_labels", domain_name_labels)

	return nil
}