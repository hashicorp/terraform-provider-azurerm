// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/ddosprotectionplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var VirtualNetworkResourceName = "azurerm_virtual_network"

func resourceVirtualNetwork() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkCreate,
		Read:   resourceVirtualNetworkRead,
		Update: resourceVirtualNetworkUpdate,
		Delete: resourceVirtualNetworkDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseVirtualNetworkID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceVirtualNetworkSchema(),
	}
}

func resourceVirtualNetworkSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"address_space": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		// Optional
		"bgp_community": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.VirtualNetworkBgpCommunity,
		},

		"ddos_protection_plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
					},

					"enable": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enforcement": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualnetworks.VirtualNetworkEncryptionEnforcementDropUnencrypted),
							string(virtualnetworks.VirtualNetworkEncryptionEnforcementAllowUnencrypted),
						}, false),
					},
				},
			},
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"flow_timeout_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			// TODO 5.0 Remove Computed and ConfigModeAttr and recommend adding this block to ignore_changes
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefix": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"security_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
			Set: resourceAzureSubnetHash,
		},

		"tags": commonschema.Tags(),
	}

	if !features.FourPointOhBeta() {
		s["address_space"] = &pluginsdk.Schema{
			Type:             pluginsdk.TypeList,
			Required:         true,
			MinItems:         1,
			DiffSuppressFunc: suppress.ListOrder,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		}
	}

	return s
}

func resourceVirtualNetworkCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_network", id.ID())
	}

	vnetProperties, err := expandVirtualNetworkProperties(ctx, *client, id, d)
	if err != nil {
		return err
	}

	vnet := virtualnetworks.VirtualNetwork{
		Name:             pointer.To(id.VirtualNetworkName),
		ExtendedLocation: expandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Properties:       vnetProperties,
		Tags:             tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("flow_timeout_in_minutes"); ok {
		vnet.Properties.FlowTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	networkSecurityGroupNames := make([]string, 0)
	for _, subnet := range *vnet.Properties.Subnets {
		if subnet.Properties != nil && subnet.Properties.NetworkSecurityGroup != nil {
			parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*subnet.Properties.NetworkSecurityGroup.Id)
			if err != nil {
				return err
			}

			networkSecurityGroupName := parsedNsgID.NetworkSecurityGroupName
			if !utils.SliceContainsValue(networkSecurityGroupNames, networkSecurityGroupName) {
				networkSecurityGroupNames = append(networkSecurityGroupNames, networkSecurityGroupName)
			}
		}
	}

	locks.MultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer locks.UnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, id, vnet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkRead(d, meta)
}

func resourceVirtualNetworkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneModel(model.ExtendedLocation))

		if props := model.Properties; props != nil {
			d.Set("guid", props.ResourceGuid)
			d.Set("flow_timeout_in_minutes", props.FlowTimeoutInMinutes)

			if space := props.AddressSpace; space != nil {
				if !features.FourPointOhBeta() {
					d.Set("address_space", utils.FlattenStringSlice(space.AddressPrefixes))
				} else {
					if err = d.Set("address_space", space.AddressPrefixes); err != nil {
						return fmt.Errorf("setting `address_space`: %+v", err)
					}
				}
			}

			if err := d.Set("ddos_protection_plan", flattenVirtualNetworkDDoSProtectionPlan(props)); err != nil {
				return fmt.Errorf("setting `ddos_protection_plan`: %+v", err)
			}

			if err := d.Set("encryption", flattenVirtualNetworkEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}

			if err := d.Set("subnet", flattenVirtualNetworkSubnets(props.Subnets)); err != nil {
				return fmt.Errorf("setting `subnets`: %+v", err)
			}

			if err := d.Set("dns_servers", flattenVirtualNetworkDNSServers(props.DhcpOptions)); err != nil {
				return fmt.Errorf("setting `dns_servers`: %+v", err)
			}

			bgpCommunity := ""
			if p := props.BgpCommunities; p != nil {
				bgpCommunity = p.VirtualNetworkCommunity
			}
			if err := d.Set("bgp_community", bgpCommunity); err != nil {
				return fmt.Errorf("setting `bgp_community`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceVirtualNetworkUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("address_space") {
		if payload.Properties.AddressSpace == nil {
			payload.Properties.AddressSpace = &virtualnetworks.AddressSpace{}
		}
		if !features.FourPointOhBeta() {
			payload.Properties.AddressSpace.AddressPrefixes = utils.ExpandStringSlice(d.Get("address_space").([]interface{}))
		} else {
			payload.Properties.AddressSpace.AddressPrefixes = utils.ExpandStringSlice(d.Get("address_space").(*pluginsdk.Set).List())
		}
	}

	if d.HasChange("bgp_community") {
		// nil out the current values in case `bgp_community` has been removed from the config file
		payload.Properties.BgpCommunities = nil
		if v := d.Get("bgp_community"); v.(string) != "" {
			payload.Properties.BgpCommunities = &virtualnetworks.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: v.(string)}
		}
	}

	if d.HasChange("ddos_protection_plan") {
		ddosProtectionPlanId, enabled := expandVirtualNetworkDdosProtectionPlan(d.Get("ddos_protection_plan").([]interface{}))
		payload.Properties.DdosProtectionPlan = ddosProtectionPlanId
		payload.Properties.EnableDdosProtection = enabled
	}

	if d.HasChange("encryption") {
		// nil out the current values in case `encryption` has been removed from the config file
		payload.Properties.Encryption = expandVirtualNetworkEncryption(d.Get("encryption").([]interface{}))
	}

	if d.HasChange("dns_servers") {
		if payload.Properties.DhcpOptions == nil {
			payload.Properties.DhcpOptions = &virtualnetworks.DhcpOptions{}
		}

		payload.Properties.DhcpOptions.DnsServers = utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
	}

	if d.HasChange("flow_timeout_in_minutes") {
		payload.Properties.FlowTimeoutInMinutes = nil
		if v := d.Get("flow_timeout_in_minutes"); v.(int) != 0 {
			payload.Properties.FlowTimeoutInMinutes = utils.Int64(int64(v.(int)))
		}
	}

	if d.HasChange("subnet") {
		subnets, err := expandVirtualNetworkSubnets(ctx, *client, d.Get("subnet").(*pluginsdk.Set).List(), *id)
		if err != nil {
			return fmt.Errorf("expanding `subnet`: %+v", err)
		}
		payload.Properties.Subnets = subnets
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	networkSecurityGroupNames := make([]string, 0)
	if payload.Properties != nil && payload.Properties.Subnets != nil {
		for _, subnet := range *payload.Properties.Subnets {
			if subnet.Properties != nil && subnet.Properties.NetworkSecurityGroup != nil && subnet.Properties.NetworkSecurityGroup.Id != nil {
				parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*subnet.Properties.NetworkSecurityGroup.Id)
				if err != nil {
					return err
				}

				networkSecurityGroupName := parsedNsgID.NetworkSecurityGroupName
				if !utils.SliceContainsValue(networkSecurityGroupNames, networkSecurityGroupName) {
					networkSecurityGroupNames = append(networkSecurityGroupNames, networkSecurityGroupName)
				}
			}
		}
	}

	locks.MultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)
	defer locks.UnlockMultipleByName(&networkSecurityGroupNames, networkSecurityGroupResourceName)

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, meta.(*clients.Client).Network.VirtualNetworks, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkRead(d, meta)
}

func resourceVirtualNetworkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	nsgNames, err := expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d)
	if err != nil {
		return fmt.Errorf("parsing Network Security Group ID's: %+v", err)
	}

	locks.MultipleByName(&nsgNames, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&nsgNames, VirtualNetworkResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualNetworkDdosProtectionPlan(input []interface{}) (*virtualnetworks.SubResource, *bool) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	var id string
	var enabled bool

	ddosPPlan := input[0].(map[string]interface{})

	if v, ok := ddosPPlan["id"]; ok {
		id = v.(string)
	}

	if v, ok := ddosPPlan["enable"]; ok {
		enabled = v.(bool)
	}

	return &virtualnetworks.SubResource{
		Id: pointer.To(id),
	}, &enabled
}

func expandVirtualNetworkEncryption(input []interface{}) *virtualnetworks.VirtualNetworkEncryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attr := input[0].(map[string]interface{})
	return &virtualnetworks.VirtualNetworkEncryption{
		Enabled:     true,
		Enforcement: pointer.To(virtualnetworks.VirtualNetworkEncryptionEnforcement(attr["enforcement"].(string))),
	}
}

func expandVirtualNetworkSubnets(ctx context.Context, client virtualnetworks.VirtualNetworksClient, input []interface{}, id commonids.VirtualNetworkId) (*[]virtualnetworks.Subnet, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	if len(input) == 0 {
		return &subnets, nil
	}

	for _, subnetRaw := range input {
		if subnetRaw == nil {
			continue
		}
		subnet := subnetRaw.(map[string]interface{})

		name := subnet["name"].(string)
		log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
		// since subnets can also be created outside of vNet definition (as root objects)
		// do a GET on subnet properties from the server before setting them
		subnetObj, err := getExistingSubnet(ctx, client, id, name)
		if err != nil {
			return nil, err
		}
		log.Printf("[INFO] Completed GET of Subnet props ")

		prefix := subnet["address_prefix"].(string)
		secGroup := subnet["security_group"].(string)

		// set the props from config and leave the rest intact
		subnetObj.Name = &name
		if subnetObj.Properties == nil {
			subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
		}

		subnetObj.Properties.AddressPrefix = &prefix

		if secGroup != "" {
			subnetObj.Properties.NetworkSecurityGroup = &virtualnetworks.NetworkSecurityGroup{
				Id: &secGroup,
			}
		} else {
			subnetObj.Properties.NetworkSecurityGroup = nil
		}

		subnets = append(subnets, *subnetObj)
	}

	return &subnets, nil
}

func expandVirtualNetworkProperties(ctx context.Context, client virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId, d *pluginsdk.ResourceData) (*virtualnetworks.VirtualNetworkPropertiesFormat, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	if subs := d.Get("subnet").(*pluginsdk.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			// since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			subnetObj, err := getExistingSubnet(ctx, client, id, name)
			if err != nil {
				return nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props ")

			prefix := subnet["address_prefix"].(string)
			secGroup := subnet["security_group"].(string)

			// set the props from config and leave the rest intact
			subnetObj.Name = &name
			if subnetObj.Properties == nil {
				subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
			}

			subnetObj.Properties.AddressPrefix = &prefix

			if secGroup != "" {
				subnetObj.Properties.NetworkSecurityGroup = &virtualnetworks.NetworkSecurityGroup{
					Id: &secGroup,
				}
			} else {
				subnetObj.Properties.NetworkSecurityGroup = nil
			}

			subnets = append(subnets, *subnetObj)
		}
	}

	properties := &virtualnetworks.VirtualNetworkPropertiesFormat{
		AddressSpace: &virtualnetworks.AddressSpace{},
		DhcpOptions: &virtualnetworks.DhcpOptions{
			DnsServers: utils.ExpandStringSlice(d.Get("dns_servers").([]interface{})),
		},
		Subnets: &subnets,
	}

	if !features.FourPointOhBeta() {
		properties.AddressSpace.AddressPrefixes = utils.ExpandStringSlice(d.Get("address_space").([]interface{}))
	} else {
		properties.AddressSpace.AddressPrefixes = utils.ExpandStringSlice(d.Get("address_space").(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("ddos_protection_plan"); ok {
		rawList := v.([]interface{})

		var ddosPPlan map[string]interface{}
		if len(rawList) > 0 {
			ddosPPlan = rawList[0].(map[string]interface{})
		}

		if v, ok := ddosPPlan["id"]; ok {
			id := v.(string)
			properties.DdosProtectionPlan = &virtualnetworks.SubResource{
				Id: &id,
			}
		}

		if v, ok := ddosPPlan["enable"]; ok {
			enable := v.(bool)
			properties.EnableDdosProtection = &enable
		}
	}

	if v, ok := d.GetOk("encryption"); ok {
		if vList := v.([]interface{}); len(vList) > 0 && vList[0] != nil {
			encryptionConf := vList[0].(map[string]interface{})
			properties.Encryption = &virtualnetworks.VirtualNetworkEncryption{
				Enabled:     true,
				Enforcement: pointer.To(virtualnetworks.VirtualNetworkEncryptionEnforcement(encryptionConf["enforcement"].(string))),
			}
		}
	}

	if v, ok := d.GetOk("bgp_community"); ok {
		properties.BgpCommunities = &virtualnetworks.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: v.(string)}
	}

	return properties, nil
}

func flattenVirtualNetworkDDoSProtectionPlan(input *virtualnetworks.VirtualNetworkPropertiesFormat) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.DdosProtectionPlan == nil || input.DdosProtectionPlan.Id == nil || input.EnableDdosProtection == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id":     *input.DdosProtectionPlan.Id,
			"enable": *input.EnableDdosProtection,
		},
	}
}

func flattenVirtualNetworkEncryption(encryption *virtualnetworks.VirtualNetworkEncryption) interface{} {
	if encryption == nil || !encryption.Enabled {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"enforcement": encryption.Enforcement,
		},
	}
}

func flattenVirtualNetworkSubnets(input *[]virtualnetworks.Subnet) *pluginsdk.Set {
	results := &pluginsdk.Set{
		F: resourceAzureSubnetHash,
	}

	if subnets := input; subnets != nil {
		for _, subnet := range *input {
			output := map[string]interface{}{}

			if id := subnet.Id; id != nil {
				output["id"] = *id
			}

			if name := subnet.Name; name != nil {
				output["name"] = *name
			}

			if props := subnet.Properties; props != nil {
				if prefix := props.AddressPrefix; prefix != nil {
					output["address_prefix"] = *prefix
				}

				if nsg := props.NetworkSecurityGroup; nsg != nil {
					if nsg.Id != nil {
						output["security_group"] = *nsg.Id
					}
				}
			}

			results.Add(output)
		}
	}

	return results
}

func flattenVirtualNetworkDNSServers(input *virtualnetworks.DhcpOptions) []string {
	results := make([]string, 0)

	if input != nil {
		if servers := input.DnsServers; servers != nil {
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

	return pluginsdk.HashString(buf.String())
}

func getExistingSubnet(ctx context.Context, client virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId, subnetName string) (*virtualnetworks.Subnet, error) {
	resp, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		// The Subnet doesn't exist when the Virtual Network doesn't exist
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(virtualnetworks.Subnet{}), nil
		}
		// raise an error if there was an issue other than 404 in getting subnet properties
		return nil, err
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			for _, subnet := range pointer.From(props.Subnets) {
				if subnetName == pointer.From(subnet.Name) {
					// Return it directly rather than copy the fields to prevent potential uncovered properties (for example, `ServiceEndpoints` mentioned in #1619)
					return pointer.To(subnet), nil
				}
			}
		}
	}

	// TODO 4.0: Return empty object when the Subnet isn't found
	return pointer.To(virtualnetworks.Subnet{}), nil
}

func expandAzureRmVirtualNetworkVirtualNetworkSecurityGroupNames(d *pluginsdk.ResourceData) ([]string, error) {
	nsgNames := make([]string, 0)

	if v, ok := d.GetOk("subnet"); ok {
		subnets := v.(*pluginsdk.Set).List()
		for _, subnet := range subnets {
			subnet, ok := subnet.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("[ERROR] Subnet should be a Hash - was '%+v'", subnet)
			}

			networkSecurityGroupId := subnet["security_group"].(string)
			if networkSecurityGroupId != "" {
				parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(networkSecurityGroupId)
				if err != nil {
					return nil, err
				}

				networkSecurityGroupName := parsedNsgID.NetworkSecurityGroupName
				if !utils.SliceContainsValue(nsgNames, networkSecurityGroupName) {
					nsgNames = append(nsgNames, networkSecurityGroupName)
				}
			}
		}
	}

	return nsgNames, nil
}

func VirtualNetworkProvisioningStateRefreshFunc(ctx context.Context, client *virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties != nil {
			return res, string(pointer.From(res.Model.Properties.ProvisioningState)), nil
		}
		return res, "", fmt.Errorf("polling for %s: %+v", id, err)
	}
}
