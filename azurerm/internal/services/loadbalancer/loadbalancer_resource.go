package loadbalancer

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
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

						"availability_zone": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							//Default:  "Zone-Redundant",
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"No-Zone",
								"1",
								"2",
								"3",
								"Zone-Redundant",
							}, false),
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

						// TODO - 3.0 make Computed only
						"zones": {
							Type:       pluginsdk.TypeList,
							Optional:   true,
							Computed:   true,
							ForceNew:   true,
							Deprecated: "This property has been deprecated in favour of `availability_zone` due to a breaking behavioural change in Azure: https://azure.microsoft.com/en-us/updates/zone-behavior-change/",
							MaxItems:   1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

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

			"tags": tags.Schema(),
		},
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
				return fmt.Errorf("checking for presence of existing Load Balancer %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_lb", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := network.LoadBalancerSku{
		Name: network.LoadBalancerSkuName(d.Get("sku").(string)),
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
		return fmt.Errorf("creating/updating Load Balancer %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating/Updating Load Balancer %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.Name)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer By ID: %+v", err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			if err := d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs)); err != nil {
				return fmt.Errorf("Error flattening `frontend_ip_configuration`: %+v", err)
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
		return fmt.Errorf("deleting Load Balancer %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Load Balancer %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

		name := data["name"].(string)
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
		if !strings.EqualFold(sku, string(network.LoadBalancerSkuNameStandard)) {
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
		frontEndConfig := network.FrontendIPConfiguration{
			Name:                                    &name,
			FrontendIPConfigurationPropertiesFormat: &properties,
			Zones:                                   zones,
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
		ipConfig := make(map[string]interface{})

		if config.Name != nil {
			ipConfig["name"] = *config.Name
		}

		if config.ID != nil {
			ipConfig["id"] = *config.ID
		}

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
		ipConfig["availability_zone"] = availabilityZones
		ipConfig["zones"] = zonesDeprecated

		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			ipConfig["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				ipConfig["subnet_id"] = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				ipConfig["private_ip_address"] = *pip
			}

			if props.PrivateIPAddressVersion != "" {
				ipConfig["private_ip_address_version"] = string(props.PrivateIPAddressVersion)
			}

			if pip := props.PublicIPAddress; pip != nil {
				ipConfig["public_ip_address_id"] = *pip.ID
			}

			if pip := props.PublicIPPrefix; pip != nil {
				ipConfig["public_ip_prefix_id"] = *pip.ID
			}

			loadBalancingRules := make([]interface{}, 0)
			if rules := props.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					loadBalancingRules = append(loadBalancingRules, *rule.ID)
				}
			}
			ipConfig["load_balancer_rules"] = pluginsdk.NewSet(pluginsdk.HashString, loadBalancingRules)

			inboundNatRules := make([]interface{}, 0)
			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, *rule.ID)
				}
			}
			ipConfig["inbound_nat_rules"] = pluginsdk.NewSet(pluginsdk.HashString, inboundNatRules)

			outboundRules := make([]interface{}, 0)
			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					outboundRules = append(outboundRules, *rule.ID)
				}
			}
			ipConfig["outbound_rules"] = pluginsdk.NewSet(pluginsdk.HashString, outboundRules)
		}

		result = append(result, ipConfig)
	}
	return result
}
