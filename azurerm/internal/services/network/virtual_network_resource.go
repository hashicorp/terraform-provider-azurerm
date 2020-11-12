package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var VirtualNetworkResourceName = "azurerm_virtual_network"

func resourceArmVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkCreateUpdate,
		Read:   resourceArmVirtualNetworkRead,
		Update: resourceArmVirtualNetworkCreateUpdate,
		Delete: resourceArmVirtualNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"bgp_community": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.VirtualNetworkBgpCommunity,
			},

			"ddos_protection_plan": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"enable": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"vm_protection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"address_prefix": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"security_group": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: resourceAzureSubnetHash,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualNetworkCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Network %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_network", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	vnetProperties, vnetPropsErr := expandVirtualNetworkProperties(ctx, d, meta)
	if vnetPropsErr != nil {
		return vnetPropsErr
	}

	vnet := network.VirtualNetwork{
		Name:                           &name,
		Location:                       &location,
		VirtualNetworkPropertiesFormat: vnetProperties,
		Tags:                           tags.Expand(t),
	}

	networkSecurityGroupNames := make([]string, 0)
	for _, subnet := range *vnet.VirtualNetworkPropertiesFormat.Subnets {
		if subnet.NetworkSecurityGroup != nil {
			parsedNsgID, err := azure.ParseAzureResourceID(*subnet.NetworkSecurityGroup.ID)
			if err != nil {
				return fmt.Errorf("Error parsing Network Security Group ID %q: %+v", *subnet.NetworkSecurityGroup.ID, err)
			}

			networkSecurityGroupName := parsedNsgID.Path["networkSecurityGroups"]

			if !utils.SliceContainsValue(networkSecurityGroupNames, networkSecurityGroupName) {
				networkSecurityGroupNames = append(networkSecurityGroupNames, networkSecurityGroupName)
			}
		}
	}

	locks.MultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer locks.UnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	future, err := client.CreateOrUpdate(ctx, resGroup, name, vnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
	client := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VirtualNetworkPropertiesFormat; props != nil {
		d.Set("guid", props.ResourceGUID)

		if space := props.AddressSpace; space != nil {
			d.Set("address_space", utils.FlattenStringSlice(space.AddressPrefixes))
		}

		if err := d.Set("ddos_protection_plan", flattenVirtualNetworkDDoSProtectionPlan(props)); err != nil {
			return fmt.Errorf("Error setting `ddos_protection_plan`: %+v", err)
		}

		if err := d.Set("subnet", flattenVirtualNetworkSubnets(props.Subnets)); err != nil {
			return fmt.Errorf("Error setting `subnets`: %+v", err)
		}

		if err := d.Set("dns_servers", flattenVirtualNetworkDNSServers(props.DhcpOptions)); err != nil {
			return fmt.Errorf("Error setting `dns_servers`: %+v", err)
		}

		bgpCommunity := ""
		if p := props.BgpCommunities; p != nil {
			if v := p.VirtualNetworkCommunity; v != nil {
				bgpCommunity = *v
			}
		}
		if err := d.Set("bgp_community", bgpCommunity); err != nil {
			return fmt.Errorf("Error setting `bgp_community`: %+v", err)
		}

		d.Set("vm_protection_enabled", props.EnableVMProtection)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["virtualNetworks"]

	nsgNames, err := expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d)
	if err != nil {
		return fmt.Errorf("Error parsing Network Security Group ID's: %+v", err)
	}

	locks.MultipleByName(&nsgNames, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&nsgNames, VirtualNetworkResourceName)

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandVirtualNetworkProperties(ctx context.Context, d *schema.ResourceData, meta interface{}) (*network.VirtualNetworkPropertiesFormat, error) {
	subnets := make([]network.Subnet, 0)
	if subs := d.Get("subnet").(*schema.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			// since subnets can also be created outside of vNet definition (as root objects)
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

			// set the props from config and leave the rest intact
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
			AddressPrefixes: utils.ExpandStringSlice(d.Get("address_space").([]interface{})),
		},
		DhcpOptions: &network.DhcpOptions{
			DNSServers: utils.ExpandStringSlice(d.Get("dns_servers").(*schema.Set).List()),
		},
		EnableVMProtection: utils.Bool(d.Get("vm_protection_enabled").(bool)),
		Subnets:            &subnets,
	}

	if v, ok := d.GetOk("ddos_protection_plan"); ok {
		rawList := v.([]interface{})

		var ddosPPlan map[string]interface{}
		if len(rawList) > 0 {
			ddosPPlan = rawList[0].(map[string]interface{})
		}

		if v, ok := ddosPPlan["id"]; ok {
			id := v.(string)
			properties.DdosProtectionPlan = &network.SubResource{
				ID: &id,
			}
		}

		if v, ok := ddosPPlan["enable"]; ok {
			enable := v.(bool)
			properties.EnableDdosProtection = &enable
		}
	}

	if v, ok := d.GetOk("bgp_community"); ok {
		properties.BgpCommunities = &network.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: utils.String(v.(string))}
	}

	return properties, nil
}

func flattenVirtualNetworkDDoSProtectionPlan(input *network.VirtualNetworkPropertiesFormat) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.DdosProtectionPlan == nil || input.DdosProtectionPlan.ID == nil || input.EnableDdosProtection == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id":     *input.DdosProtectionPlan.ID,
			"enable": *input.EnableDdosProtection,
		},
	}
}

func flattenVirtualNetworkSubnets(input *[]network.Subnet) *schema.Set {
	results := &schema.Set{
		F: resourceAzureSubnetHash,
	}

	if subnets := input; subnets != nil {
		for _, subnet := range *input {
			output := map[string]interface{}{}

			if id := subnet.ID; id != nil {
				output["id"] = *id
			}

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
			results = *servers
		}
	}

	return results
}

func resourceAzureSubnetHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
		if v, ok := m["address_prefix"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["security_group"]; ok {
			buf.WriteString(v.(string))
		}
	}

	return hashcode.String(buf.String())
}

func getExistingSubnet(ctx context.Context, resGroup string, vnetName string, subnetName string, meta interface{}) (*network.Subnet, error) {
	subnetClient := meta.(*clients.Client).Network.SubnetsClient
	resp, err := subnetClient.Get(ctx, resGroup, vnetName, subnetName, "")

	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return &network.Subnet{}, nil
		}
		// raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	// Return it directly rather than copy the fields to prevent potential uncovered properties (for example, `ServiceEndpoints` mentioned in #1619)
	return &resp, nil
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
				parsedNsgID, err := azure.ParseAzureResourceID(networkSecurityGroupId)
				if err != nil {
					return nil, fmt.Errorf("Error parsing Network Security Group ID %q: %+v", networkSecurityGroupId, err)
				}

				networkSecurityGroupName := parsedNsgID.Path["networkSecurityGroups"]

				if !utils.SliceContainsValue(nsgNames, networkSecurityGroupName) {
					nsgNames = append(nsgNames, networkSecurityGroupName)
				}
			}
		}
	}

	return nsgNames, nil
}
