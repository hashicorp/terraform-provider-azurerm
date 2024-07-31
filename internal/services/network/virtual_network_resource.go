// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/ddosprotectionplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
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
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefixes": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"default_outbound_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},

					"delegation": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						MaxItems:   1,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"service_delegation": {
									Type:       pluginsdk.TypeList,
									Required:   true,
									MaxItems:   1,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(subnetDelegationServiceNames, false),
											},

											"actions": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														"Microsoft.Network/networkinterfaces/*",
														"Microsoft.Network/publicIPAddresses/join/action",
														"Microsoft.Network/publicIPAddresses/read",
														"Microsoft.Network/virtualNetworks/read",
														"Microsoft.Network/virtualNetworks/subnets/action",
														"Microsoft.Network/virtualNetworks/subnets/join/action",
														"Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
														"Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},

					"private_endpoint_network_policies": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
						ValidateFunc: validation.StringInSlice(subnets.PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies(), false),
					},

					"private_link_service_network_policies_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"route_table_id": commonschema.ResourceIDReferenceOptional(&routetables.RouteTableId{}),

					"security_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"service_endpoints": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"service_endpoint_policy_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
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
		s["subnet"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
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

	vnetProperties, routeTables, err := expandVirtualNetworkProperties(ctx, *client, id, d)
	if err != nil {
		return err
	}

	locks.MultipleByName(routeTables, routeTableResourceName)
	defer locks.UnlockMultipleByName(routeTables, routeTableResourceName)

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

			subnet, err := flattenVirtualNetworkSubnets(props.Subnets)
			if err != nil {
				return fmt.Errorf("flattening `subnet`: %+v", err)
			}
			if err := d.Set("subnet", subnet); err != nil {
				return fmt.Errorf("setting `subnet`: %+v", err)
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
		subnets, routeTables, err := expandVirtualNetworkSubnets(ctx, *client, d.Get("subnet").(*pluginsdk.Set).List(), *id)
		if err != nil {
			return fmt.Errorf("expanding `subnet`: %+v", err)
		}
		payload.Properties.Subnets = subnets

		locks.MultipleByName(routeTables, routeTableResourceName)
		defer locks.UnlockMultipleByName(routeTables, routeTableResourceName)
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

	nsgNames, routeTableNames, err := expandResourcesForLocking(d)
	if err != nil {
		return fmt.Errorf("parsing Network Security Group ID's: %+v", err)
	}

	locks.MultipleByName(&nsgNames, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&nsgNames, VirtualNetworkResourceName)

	locks.MultipleByName(&routeTableNames, routeTableResourceName)
	defer locks.UnlockMultipleByName(&routeTableNames, routeTableResourceName)

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

func expandVirtualNetworkSubnets(ctx context.Context, client virtualnetworks.VirtualNetworksClient, input []interface{}, id commonids.VirtualNetworkId) (*[]virtualnetworks.Subnet, *[]string, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	routeTables := make([]string, 0)

	if len(input) == 0 {
		return &subnets, &routeTables, nil
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
			return nil, nil, err
		}
		log.Printf("[INFO] Completed GET of Subnet props")

		// set the props from config and leave the rest intact
		subnetObj.Name = pointer.To(name)
		if subnetObj.Properties == nil {
			subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
		}

		if !features.FourPointOhBeta() {
			subnetObj.Properties.AddressPrefix = pointer.To(subnet["address_prefix"].(string))
		}

		if features.FourPointOhBeta() {
			addressPrefixes := make([]string, 0)
			for _, prefix := range subnet["address_prefixes"].([]interface{}) {
				addressPrefixes = append(addressPrefixes, prefix.(string))
			}

			if len(addressPrefixes) == 1 {
				subnetObj.Properties.AddressPrefix = pointer.To(addressPrefixes[0])
				subnetObj.Properties.AddressPrefixes = nil
			} else {
				subnetObj.Properties.AddressPrefixes = pointer.To(addressPrefixes)
				subnetObj.Properties.AddressPrefix = nil
			}

			privateEndpointNetworkPolicies := virtualnetworks.VirtualNetworkPrivateEndpointNetworkPolicies(subnet["private_endpoint_network_policies"].(string))
			privateLinkServiceNetworkPolicies := virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled
			if subnet["private_link_service_network_policies_enabled"].(bool) {
				privateLinkServiceNetworkPolicies = virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled
			}
			subnetObj.Properties.DefaultOutboundAccess = pointer.To(subnet["default_outbound_access_enabled"].(bool))
			subnetObj.Properties.Delegations = expandVirtualNetworkSubnetDelegation(subnet["delegation"].([]interface{}))
			subnetObj.Properties.PrivateEndpointNetworkPolicies = pointer.To(privateEndpointNetworkPolicies)
			subnetObj.Properties.PrivateLinkServiceNetworkPolicies = pointer.To(privateLinkServiceNetworkPolicies)

			if routeTableId := subnet["route_table_id"].(string); routeTableId != "" {
				id, err := routetables.ParseRouteTableID(routeTableId)
				if err != nil {
					return nil, nil, err
				}

				// Collecting a list of route tables to lock on outside of this function
				routeTables = append(routeTables, id.RouteTableName)
				subnetObj.Properties.RouteTable = &virtualnetworks.RouteTable{
					Id: pointer.To(id.ID()),
				}
			} else {
				subnetObj.Properties.RouteTable = nil
			}

			subnetObj.Properties.ServiceEndpointPolicies = expandVirtualNetworkSubnetServiceEndpointPolicies(subnet["service_endpoint_policy_ids"].(*pluginsdk.Set).List())
			subnetObj.Properties.ServiceEndpoints = expandVirtualNetworkSubnetServiceEndpoints(subnet["service_endpoints"].(*pluginsdk.Set).List())
		}

		if secGroup := subnet["security_group"].(string); secGroup != "" {
			subnetObj.Properties.NetworkSecurityGroup = &virtualnetworks.NetworkSecurityGroup{
				Id: &secGroup,
			}
		} else {
			subnetObj.Properties.NetworkSecurityGroup = nil
		}

		subnets = append(subnets, *subnetObj)
	}

	return &subnets, &routeTables, nil
}

func expandVirtualNetworkProperties(ctx context.Context, client virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId, d *pluginsdk.ResourceData) (*virtualnetworks.VirtualNetworkPropertiesFormat, *[]string, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	routeTables := make([]string, 0)
	if subs := d.Get("subnet").(*pluginsdk.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]interface{})

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			// since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			subnetObj, err := getExistingSubnet(ctx, client, id, name)
			if err != nil {
				return nil, nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props")

			// set the props from config and leave the rest intact
			subnetObj.Name = pointer.To(name)
			if subnetObj.Properties == nil {
				subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
			}

			if !features.FourPointOhBeta() {
				subnetObj.Properties.AddressPrefix = pointer.To(subnet["address_prefix"].(string))
			}

			if features.FourPointOhBeta() {
				addressPrefixes := make([]string, 0)
				for _, prefix := range subnet["address_prefixes"].([]interface{}) {
					addressPrefixes = append(addressPrefixes, prefix.(string))
				}

				if len(addressPrefixes) == 1 {
					subnetObj.Properties.AddressPrefix = pointer.To(addressPrefixes[0])
				} else {
					subnetObj.Properties.AddressPrefixes = pointer.To(addressPrefixes)
				}

				privateEndpointNetworkPolicies := virtualnetworks.VirtualNetworkPrivateEndpointNetworkPolicies(subnet["private_endpoint_network_policies"].(string))
				privateLinkServiceNetworkPolicies := virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled
				if subnet["private_link_service_network_policies_enabled"].(bool) {
					privateLinkServiceNetworkPolicies = virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled
				}
				subnetObj.Properties.DefaultOutboundAccess = pointer.To(subnet["default_outbound_access_enabled"].(bool))
				subnetObj.Properties.Delegations = expandVirtualNetworkSubnetDelegation(subnet["delegation"].([]interface{}))
				subnetObj.Properties.PrivateEndpointNetworkPolicies = pointer.To(privateEndpointNetworkPolicies)
				subnetObj.Properties.PrivateLinkServiceNetworkPolicies = pointer.To(privateLinkServiceNetworkPolicies)

				if routeTableId := subnet["route_table_id"].(string); routeTableId != "" {
					id, err := routetables.ParseRouteTableID(routeTableId)
					if err != nil {
						return nil, nil, err
					}

					// Collecting a list of route tables to lock on outside of this function
					routeTables = append(routeTables, id.RouteTableName)
					subnetObj.Properties.RouteTable = &virtualnetworks.RouteTable{
						Id: pointer.To(id.ID()),
					}
				}

				subnetObj.Properties.ServiceEndpointPolicies = expandVirtualNetworkSubnetServiceEndpointPolicies(subnet["service_endpoint_policy_ids"].(*pluginsdk.Set).List())
				subnetObj.Properties.ServiceEndpoints = expandVirtualNetworkSubnetServiceEndpoints(subnet["service_endpoints"].(*pluginsdk.Set).List())
			}

			if secGroup := subnet["security_group"].(string); secGroup != "" {
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

	return properties, &routeTables, nil
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

func flattenVirtualNetworkSubnets(input *[]virtualnetworks.Subnet) (*pluginsdk.Set, error) {
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
				if !features.FourPointOhBeta() {
					if prefix := props.AddressPrefix; prefix != nil {
						output["address_prefix"] = *prefix
					}
				}

				if nsg := props.NetworkSecurityGroup; nsg != nil {
					if nsg.Id != nil {
						output["security_group"] = *nsg.Id
					}
				}

				if features.FourPointOhBeta() {
					if props.AddressPrefixes == nil {
						if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
							output["address_prefixes"] = []string{*props.AddressPrefix}
						} else {
							output["address_prefixes"] = []string{}
						}
					} else {
						output["address_prefixes"] = props.AddressPrefixes
					}
					output["delegation"] = flattenVirtualNetworkSubnetDelegation(props.Delegations)
					output["default_outbound_access_enabled"] = pointer.From(props.DefaultOutboundAccess)
					output["private_endpoint_network_policies"] = string(pointer.From(props.PrivateEndpointNetworkPolicies))
					output["private_link_service_network_policies_enabled"] = strings.EqualFold(string(pointer.From(props.PrivateLinkServiceNetworkPolicies)), string(virtualnetworks.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled))
					routeTableId := ""
					if props.RouteTable != nil && props.RouteTable.Id != nil {
						id, err := routetables.ParseRouteTableID(*props.RouteTable.Id)
						if err != nil {
							return nil, err
						}
						routeTableId = id.ID()
					}
					output["route_table_id"] = routeTableId
					output["service_endpoints"] = flattenVirtualNetworkSubnetServiceEndpoints(props.ServiceEndpoints)
					output["service_endpoint_policy_ids"] = flattenVirtualNetworkSubnetServiceEndpointPolicies(props.ServiceEndpointPolicies)
				}
			}

			results.Add(output)
		}
	}

	return results, nil
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

	return pointer.To(virtualnetworks.Subnet{}), nil
}

func expandResourcesForLocking(d *pluginsdk.ResourceData) ([]string, []string, error) {
	nsgNames := make([]string, 0)
	routeTableNames := make([]string, 0)

	if v, ok := d.GetOk("subnet"); ok {
		subnets := v.(*pluginsdk.Set).List()
		for _, subnet := range subnets {
			subnet, ok := subnet.(map[string]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("[ERROR] Subnet should be a Hash - was '%+v'", subnet)
			}

			networkSecurityGroupId := subnet["security_group"].(string)
			if networkSecurityGroupId != "" {
				parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(networkSecurityGroupId)
				if err != nil {
					return nil, nil, err
				}

				networkSecurityGroupName := parsedNsgID.NetworkSecurityGroupName
				if !utils.SliceContainsValue(nsgNames, networkSecurityGroupName) {
					nsgNames = append(nsgNames, networkSecurityGroupName)
				}
			}

			if features.FourPointOhBeta() {
				routeTableId := subnet["route_table_id"].(string)
				if routeTableId != "" {
					parsedRouteTableID, err := routetables.ParseRouteTableID(routeTableId)
					if err != nil {
						return nil, nil, err
					}
					routeTableName := parsedRouteTableID.RouteTableName
					if !utils.SliceContainsValue(routeTableNames, routeTableName) {
						routeTableNames = append(routeTableNames, routeTableName)
					}
				}
			}
		}
	}

	return nsgNames, routeTableNames, nil
}

func expandVirtualNetworkSubnetServiceEndpointPolicies(input []interface{}) *[]virtualnetworks.ServiceEndpointPolicy {
	output := make([]virtualnetworks.ServiceEndpointPolicy, 0)
	for _, policy := range input {
		policy := policy.(string)
		output = append(output, virtualnetworks.ServiceEndpointPolicy{Id: &policy})
	}
	return &output
}

func expandVirtualNetworkSubnetServiceEndpoints(input []interface{}) *[]virtualnetworks.ServiceEndpointPropertiesFormat {
	endpoints := make([]virtualnetworks.ServiceEndpointPropertiesFormat, 0)

	for _, svcEndpointRaw := range input {
		if svc, ok := svcEndpointRaw.(string); ok {
			endpoint := virtualnetworks.ServiceEndpointPropertiesFormat{
				Service: &svc,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return &endpoints
}

func expandVirtualNetworkSubnetDelegation(input []interface{}) *[]virtualnetworks.Delegation {
	retDelegations := make([]virtualnetworks.Delegation, 0)

	for _, deleValue := range input {
		deleData := deleValue.(map[string]interface{})
		deleName := deleData["name"].(string)
		srvDelegations := deleData["service_delegation"].([]interface{})
		srvDelegation := srvDelegations[0].(map[string]interface{})
		srvName := srvDelegation["name"].(string)

		srvActions := srvDelegation["actions"].(*pluginsdk.Set).List()

		retSrvActions := make([]string, 0)
		for _, srvAction := range srvActions {
			srvActionData := srvAction.(string)
			retSrvActions = append(retSrvActions, srvActionData)
		}

		retDelegation := virtualnetworks.Delegation{
			Name: &deleName,
			Properties: &virtualnetworks.ServiceDelegationPropertiesFormat{
				ServiceName: &srvName,
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return &retDelegations
}

func flattenVirtualNetworkSubnetServiceEndpointPolicies(input *[]virtualnetworks.ServiceEndpointPolicy) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, policy := range *input {
		id := ""
		if policy.Id != nil {
			id = *policy.Id
		}
		output = append(output, id)
	}
	return output
}

func flattenVirtualNetworkSubnetServiceEndpoints(serviceEndpoints *[]virtualnetworks.ServiceEndpointPropertiesFormat) []interface{} {
	endpoints := make([]interface{}, 0)

	if serviceEndpoints == nil {
		return endpoints
	}

	for _, endpoint := range *serviceEndpoints {
		if endpoint.Service != nil {
			endpoints = append(endpoints, *endpoint.Service)
		}
	}

	return endpoints
}

func flattenVirtualNetworkSubnetDelegation(delegations *[]virtualnetworks.Delegation) []interface{} {
	if delegations == nil {
		return []interface{}{}
	}

	retDeles := make([]interface{}, 0)

	normalizeServiceName := map[string]string{}
	for _, normName := range subnetDelegationServiceNames {
		normalizeServiceName[strings.ToLower(normName)] = normName
	}

	for _, dele := range *delegations {
		retDele := make(map[string]interface{})
		if v := dele.Name; v != nil {
			retDele["name"] = *v
		}

		svcDeles := make([]interface{}, 0)
		svcDele := make(map[string]interface{})
		if props := dele.Properties; props != nil {
			if v := props.ServiceName; v != nil {
				name := *v
				if nv, ok := normalizeServiceName[strings.ToLower(name)]; ok {
					name = nv
				}
				svcDele["name"] = name
			}

			if v := props.Actions; v != nil {
				svcDele["actions"] = *v
			}
		}

		svcDeles = append(svcDeles, svcDele)

		retDele["service_delegation"] = svcDeles

		retDeles = append(retDeles, retDele)
	}

	return retDeles
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
