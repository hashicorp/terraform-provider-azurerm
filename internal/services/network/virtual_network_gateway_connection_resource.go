// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgatewayconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualNetworkGatewayConnection() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayConnectionCreate,
		Read:   resourceVirtualNetworkGatewayConnectionRead,
		Update: resourceVirtualNetworkGatewayConnectionUpdate,
		Delete: resourceVirtualNetworkGatewayConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworkgatewayconnections.ParseConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute),
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec),
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet),
				}, false),
			},

			"virtual_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
			},

			"shared_key": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C the API generates a key for the user if not supplied
				Computed:  true,
				Sensitive: true,
			},

			"authorization_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"dpd_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"express_route_circuit_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: expressroutecircuits.ValidateExpressRouteCircuitID,
			},

			"egress_nat_rule_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayNatRuleID,
				},
			},

			"ingress_nat_rule_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayNatRuleID,
				},
			},

			"peer_virtual_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: virtualnetworkgateways.ValidateVirtualNetworkGatewayID,
			},

			"local_azure_ip_address_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"local_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: localnetworkgateways.ValidateLocalNetworkGatewayID,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_bgp": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"use_policy_based_traffic_selectors": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"routing_weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 32000),
			},

			"express_route_gateway_bypass": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"private_link_fast_path_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"connection_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionProtocolIKEvOne),
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionProtocolIKEvTwo),
				}, false),
			},

			"connection_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeInitiatorOnly),
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeResponderOnly),
					string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeDefault),
				}, false),
				Default: string(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionModeDefault),
			},

			"traffic_selector_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"remote_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"custom_bgp_addresses": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"primary": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.Any(validation.IsIPv4Address),
						},
						"secondary": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.Any(validation.IsIPv4Address),
						},
					},
				},
			},

			"ipsec_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dh_group": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.DhGroupDHGroupOne),
								string(virtualnetworkgatewayconnections.DhGroupDHGroupOneFour),
								string(virtualnetworkgatewayconnections.DhGroupDHGroupTwo),
								string(virtualnetworkgatewayconnections.DhGroupDHGroupTwoZeroFourEight),
								string(virtualnetworkgatewayconnections.DhGroupDHGroupTwoFour),
								string(virtualnetworkgatewayconnections.DhGroupECPTwoFiveSix),
								string(virtualnetworkgatewayconnections.DhGroupECPThreeEightFour),
								string(virtualnetworkgatewayconnections.DhGroupNone),
							}, false),
						},

						"ike_encryption": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.IkeEncryptionAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IkeEncryptionAESOneNineTwo),
								string(virtualnetworkgatewayconnections.IkeEncryptionAESTwoFiveSix),
								string(virtualnetworkgatewayconnections.IkeEncryptionDES),
								string(virtualnetworkgatewayconnections.IkeEncryptionDESThree),
								string(virtualnetworkgatewayconnections.IkeEncryptionGCMAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IkeEncryptionGCMAESTwoFiveSix),
							}, false),
						},

						"ike_integrity": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.IkeIntegrityGCMAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IkeIntegrityGCMAESTwoFiveSix),
								string(virtualnetworkgatewayconnections.IkeIntegrityMDFive),
								string(virtualnetworkgatewayconnections.IkeIntegritySHAOne),
								string(virtualnetworkgatewayconnections.IkeIntegritySHATwoFiveSix),
								string(virtualnetworkgatewayconnections.IkeIntegritySHAThreeEightFour),
							}, false),
						},

						"ipsec_encryption": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.IPsecEncryptionAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IPsecEncryptionAESOneNineTwo),
								string(virtualnetworkgatewayconnections.IPsecEncryptionAESTwoFiveSix),
								string(virtualnetworkgatewayconnections.IPsecEncryptionDES),
								string(virtualnetworkgatewayconnections.IPsecEncryptionDESThree),
								string(virtualnetworkgatewayconnections.IPsecEncryptionGCMAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IPsecEncryptionGCMAESOneNineTwo),
								string(virtualnetworkgatewayconnections.IPsecEncryptionGCMAESTwoFiveSix),
								string(virtualnetworkgatewayconnections.IPsecEncryptionNone),
							}, false),
						},

						"ipsec_integrity": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.IPsecIntegrityGCMAESOneTwoEight),
								string(virtualnetworkgatewayconnections.IPsecIntegrityGCMAESOneNineTwo),
								string(virtualnetworkgatewayconnections.IPsecIntegrityGCMAESTwoFiveSix),
								string(virtualnetworkgatewayconnections.IPsecIntegrityMDFive),
								string(virtualnetworkgatewayconnections.IPsecIntegritySHAOne),
								string(virtualnetworkgatewayconnections.IPsecIntegritySHATwoFiveSix),
							}, false),
						},

						"pfs_group": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgatewayconnections.PfsGroupECPTwoFiveSix),
								string(virtualnetworkgatewayconnections.PfsGroupECPThreeEightFour),
								string(virtualnetworkgatewayconnections.PfsGroupNone),
								string(virtualnetworkgatewayconnections.PfsGroupPFSOne),
								string(virtualnetworkgatewayconnections.PfsGroupPFSOneFour),
								string(virtualnetworkgatewayconnections.PfsGroupPFSTwo),
								string(virtualnetworkgatewayconnections.PfsGroupPFSTwoZeroFourEight),
								string(virtualnetworkgatewayconnections.PfsGroupPFSTwoFour),
								string(virtualnetworkgatewayconnections.PfsGroupPFSMM),
							}, false),
						},

						"sa_datasize": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 2147483647),
						},

						"sa_lifetime": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(300),
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["shared_key"] = &pluginsdk.Schema{
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Optional:  true,
			Sensitive: true,
		}
	}
	return resource
}

func resourceVirtualNetworkGatewayConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGatewayConnections
	vnetGatewayClient := meta.(*clients.Client).Network.VirtualNetworkGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway Connection creation.")

	id := virtualnetworkgatewayconnections.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_network_gateway_connection", id.ID())
	}

	var virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway
	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return err
		}

		resp, err := vnetGatewayClient.Get(ctx, *gwid)
		if err != nil {
			return err
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", gwid, err)
		}

		virtualNetworkGateway = *resp.Model
	}

	properties, err := getVirtualNetworkGatewayConnectionProperties(d, virtualNetworkGateway)
	if err != nil {
		return err
	}

	connection := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnection{
		Name:       pointer.To(id.ConnectionName),
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: *properties,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, connection); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if properties.SharedKey != nil && !d.IsNewResource() {
		if err := client.SetSharedKeyThenPoll(ctx, id, virtualnetworkgatewayconnections.ConnectionSharedKey{
			Value: pointer.From(properties.SharedKey),
		}); err != nil {
			return fmt.Errorf("updating Shared Key for %s: %+v", id, err)
		}

		// Once this issue https://github.com/Azure/azure-rest-api-specs/issues/26660 is fixed, below this part will be removed
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(virtualnetworkgatewayconnections.ProvisioningStateUpdating)},
			Target:     []string{string(virtualnetworkgatewayconnections.ProvisioningStateSucceeded)},
			Refresh:    virtualNetworkGatewayConnectionStateRefreshFunc(ctx, client, id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGatewayConnections
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgatewayconnections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("resource_group_name", id.ResourceGroupName)

	respKey, err := client.GetSharedKey(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving Shared Key for %s: %+v", id, err)
	}

	if model := respKey.Model; model != nil {
		if model.Value != "" {
			d.Set("shared_key", model.Value)
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties

		if string(props.ConnectionType) != "" {
			d.Set("type", string(props.ConnectionType))
		}

		d.Set("virtual_network_gateway_id", props.VirtualNetworkGateway1.Id)

		// not returned from api - getting from state
		d.Set("authorization_key", d.Get("authorization_key").(string))

		if props.DpdTimeoutSeconds != nil {
			d.Set("dpd_timeout_seconds", props.DpdTimeoutSeconds)
		}

		if props.Peer != nil {
			d.Set("express_route_circuit_id", props.Peer.Id)
		}

		if props.VirtualNetworkGateway2 != nil {
			d.Set("peer_virtual_network_gateway_id", props.VirtualNetworkGateway2.Id)
		}

		if props.UseLocalAzureIPAddress != nil {
			d.Set("local_azure_ip_address_enabled", props.UseLocalAzureIPAddress)
		}

		if props.LocalNetworkGateway2 != nil {
			d.Set("local_network_gateway_id", props.LocalNetworkGateway2.Id)
		}

		if props.EnableBgp != nil {
			d.Set("enable_bgp", props.EnableBgp)
		}

		if props.UsePolicyBasedTrafficSelectors != nil {
			d.Set("use_policy_based_traffic_selectors", props.UsePolicyBasedTrafficSelectors)
		}

		if props.RoutingWeight != nil {
			d.Set("routing_weight", props.RoutingWeight)
		}

		if err := d.Set("custom_bgp_addresses", flattenGatewayCustomBgpIPAddresses(props.GatewayCustomBgpIPAddresses)); err != nil {
			return fmt.Errorf("setting `custom_bgp_addresses`: %+v", err)
		}

		d.Set("connection_protocol", string(pointer.From(props.ConnectionProtocol)))

		d.Set("connection_mode", string(pointer.From(props.ConnectionMode)))

		if props.ExpressRouteGatewayBypass != nil {
			d.Set("express_route_gateway_bypass", props.ExpressRouteGatewayBypass)
		}

		if props.EnablePrivateLinkFastPath != nil {
			d.Set("private_link_fast_path_enabled", props.EnablePrivateLinkFastPath)
		}

		if props.IPsecPolicies != nil {
			ipsecPolicies := flattenVirtualNetworkGatewayConnectionIpsecPolicies(props.IPsecPolicies)

			if err := d.Set("ipsec_policy", ipsecPolicies); err != nil {
				return fmt.Errorf("setting `ipsec_policy`: %+v", err)
			}
		}

		trafficSelectorPolicies := flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(props.TrafficSelectorPolicies)
		if err := d.Set("traffic_selector_policy", trafficSelectorPolicies); err != nil {
			return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
		}

		if err := d.Set("egress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(props.EgressNatRules)); err != nil {
			return fmt.Errorf("setting `egress_nat_rule_ids`: %+v", err)
		}

		if err := d.Set("ingress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(props.IngressNatRules)); err != nil {
			return fmt.Errorf("setting `ingress_nat_rule_ids`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceVirtualNetworkGatewayConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGatewayConnections
	vnetGatewayClient := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway Connection update.")

	id, err := virtualnetworkgatewayconnections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	var virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway
	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return err
		}

		resp, err := vnetGatewayClient.Get(ctx, *gwid)
		if err != nil {
			return err
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", gwid, err)
		}

		virtualNetworkGateway = *resp.Model
	}

	if d.HasChange("authorization_key") {
		payload.Properties.AuthorizationKey = pointer.To(d.Get("authorization_key").(string))
	}

	if d.HasChange("egress_nat_rule_ids") {
		payload.Properties.EgressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(d.Get("egress_nat_rule_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("ingress_nat_rule_ids") {
		payload.Properties.IngressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(d.Get("ingress_nat_rule_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("local_network_gateway_id") {
		localNetworkGatewayId := d.Get("local_network_gateway_id").(string)
		name, err := localNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return fmt.Errorf("Getting LocalNetworkGateway Name and Group:: %+v", err)
		}

		payload.Properties.LocalNetworkGateway2 = &virtualnetworkgatewayconnections.LocalNetworkGateway{
			Id:   &localNetworkGatewayId,
			Name: &name,
			Properties: virtualnetworkgatewayconnections.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &virtualnetworkgatewayconnections.AddressSpace{},
			},
		}
	}

	if d.HasChange("enable_bgp") {
		payload.Properties.EnableBgp = pointer.To(d.Get("enable_bgp").(bool))
	}

	if d.HasChange("use_policy_based_traffic_selectors") {
		payload.Properties.UsePolicyBasedTrafficSelectors = pointer.To(d.Get("use_policy_based_traffic_selectors").(bool))
	}

	if d.HasChange("routing_weight") {
		payload.Properties.RoutingWeight = pointer.To(int64(d.Get("routing_weight").(int)))
	}

	if d.HasChange("express_route_gateway_bypass") {
		payload.Properties.ExpressRouteGatewayBypass = pointer.To(d.Get("express_route_gateway_bypass").(bool))
	}

	if d.HasChange("private_link_fast_path_enabled") {
		payload.Properties.EnablePrivateLinkFastPath = pointer.To(d.Get("private_link_fast_path_enabled").(bool))
	}

	if d.HasChange("traffic_selector_policy") {
		payload.Properties.TrafficSelectorPolicies = expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(d.Get("traffic_selector_policy").([]interface{}))
	}

	if d.HasChange("custom_bgp_addresses") {
		if virtualNetworkGateway.Properties.BgpSettings == nil || virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses == nil {
			return fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway` %s (%s) failed: get nil", *virtualNetworkGateway.Name, *virtualNetworkGateway.Id)
		}

		gatewayCustomBgpIPAddresses, err := expandGatewayCustomBgpIPAddresses(d, virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses)
		if err != nil {
			return err
		}

		payload.Properties.GatewayCustomBgpIPAddresses = gatewayCustomBgpIPAddresses
	}

	if d.HasChange("ipsec_policy") {
		payload.Properties.IPsecPolicies = expandVirtualNetworkGatewayConnectionIpsecPolicies(d.Get("ipsec_policy").([]interface{}))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if d.HasChange("shared_key") {
		if err := client.SetSharedKeyThenPoll(ctx, *id, virtualnetworkgatewayconnections.ConnectionSharedKey{
			Value: d.Get("shared_key").(string),
		}); err != nil {
			return fmt.Errorf("updating Shared Key for %s: %+v", id, err)
		}

		// Once this issue https://github.com/Azure/azure-rest-api-specs/issues/26660 is fixed, below this part will be removed
		stateConf := &pluginsdk.StateChangeConf{
			Pending:    []string{string(virtualnetworkgatewayconnections.ProvisioningStateUpdating)},
			Target:     []string{string(virtualnetworkgatewayconnections.ProvisioningStateSucceeded)},
			Refresh:    virtualNetworkGatewayConnectionStateRefreshFunc(ctx, client, *id),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
		}

		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceVirtualNetworkGatewayConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGatewayConnections
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgatewayconnections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func virtualNetworkGatewayConnectionStateRefreshFunc(ctx context.Context, client *virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionsClient, id virtualnetworkgatewayconnections.ConnectionId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties.ProvisioningState != nil {
			return res, string(pointer.From(res.Model.Properties.ProvisioningState)), nil
		}

		return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
	}
}

func getVirtualNetworkGatewayConnectionProperties(d *pluginsdk.ResourceData, virtualNetworkGateway virtualnetworkgateways.VirtualNetworkGateway) (*virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionType(d.Get("type").(string))
	connectionMode := virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionMode(d.Get("connection_mode").(string))

	props := &virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType:                 connectionType,
		ConnectionMode:                 pointer.To(connectionMode),
		EnableBgp:                      pointer.To(d.Get("enable_bgp").(bool)),
		EnablePrivateLinkFastPath:      pointer.To(d.Get("private_link_fast_path_enabled").(bool)),
		ExpressRouteGatewayBypass:      pointer.To(d.Get("express_route_gateway_bypass").(bool)),
		UsePolicyBasedTrafficSelectors: pointer.To(d.Get("use_policy_based_traffic_selectors").(bool)),
	}

	if virtualNetworkGateway.Name != nil && virtualNetworkGateway.Id != nil {
		props.VirtualNetworkGateway1 = virtualnetworkgatewayconnections.VirtualNetworkGateway{
			Id:   virtualNetworkGateway.Id,
			Name: virtualNetworkGateway.Name,
			Properties: virtualnetworkgatewayconnections.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]virtualnetworkgatewayconnections.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		authorizationKey := v.(string)
		props.AuthorizationKey = &authorizationKey
	}

	if v, ok := d.GetOk("dpd_timeout_seconds"); ok {
		props.DpdTimeoutSeconds = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("express_route_circuit_id"); ok {
		expressRouteCircuitId := v.(string)
		props.Peer = &virtualnetworkgatewayconnections.SubResource{
			Id: &expressRouteCircuitId,
		}
	}

	if v, ok := d.GetOk("egress_nat_rule_ids"); ok {
		props.EgressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("ingress_nat_rule_ids"); ok {
		props.IngressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("peer_virtual_network_gateway_id"); ok {
		gwid, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(v.(string))
		if err != nil {
			return nil, err
		}
		props.VirtualNetworkGateway2 = &virtualnetworkgatewayconnections.VirtualNetworkGateway{
			Id:   pointer.To(gwid.ID()),
			Name: &gwid.VirtualNetworkGatewayName,
			Properties: virtualnetworkgatewayconnections.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]virtualnetworkgatewayconnections.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("local_azure_ip_address_enabled"); ok {
		props.UseLocalAzureIPAddress = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("local_network_gateway_id"); ok {
		localNetworkGatewayId := v.(string)
		name, err := localNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("Getting LocalNetworkGateway Name and Group:: %+v", err)
		}

		props.LocalNetworkGateway2 = &virtualnetworkgatewayconnections.LocalNetworkGateway{
			Id:   &localNetworkGatewayId,
			Name: &name,
			Properties: virtualnetworkgatewayconnections.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &virtualnetworkgatewayconnections.AddressSpace{},
			},
		}
	}

	if v, ok := d.GetOk("routing_weight"); ok {
		props.RoutingWeight = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("shared_key"); ok {
		props.SharedKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("connection_protocol"); ok {
		connectionProtocol := v.(string)
		props.ConnectionProtocol = pointer.To(virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionProtocol(connectionProtocol))
	}

	if v, ok := d.GetOk("traffic_selector_policy"); ok {
		props.TrafficSelectorPolicies = expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(v.([]interface{}))
	}

	if v, ok := d.GetOk("ipsec_policy"); ok {
		props.IPsecPolicies = expandVirtualNetworkGatewayConnectionIpsecPolicies(v.([]interface{}))
	}

	if utils.NormaliseNilableBool(props.EnableBgp) {
		if _, ok := d.GetOk("custom_bgp_addresses"); ok {
			if virtualNetworkGateway.Properties.BgpSettings == nil || virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses == nil {
				return nil, fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway` %s (%s) failed: get nil", *virtualNetworkGateway.Name, *virtualNetworkGateway.Id)
			}

			gatewayCustomBgpIPAddresses, err := expandGatewayCustomBgpIPAddresses(d, virtualNetworkGateway.Properties.BgpSettings.BgpPeeringAddresses)
			if err != nil {
				return nil, err
			}

			props.GatewayCustomBgpIPAddresses = gatewayCustomBgpIPAddresses
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeExpressRoute {
		if props.Peer == nil || props.Peer.Id == nil {
			return nil, fmt.Errorf("`express_route_circuit_id` must be specified when `type` is set to `ExpressRoute`")
		}
		if d.Get("private_link_fast_path_enabled").(bool) && !d.Get("express_route_gateway_bypass").(bool) {
			return nil, fmt.Errorf("`express_route_gateway_bypass` must be enabled when `private_link_fast_path_enabled` is set to `true`")
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec {
		if props.LocalNetworkGateway2 == nil || props.LocalNetworkGateway2.Id == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` must be specified when `type` is set to `IPsec`")
		}
	}

	if props.ConnectionType == virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeVnetTwoVnet {
		if props.VirtualNetworkGateway2 == nil || props.VirtualNetworkGateway2.Id == nil {
			return nil, fmt.Errorf("`peer_virtual_network_gateway_id` must be specified when `type` is set to `Vnet2Vnet`")
		}
	}

	if props.GatewayCustomBgpIPAddresses != nil && props.ConnectionType != virtualnetworkgatewayconnections.VirtualNetworkGatewayConnectionTypeIPsec {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `type` is set to `IPsec`")
	}

	if props.GatewayCustomBgpIPAddresses != nil && virtualNetworkGateway.Properties.ActiveActive == pointer.To(false) {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `azurerm_virtual_network_gateway` `active_active` is set enabled`")
	}

	return props, nil
}

func expandVirtualNetworkGatewayConnectionIpsecPolicies(schemaIpsecPolicies []interface{}) *[]virtualnetworkgatewayconnections.IPsecPolicy {
	ipsecPolicies := make([]virtualnetworkgatewayconnections.IPsecPolicy, 0, len(schemaIpsecPolicies))

	for _, d := range schemaIpsecPolicies {
		schemaIpsecPolicy := d.(map[string]interface{})
		ipsecPolicy := &virtualnetworkgatewayconnections.IPsecPolicy{}

		if dhGroup, ok := schemaIpsecPolicy["dh_group"].(string); ok && dhGroup != "" {
			ipsecPolicy.DhGroup = virtualnetworkgatewayconnections.DhGroup(dhGroup)
		}

		if ikeEncryption, ok := schemaIpsecPolicy["ike_encryption"].(string); ok && ikeEncryption != "" {
			ipsecPolicy.IkeEncryption = virtualnetworkgatewayconnections.IkeEncryption(ikeEncryption)
		}

		if ikeIntegrity, ok := schemaIpsecPolicy["ike_integrity"].(string); ok && ikeIntegrity != "" {
			ipsecPolicy.IkeIntegrity = virtualnetworkgatewayconnections.IkeIntegrity(ikeIntegrity)
		}

		if ipsecEncryption, ok := schemaIpsecPolicy["ipsec_encryption"].(string); ok && ipsecEncryption != "" {
			ipsecPolicy.IPsecEncryption = virtualnetworkgatewayconnections.IPsecEncryption(ipsecEncryption)
		}

		if ipsecIntegrity, ok := schemaIpsecPolicy["ipsec_integrity"].(string); ok && ipsecIntegrity != "" {
			ipsecPolicy.IPsecIntegrity = virtualnetworkgatewayconnections.IPsecIntegrity(ipsecIntegrity)
		}

		if pfsGroup, ok := schemaIpsecPolicy["pfs_group"].(string); ok && pfsGroup != "" {
			ipsecPolicy.PfsGroup = virtualnetworkgatewayconnections.PfsGroup(pfsGroup)
		}

		if v, ok := schemaIpsecPolicy["sa_datasize"].(int); ok {
			ipsecPolicy.SaDataSizeKilobytes = int64(v)
		}

		if v, ok := schemaIpsecPolicy["sa_lifetime"].(int); ok {
			ipsecPolicy.SaLifeTimeSeconds = int64(v)
		}

		ipsecPolicies = append(ipsecPolicies, *ipsecPolicy)
	}

	return &ipsecPolicies
}

func expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(schemaTrafficSelectorPolicies []interface{}) *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy {
	trafficSelectorPolicies := make([]virtualnetworkgatewayconnections.TrafficSelectorPolicy, 0, len(schemaTrafficSelectorPolicies))

	for _, d := range schemaTrafficSelectorPolicies {
		schemaTrafficSelectorPolicy := d.(map[string]interface{})
		trafficSelectorPolicy := &virtualnetworkgatewayconnections.TrafficSelectorPolicy{}
		if localAddressRanges, ok := schemaTrafficSelectorPolicy["local_address_cidrs"].([]interface{}); ok {
			trafficSelectorPolicy.LocalAddressRanges = pointer.From(utils.ExpandStringSlice(localAddressRanges))
		}
		if remoteAddressRanges, ok := schemaTrafficSelectorPolicy["remote_address_cidrs"].([]interface{}); ok {
			trafficSelectorPolicy.RemoteAddressRanges = pointer.From(utils.ExpandStringSlice(remoteAddressRanges))
		}

		trafficSelectorPolicies = append(trafficSelectorPolicies, *trafficSelectorPolicy)
	}

	return &trafficSelectorPolicies
}

func expandGatewayCustomBgpIPAddresses(d *pluginsdk.ResourceData, bgpPeeringAddresses *[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress) (*[]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration, error) {
	customBgpIpAddresses := make([]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration, 0)

	bgpAddresses := d.Get("custom_bgp_addresses").([]interface{})
	if len(bgpAddresses) == 0 {
		return &customBgpIpAddresses, nil
	}

	bgAs := bgpAddresses[0].(map[string]interface{})
	primaryAddress := bgAs["primary"].(string)
	secondaryAddress := bgAs["secondary"].(string)

	var primaryIpConfiguration string
	var secondaryIpConfiguration string

	if bgpPeeringAddresses == nil {
		return &customBgpIpAddresses, fmt.Errorf("retrieving BGP peering address from `virtual_network_gateway`, addresses returned nil")
	}

	for _, address := range *bgpPeeringAddresses {
		if address.CustomBgpIPAddresses == nil {
			continue
		}

		for _, ip := range *address.CustomBgpIPAddresses {
			if address.IPconfigurationId == nil {
				continue
			}

			if ip == primaryAddress {
				primaryIpConfiguration = *address.IPconfigurationId
			} else if ip == secondaryAddress {
				secondaryIpConfiguration = *address.IPconfigurationId
			}
		}
	}

	if len(primaryIpConfiguration) == 0 || (secondaryAddress != "" && len(secondaryIpConfiguration) == 0) {
		return &customBgpIpAddresses, fmt.Errorf("primary or secondary address not found at `virtual_network_gateway` configuration `bgp_settings` `peering_addresses`")
	}

	customBgpIpAddresses = append(customBgpIpAddresses, virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration{
		IPConfigurationId:  primaryIpConfiguration,
		CustomBgpIPAddress: primaryAddress,
	})

	if secondaryAddress != "" {
		customBgpIpAddresses = append(customBgpIpAddresses, virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration{
			IPConfigurationId:  secondaryIpConfiguration,
			CustomBgpIPAddress: secondaryAddress,
		})
	}

	return &customBgpIpAddresses, nil
}

func flattenVirtualNetworkGatewayConnectionIpsecPolicies(ipsecPolicies *[]virtualnetworkgatewayconnections.IPsecPolicy) []interface{} {
	schemaIpsecPolicies := make([]interface{}, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]interface{})

			schemaIpsecPolicy["dh_group"] = string(ipsecPolicy.DhGroup)
			schemaIpsecPolicy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
			schemaIpsecPolicy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
			schemaIpsecPolicy["ipsec_encryption"] = string(ipsecPolicy.IPsecEncryption)
			schemaIpsecPolicy["ipsec_integrity"] = string(ipsecPolicy.IPsecIntegrity)
			schemaIpsecPolicy["pfs_group"] = string(ipsecPolicy.PfsGroup)
			schemaIpsecPolicy["sa_datasize"] = int(ipsecPolicy.SaDataSizeKilobytes)
			schemaIpsecPolicy["sa_lifetime"] = int(ipsecPolicy.SaLifeTimeSeconds)

			schemaIpsecPolicies = append(schemaIpsecPolicies, schemaIpsecPolicy)
		}
	}

	return schemaIpsecPolicies
}

func flattenGatewayCustomBgpIPAddresses(gatewayCustomBgpIPAddresses *[]virtualnetworkgatewayconnections.GatewayCustomBgpIPAddressIPConfiguration) interface{} {
	customBgpIpAddresses := make([]interface{}, 0)
	if gatewayCustomBgpIPAddresses == nil || len(*gatewayCustomBgpIPAddresses) == 0 {
		return customBgpIpAddresses
	}

	customBgpIpAddress := map[string]interface{}{}
	for k, v := range *gatewayCustomBgpIPAddresses {
		if k == 0 && v.CustomBgpIPAddress != "" {
			customBgpIpAddress["primary"] = v.CustomBgpIPAddress
		} else if k == 1 && v.CustomBgpIPAddress != "" {
			customBgpIpAddress["secondary"] = v.CustomBgpIPAddress
		}
	}

	return append(customBgpIpAddresses, customBgpIpAddress)
}

func flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(trafficSelectorPolicies *[]virtualnetworkgatewayconnections.TrafficSelectorPolicy) []interface{} {
	schemaTrafficSelectorPolicies := make([]interface{}, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]interface{}{
				"local_address_cidrs":  trafficSelectorPolicy.LocalAddressRanges,
				"remote_address_cidrs": trafficSelectorPolicy.RemoteAddressRanges,
			})
		}
	}

	return schemaTrafficSelectorPolicies
}

func expandVirtualNetworkGatewayConnectionNatRuleIds(input []interface{}) *[]virtualnetworkgatewayconnections.SubResource {
	results := make([]virtualnetworkgatewayconnections.SubResource, 0)

	for _, item := range input {
		results = append(results, virtualnetworkgatewayconnections.SubResource{
			Id: pointer.To(item.(string)),
		})
	}

	return &results
}

func flattenVirtualNetworkGatewayConnectionNatRuleIds(input *[]virtualnetworkgatewayconnections.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var id string
		if item.Id != nil {
			id = *item.Id
		}

		results = append(results, id)
	}

	return results
}
