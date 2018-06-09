package azurerm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var virtualNetworkResourceName = "azurerm_virtual_network"

func resourceArmVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkCreate,
		Read:   resourceArmVirtualNetworkRead,
		Update: resourceArmVirtualNetworkCreate,
		Delete: resourceArmVirtualNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"subnet": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address_prefix": {
							Type:     schema.TypeString,
							Required: true,
						},
						"security_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceAzureSubnetHash,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	vnetProperties, vnetPropsErr := expandVirtualNetworkProperties(ctx, d, meta)
	if vnetPropsErr != nil {
		return vnetPropsErr
	}

	vnet := network.VirtualNetwork{
		Name:                           &name,
		Location:                       &location,
		VirtualNetworkPropertiesFormat: vnetProperties,
		Tags: expandTags(tags),
	}

	networkSecurityGroupNames := make([]string, 0)
	for _, subnet := range *vnet.VirtualNetworkPropertiesFormat.Subnets {
		if subnet.NetworkSecurityGroup != nil {
			nsgName, err := parseNetworkSecurityGroupName(*subnet.NetworkSecurityGroup.ID)
			if err != nil {
				return err
			}

			if !sliceContainsValue(networkSecurityGroupNames, nsgName) {
				networkSecurityGroupNames = append(networkSecurityGroupNames, nsgName)
			}
		}
	}

	azureRMLockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer azureRMUnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	future, err := client.CreateOrUpdate(ctx, resGroup, name, vnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Network %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkRead(d, meta)
}

func resourceArmVirtualNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualNetworks"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		if space := props.AddressSpace; space != nil {
			d.Set("address_space", space.AddressPrefixes)
		}

		subnets := flattenVirtualNetworkSubnets(props.Subnets)
		if err := d.Set("subnet", subnets); err != nil {
			return fmt.Errorf("Error setting `subnets`: %+v", err)
		}

		dnsServers := flattenVirtualNetworkDNSServers(props.DhcpOptions)
		if err := d.Set("dns_servers", dnsServers); err != nil {
			return fmt.Errorf("Error setting `dns_servers`: %+v", err)
		}

	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualNetworks"]

	nsgNames, err := expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Error parsing Network Security Group ID's: %+v", err)
	}

	azureRMLockMultipleByName(&nsgNames, virtualNetworkResourceName)
	defer azureRMUnlockMultipleByName(&nsgNames, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandVirtualNetworkProperties(ctx context.Context, d *schema.ResourceData, meta interface{}) (*network.VirtualNetworkPropertiesFormat, error) {
	// first; get address space prefixes:
	prefixes := []string{}
	for _, prefix := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, prefix.(string))
	}

	// then; the dns servers:
	dnses := []string{}
	for _, dns := range d.Get("dns_servers").([]interface{}) {
		dnses = append(dnses, dns.(string))
	}

	// then; the subnets:
	subnets := []network.Subnet{}
	if subs := d.Get("subnet").(*schema.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			//since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			resGroup := d.Get("resource_group_name").(string)
			vnetName := d.Get("name").(string)
			subnetObj, err := getExistingSubnet(ctx, resGroup, vnetName, name, meta)
			if err != nil {
				return nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props ")

			prefix := subnet["address_prefix"].(string)
			secGroup := subnet["security_group"].(string)

			//set the props from config and leave the rest intact
			subnetObj.Name = &name
			if subnetObj.SubnetPropertiesFormat == nil {
				subnetObj.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{}
			}

			subnetObj.SubnetPropertiesFormat.AddressPrefix = &prefix

			if secGroup != "" {
				subnetObj.SubnetPropertiesFormat.NetworkSecurityGroup = &network.SecurityGroup{
					ID: &secGroup,
				}
			} else {
				subnetObj.SubnetPropertiesFormat.NetworkSecurityGroup = nil
			}

			subnets = append(subnets, *subnetObj)
		}
	}

	properties := &network.VirtualNetworkPropertiesFormat{
		AddressSpace: &network.AddressSpace{
			AddressPrefixes: &prefixes,
		},
		DhcpOptions: &network.DhcpOptions{
			DNSServers: &dnses,
		},
		Subnets: &subnets,
	}
	// finally; return the struct:
	return properties, nil
}
func flattenVirtualNetworkSubnets(input *[]network.Subnet) *schema.Set {
	results := &schema.Set{
		F: resourceAzureSubnetHash,
	}

	if subnets := input; subnets != nil {
		for _, subnet := range *input {
			output := map[string]interface{}{}

			if name := subnet.Name; name != nil {
				output["name"] = *name
			}

			if props := subnet.SubnetPropertiesFormat; props != nil {
				if prefix := props.AddressPrefix; prefix != nil {
					output["address_prefix"] = *prefix
				}

				if nsg := props.NetworkSecurityGroup; nsg != nil {
					if nsg.ID != nil {
						output["security_group"] = *nsg.ID
					}
				}
			}

			results.Add(output)
		}
	}

	return results
}

func flattenVirtualNetworkDNSServers(input *network.DhcpOptions) []string {
	results := make([]string, 0)

	if input != nil {
		if servers := input.DNSServers; servers != nil {
			for _, dns := range *servers {
				results = append(results, dns)
			}
		}
	}

	return results
}

func resourceAzureSubnetHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(fmt.Sprintf("%s", m["name"].(string)))
		buf.WriteString(fmt.Sprintf("%s", m["address_prefix"].(string)))

		if v, ok := m["security_group"]; ok {
			buf.WriteString(v.(string))
		}
	}

	return hashcode.String(buf.String())
}

func getExistingSubnet(ctx context.Context, resGroup string, vnetName string, subnetName string, meta interface{}) (*network.Subnet, error) {
	//attempt to retrieve existing subnet from the server
	existingSubnet := network.Subnet{}
	subnetClient := meta.(*ArmClient).subnetClient
	resp, err := subnetClient.Get(ctx, resGroup, vnetName, subnetName, "")

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return &existingSubnet, nil
		}
		//raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	existingSubnet.SubnetPropertiesFormat = &network.SubnetPropertiesFormat{
		AddressPrefix: resp.SubnetPropertiesFormat.AddressPrefix,
	}

	if resp.SubnetPropertiesFormat.NetworkSecurityGroup != nil {
		existingSubnet.SubnetPropertiesFormat.NetworkSecurityGroup = resp.SubnetPropertiesFormat.NetworkSecurityGroup
	}

	if resp.SubnetPropertiesFormat.RouteTable != nil {
		existingSubnet.SubnetPropertiesFormat.RouteTable = resp.SubnetPropertiesFormat.RouteTable
	}

	if resp.SubnetPropertiesFormat.IPConfigurations != nil {
		ips := make([]string, 0)
		for _, ip := range *resp.SubnetPropertiesFormat.IPConfigurations {
			ips = append(ips, *ip.ID)
		}

		existingSubnet.SubnetPropertiesFormat.IPConfigurations = resp.SubnetPropertiesFormat.IPConfigurations
	}

	return &existingSubnet, nil
}

func expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d *schema.ResourceData) ([]string, error) {
	nsgNames := make([]string, 0)

	if v, ok := d.GetOk("subnet"); ok {
		subnets := v.(*schema.Set).List()
		for _, subnet := range subnets {
			subnet, ok := subnet.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR] Subnet should be a Hash - was '%+v'", subnet)
			}

			networkSecurityGroupId := subnet["security_group"].(string)
			if networkSecurityGroupId != "" {
				nsgName, err := parseNetworkSecurityGroupName(networkSecurityGroupId)
				if err != nil {
					return nil, err
				}

				if !sliceContainsValue(nsgNames, nsgName) {
					nsgNames = append(nsgNames, nsgName)
				}
			}
		}
	}

	return nsgNames, nil
}
