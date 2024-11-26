// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerCreateUpdate,
		Read:   resourceArmLoadBalancerRead,
		Update: resourceArmLoadBalancerCreateUpdate,
		Delete: resourceArmLoadBalancerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := loadbalancers.ParseLoadBalancerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceArmLoadBalancerSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIf("frontend_ip_configuration", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				old, new := d.GetChange("frontend_ip_configuration")
				if len(old.([]interface{})) == 0 && len(new.([]interface{})) > 0 || len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0 {
					return false
				} else {
					for i, nc := range new.([]interface{}) {
						dataNew := nc.(map[string]interface{})
						for _, oc := range old.([]interface{}) {
							dataOld := oc.(map[string]interface{})
							if dataOld["name"].(string) == dataNew["name"].(string) {
								if !reflect.DeepEqual(dataOld["zones"].(*pluginsdk.Set).List(), dataNew["zones"].(*pluginsdk.Set).List()) {
									// set ForceNew to true when the `frontend_ip_configuration.#.zones` is changed.
									d.ForceNew("frontend_ip_configuration." + strconv.Itoa(i) + ".zones")
									break
								}
							}
						}
					}
				}
				return false
			}),
		),
	}
}

func resourceArmLoadBalancerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Load Balancer creation.")

	id := loadbalancers.NewLoadBalancerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	plbId := loadbalancers.ProviderLoadBalancerId(id)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_lb", id.ID())
		}
	}

	if strings.EqualFold(d.Get("sku_tier").(string), string(loadbalancers.LoadBalancerSkuTierGlobal)) {
		if !strings.EqualFold(d.Get("sku").(string), string(loadbalancers.LoadBalancerSkuNameStandard)) {
			return fmt.Errorf("global load balancing is only supported for standard SKU load balancers")
		}
	}

	sku := loadbalancers.LoadBalancerSku{
		Name: pointer.To(loadbalancers.LoadBalancerSkuName(d.Get("sku").(string))),
		Tier: pointer.To(loadbalancers.LoadBalancerSkuTier(d.Get("sku_tier").(string))),
	}

	properties := loadbalancers.LoadBalancerPropertiesFormat{}

	if _, ok := d.GetOk("frontend_ip_configuration"); ok {
		properties.FrontendIPConfigurations = expandAzureRmLoadBalancerFrontendIpConfigurations(d)
	}

	loadBalancer := loadbalancers.LoadBalancer{
		Name:             pointer.To(id.LoadBalancerName),
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:             tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:              pointer.To(sku),
		Properties:       pointer.To(properties),
	}

	err := client.CreateOrUpdateThenPoll(ctx, plbId, loadBalancer)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRead(d, meta)
}

func resourceArmLoadBalancerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	resp, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LoadBalancerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		if sku := model.Sku; sku != nil {
			d.Set("sku", string(pointer.From(sku.Name)))
			d.Set("sku_tier", string(pointer.From(sku.Tier)))
		}

		if props := model.Properties; props != nil {
			if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
				if err := d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs)); err != nil {
					return fmt.Errorf("flattening `frontend_ip_configuration`: %+v", err)
				}

				privateIpAddress := ""
				privateIpAddresses := make([]string, 0)
				for _, config := range *feipConfigs {
					if feipProps := config.Properties; feipProps != nil {
						if ip := feipProps.PrivateIPAddress; ip != nil {
							if privateIpAddress == "" {
								privateIpAddress = pointer.From(feipProps.PrivateIPAddress)
							}

							privateIpAddresses = append(privateIpAddresses, *feipProps.PrivateIPAddress)
						}
					}
				}

				d.Set("private_ip_address", privateIpAddress)
				d.Set("private_ip_addresses", privateIpAddresses)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceArmLoadBalancerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}

	err = client.DeleteThenPoll(ctx, plbId)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRmLoadBalancerFrontendIpConfigurations(d *pluginsdk.ResourceData) *[]loadbalancers.FrontendIPConfiguration {
	configs := d.Get("frontend_ip_configuration").([]interface{})
	frontEndConfigs := make([]loadbalancers.FrontendIPConfiguration, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		properties := loadbalancers.FrontendIPConfigurationPropertiesFormat{}

		if v := data["private_ip_address_allocation"].(string); v != "" {
			properties.PrivateIPAllocationMethod = pointer.To(loadbalancers.IPAllocationMethod(v))
		}

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &loadbalancers.SubResource{
				Id: pointer.To(v),
			}
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = pointer.To(v)
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &loadbalancers.PublicIPAddress{
				Id: pointer.To(v),
			}
		}

		if v := data["public_ip_prefix_id"].(string); v != "" {
			properties.PublicIPPrefix = &loadbalancers.SubResource{
				Id: pointer.To(v),
			}
		}

		if v := data["subnet_id"].(string); v != "" {
			properties.PrivateIPAddressVersion = pointer.To(loadbalancers.IPVersionIPvFour)
			if v := data["private_ip_address_version"].(string); v != "" {
				properties.PrivateIPAddressVersion = pointer.To(loadbalancers.IPVersion(v))
			}
			properties.Subnet = &loadbalancers.Subnet{
				Id: pointer.To(v),
			}
		}

		frontEndConfig := loadbalancers.FrontendIPConfiguration{
			Name:       pointer.To(data["name"].(string)),
			Properties: pointer.To(properties),
		}

		zones := zones.ExpandUntyped(data["zones"].(*pluginsdk.Set).List())
		if len(zones) > 0 {
			frontEndConfig.Zones = &zones
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs
}

func flattenLoadBalancerFrontendIpConfiguration(ipConfigs *[]loadbalancers.FrontendIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		name := ""
		if config.Name != nil {
			name = *config.Name
		}

		id := ""
		if config.Id != nil {
			id = *config.Id
		}

		var inboundNatRules []interface{}
		var loadBalancingRules []interface{}
		var outboundRules []interface{}
		gatewayLoadBalancerId := ""
		publicIpPrefixId := ""
		privateIPAllocationMethod := ""
		publicIpAddressId := ""
		privateIpAddressVersion := ""
		subnetId := ""
		privateIpAddress := ""

		if props := config.Properties; props != nil {
			privateIPAllocationMethod = string(pointer.From(props.PrivateIPAllocationMethod))

			if props.GatewayLoadBalancer != nil {
				gatewayLoadBalancerId = pointer.From(props.GatewayLoadBalancer.Id)
			}

			if subnet := props.Subnet; subnet != nil {
				subnetId = pointer.From(subnet.Id)
			}
			privateIpAddress = pointer.From(props.PrivateIPAddress)
			privateIpAddressVersion = string(pointer.From(props.PrivateIPAddressVersion))

			if pip := props.PublicIPAddress; pip != nil {
				publicIpAddressId = pointer.From(pip.Id)
			}

			if pip := props.PublicIPPrefix; pip != nil {
				publicIpPrefixId = pointer.From(pip.Id)
			}

			if rules := props.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}

					loadBalancingRules = append(loadBalancingRules, pointer.From(rule.Id))
				}
			}

			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, pointer.From(rule.Id))
				}
			}

			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					outboundRules = append(outboundRules, pointer.From(rule.Id))
				}
			}
		}

		out := map[string]interface{}{
			"gateway_load_balancer_frontend_ip_configuration_id": gatewayLoadBalancerId,
			"id":                            id,
			"inbound_nat_rules":             pluginsdk.NewSet(pluginsdk.HashString, inboundNatRules),
			"load_balancer_rules":           pluginsdk.NewSet(pluginsdk.HashString, loadBalancingRules),
			"name":                          name,
			"outbound_rules":                pluginsdk.NewSet(pluginsdk.HashString, outboundRules),
			"public_ip_address_id":          publicIpAddressId,
			"private_ip_address":            privateIpAddress,
			"private_ip_address_version":    privateIpAddressVersion,
			"private_ip_address_allocation": privateIPAllocationMethod,
			"public_ip_prefix_id":           publicIpPrefixId,
			"subnet_id":                     subnetId,
			"zones":                         zones.FlattenUntyped(config.Zones),
		}

		result = append(result, out)
	}
	return result
}

func resourceArmLoadBalancerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(loadbalancers.LoadBalancerSkuNameStandard),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.LoadBalancerSkuNameBasic),
				string(loadbalancers.LoadBalancerSkuNameStandard),
				string(loadbalancers.LoadBalancerSkuNameGateway),
			}, false),
		},

		"sku_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(loadbalancers.LoadBalancerSkuTierRegional),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(loadbalancers.LoadBalancerSkuTierRegional),
				string(loadbalancers.LoadBalancerSkuTierGlobal),
			}, false),
		},

		"frontend_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // TODO: why is this computed?
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // TODO: remove computed in 4.0 and use ignore_changes
						ValidateFunc: validation.Any(
							validation.IsIPAddress,
							validation.StringIsEmpty,
						),
					},

					"private_ip_address_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // TODO: why is this computed?
						ValidateFunc: validation.StringInSlice([]string{
							string(loadbalancers.IPVersionIPvFour),
							string(loadbalancers.IPVersionIPvSix),
						}, false),
					},

					"public_ip_address_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // TODO: why is this computed?
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"public_ip_prefix_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
					},

					"private_ip_address_allocation": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(loadbalancers.IPAllocationMethodDynamic),
							string(loadbalancers.IPAllocationMethodStatic),
						}, true),
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"gateway_load_balancer_frontend_ip_configuration_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: loadbalancers.ValidateFrontendIPConfigurationID,
					},

					"load_balancer_rules": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
					},

					"inbound_nat_rules": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
					},

					"outbound_rules": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						Set: pluginsdk.HashString,
					},

					"zones": commonschema.ZonesMultipleOptional(),

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"private_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"private_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}
