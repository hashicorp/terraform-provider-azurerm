package loadbalancer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/state"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmLoadBalancer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerCreateUpdate,
		Read:   resourceArmLoadBalancerRead,
		Update: resourceArmLoadBalancerCreateUpdate,
		Delete: resourceArmLoadBalancerDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LoadBalancerID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(network.LoadBalancerSkuNameBasic),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.LoadBalancerSkuNameBasic),
					string(network.LoadBalancerSkuNameStandard),
					string(network.LoadBalancerSkuNameGateway),
				}, true),
				// TODO - 3.0 remove this property
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(network.LoadBalancerSkuTierRegional),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.LoadBalancerSkuTierRegional),
					string(network.LoadBalancerSkuTierGlobal),
				}, false),
			},

			"frontend_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: func() map[string]*pluginsdk.Schema {
						s := map[string]*pluginsdk.Schema{
							"name": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"subnet_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: azure.ValidateResourceIDOrEmpty,
							},

							"private_ip_address": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.Any(
									validation.IsIPAddress,
									validation.StringIsEmpty,
								),
							},

							"private_ip_address_version": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.IPVersionIPv4),
									string(network.IPVersionIPv6),
								}, false),
							},

							"public_ip_address_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: azure.ValidateResourceIDOrEmpty,
							},

							"public_ip_prefix_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: azure.ValidateResourceIDOrEmpty,
							},

							"private_ip_address_allocation": {
								Type:     pluginsdk.TypeString,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.IPAllocationMethodDynamic),
									string(network.IPAllocationMethodStatic),
								}, true),
								StateFunc:        state.IgnoreCase,
								DiffSuppressFunc: suppress.CaseDifference,
							},

							"gateway_load_balancer_frontend_ip_configuration_id": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								Computed:     true,
								ValidateFunc: validate.LoadBalancerFrontendIpConfigurationID,
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

							"id": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						}

						if features.ThreePointOhBeta() {
							s["zones"] = commonschema.ZonesMultipleOptional()
						} else {
							s["availability_zone"] = &pluginsdk.Schema{
								Type:     pluginsdk.TypeString,
								Optional: true,
								// Default:  "Zone-Redundant",
								Computed: true,
								ValidateFunc: validation.StringInSlice([]string{
									"No-Zone",
									"1",
									"2",
									"3",
									"Zone-Redundant",
								}, false),
							}

							s["zones"] = &pluginsdk.Schema{
								Type:       pluginsdk.TypeList,
								Optional:   true,
								Computed:   true,
								Deprecated: "This property has been deprecated in favour of `availability_zone` due to a breaking behavioural change in Azure: https://azure.microsoft.com/en-us/updates/zone-behavior-change/",
								MaxItems:   1,
								Elem: &pluginsdk.Schema{
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							}
						}

						return s
					}(),
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if features.ThreePointOhBeta() {
				return nil
			}
			if ok := d.HasChange("frontend_ip_configuration"); ok {
				configs := d.Get("frontend_ip_configuration").([]interface{})

				for index := range configs {
					if d.HasChange(fmt.Sprintf("frontend_ip_configuration.%d.availability_zone", index)) && !d.HasChange(fmt.Sprintf("frontend_ip_configuration.%d.name", index)) {
						return fmt.Errorf("in place change of the `frontend_ip_configuration.%[1]d.availability_zone` is not allowed. It is allowed to do this while also changing `frontend_ip_configuration.%[1]d.name`", index)
					}

					if d.HasChange(fmt.Sprintf("frontend_ip_configuration.%d.zones", index)) && !d.HasChange(fmt.Sprintf("frontend_ip_configuration.%d.name", index)) {
						return fmt.Errorf("in place change of the `frontend_ip_configuration.%[1]d.zones` is not allowed. It is allowed to do this while also changing `frontend_ip_configuration.%[1]d.name`", index)
					}
				}
			}

			return nil
		}),
	}
}

func resourceArmLoadBalancerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Load Balancer creation.")

	id := parse.NewLoadBalancerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_lb", id.ID())
		}
	}

	if strings.EqualFold(d.Get("sku_tier").(string), string(network.LoadBalancerSkuTierGlobal)) {
		if !strings.EqualFold(d.Get("sku").(string), string(network.LoadBalancerSkuNameStandard)) {
			return fmt.Errorf("global load balancing is only supported for standard SKU load balancers")
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := network.LoadBalancerSku{
		Name: network.LoadBalancerSkuName(d.Get("sku").(string)),
		Tier: network.LoadBalancerSkuTier(d.Get("sku_tier").(string)),
	}
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	properties := network.LoadBalancerPropertiesFormat{}

	if _, ok := d.GetOk("frontend_ip_configuration"); ok {
		frontendIPConfigurations, err := expandAzureRmLoadBalancerFrontendIpConfigurations(d)
		if err != nil {
			return err
		}
		properties.FrontendIPConfigurations = frontendIPConfigurations
	}

	loadBalancer := network.LoadBalancer{
		Name:                         utils.String(id.Name),
		Location:                     utils.String(location),
		Tags:                         expandedTags,
		Sku:                          &sku,
		LoadBalancerPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, loadBalancer)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerRead(d, meta)
}

func resourceArmLoadBalancerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
		d.Set("sku_tier", string(sku.Tier))
	}

	if props := resp.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			if err := d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs)); err != nil {
				return fmt.Errorf("flattening `frontend_ip_configuration`: %+v", err)
			}

			privateIpAddress := ""
			privateIpAddresses := make([]string, 0)
			for _, config := range *feipConfigs {
				if feipProps := config.FrontendIPConfigurationPropertiesFormat; feipProps != nil {
					if ip := feipProps.PrivateIPAddress; ip != nil {
						if privateIpAddress == "" {
							privateIpAddress = *feipProps.PrivateIPAddress
						}

						privateIpAddresses = append(privateIpAddresses, *feipProps.PrivateIPAddress)
					}
				}
			}

			d.Set("private_ip_address", privateIpAddress)
			d.Set("private_ip_addresses", privateIpAddresses)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLoadBalancerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRmLoadBalancerFrontendIpConfigurations(d *pluginsdk.ResourceData) (*[]network.FrontendIPConfiguration, error) {
	configs := d.Get("frontend_ip_configuration").([]interface{})
	frontEndConfigs := make([]network.FrontendIPConfiguration, 0, len(configs))
	sku := d.Get("sku").(string)

	for index, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		privateIpAllocationMethod := data["private_ip_address_allocation"].(string)
		properties := network.FrontendIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: network.IPAllocationMethod(privateIpAllocationMethod),
		}

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &network.SubResource{
				ID: utils.String(v),
			}
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		subnetSet := false
		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &network.PublicIPAddress{
				ID: &v,
			}
		}

		if v := data["public_ip_prefix_id"].(string); v != "" {
			properties.PublicIPPrefix = &network.SubResource{
				ID: &v,
			}
		}

		if v := data["subnet_id"].(string); v != "" {
			subnetSet = true
			properties.PrivateIPAddressVersion = network.IPVersionIPv4
			if v := data["private_ip_address_version"].(string); v != "" {
				properties.PrivateIPAddressVersion = network.IPVersion(v)
			}
			properties.Subnet = &network.Subnet{
				ID: &v,
			}
		}

		frontEndConfig := network.FrontendIPConfiguration{
			Name:                                    utils.String(data["name"].(string)),
			FrontendIPConfigurationPropertiesFormat: &properties,
		}

		if features.ThreePointOhBeta() {
			zones := zones.Expand(data["zones"].(*schema.Set).List())
			if len(zones) > 0 {
				frontEndConfig.Zones = &zones
			}
		} else {
			// TODO - get zone list for each location by Resource API, instead of hardcode
			zones := &[]string{"1", "2"}
			zonesSet := false
			// TODO - Remove in 3.0
			if deprecatedZonesRaw, ok := d.GetOk(fmt.Sprintf("frontend_ip_configuration.%d.zones", index)); ok {
				zonesSet = true
				deprecatedZones := azure.ExpandZones(deprecatedZonesRaw.([]interface{}))
				if deprecatedZones != nil {
					zones = deprecatedZones
				}
			}

			if availabilityZones, ok := d.GetOk(fmt.Sprintf("frontend_ip_configuration.%d.availability_zone", index)); ok {
				zonesSet = true
				switch availabilityZones.(string) {
				case "1", "2", "3":
					zones = &[]string{availabilityZones.(string)}
				case "Zone-Redundant":
					zones = &[]string{"1", "2"}
				case "No-Zone":
					zones = &[]string{}
				}
			}
			if strings.EqualFold(sku, string(network.LoadBalancerSkuNameBasic)) {
				if zonesSet && len(*zones) > 0 {
					return nil, fmt.Errorf("Availability Zones are not available on the `Basic` SKU")
				}
				zones = &[]string{}
			} else if !subnetSet {
				if zonesSet && len(*zones) > 0 {
					return nil, fmt.Errorf("Networking supports zones only for frontendIpconfigurations which reference a subnet.")
				}
				zones = &[]string{}
			}

			frontEndConfig.Zones = zones
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs, nil
}

func flattenLoadBalancerFrontendIpConfiguration(ipConfigs *[]network.FrontendIPConfiguration) []interface{} {
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
		if config.ID != nil {
			id = *config.ID
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

		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			privateIPAllocationMethod = string(props.PrivateIPAllocationMethod)

			if props.GatewayLoadBalancer != nil && props.GatewayLoadBalancer.ID != nil {
				gatewayLoadBalancerId = *props.GatewayLoadBalancer.ID
			}

			if subnet := props.Subnet; subnet != nil {
				subnetId = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				privateIpAddress = *pip
			}

			if props.PrivateIPAddressVersion != "" {
				privateIpAddressVersion = string(props.PrivateIPAddressVersion)
			}

			if pip := props.PublicIPAddress; pip != nil {
				publicIpAddressId = *pip.ID
			}

			if pip := props.PublicIPPrefix; pip != nil {
				publicIpPrefixId = *pip.ID
			}

			if rules := props.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.ID == nil {
						continue
					}

					loadBalancingRules = append(loadBalancingRules, *rule.ID)
				}
			}

			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, *rule.ID)
				}
			}

			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					outboundRules = append(outboundRules, *rule.ID)
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
		}

		if features.ThreePointOhBeta() {
			out["zones"] = pluginsdk.NewSet(pluginsdk.HashString, zones.Flatten(config.Zones))
		} else {
			availabilityZones := "No-Zone"
			zonesDeprecated := make([]string, 0)
			if config.Zones != nil {
				if len(*config.Zones) > 1 {
					availabilityZones = "Zone-Redundant"
				}
				if len(*config.Zones) == 1 {
					zones := *config.Zones
					availabilityZones = zones[0]
					zonesDeprecated = zones
				}
			}
			out["availability_zone"] = availabilityZones
			out["zones"] = zonesDeprecated
		}
		result = append(result, out)
	}
	return result
}
