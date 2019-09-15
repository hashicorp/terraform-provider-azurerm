package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerCreateUpdate,
		Read:   resourceArmLoadBalancerRead,
		Update: resourceArmLoadBalancerCreateUpdate,
		Delete: resourceArmLoadBalancerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"private_ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.IPv4AddressOrEmpty,
						},

						"public_ip_address_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"public_ip_prefix_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Dynamic),
								string(network.Static),
							}, true),
							StateFunc:        state.IgnoreCase,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"load_balancer_rules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
							Set: schema.HashString,
						},

						"inbound_nat_rules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
							Set: schema.HashString,
						},

						"outbound_rules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
							Set: schema.HashString,
						},

						"zones": azure.SchemaSingleZone(),
					},
				},
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmLoadBalancerCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Load Balancer creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Load Balancer %q (Resource Group %q): %s", name, resGroup, err)
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
		properties.FrontendIPConfigurations = expandAzureRmLoadBalancerFrontendIpConfigurations(d)
	}

	loadBalancer := network.LoadBalancer{
		Name:                         utils.String(name),
		Location:                     utils.String(location),
		Tags:                         expandedTags,
		Sku:                          &sku,
		LoadBalancerPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error Retrieving Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Load Balancer %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLoadBalancerRead(d, meta)
}

func resourceArmLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Id(), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID %q: %+v", d.Id(), err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] Load Balancer %q not found. Removing from state", d.Id())
		return nil
	}

	d.Set("name", loadBalancer.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := loadBalancer.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := loadBalancer.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
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

	return tags.FlattenAndSet(d, loadBalancer.Tags)
}

func resourceArmLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.LoadBalancersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["loadBalancers"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deleting Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandAzureRmLoadBalancerFrontendIpConfigurations(d *schema.ResourceData) *[]network.FrontendIPConfiguration {
	configs := d.Get("frontend_ip_configuration").([]interface{})
	frontEndConfigs := make([]network.FrontendIPConfiguration, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		privateIpAllocationMethod := data["private_ip_address_allocation"].(string)
		properties := network.FrontendIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: network.IPAllocationMethod(privateIpAllocationMethod),
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

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
			properties.Subnet = &network.Subnet{
				ID: &v,
			}
		}

		name := data["name"].(string)
		zones := azure.ExpandZones(data["zones"].([]interface{}))
		frontEndConfig := network.FrontendIPConfiguration{
			Name:                                    &name,
			FrontendIPConfigurationPropertiesFormat: &properties,
			Zones:                                   zones,
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs
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

		zones := make([]string, 0)
		if zs := config.Zones; zs != nil {
			zones = *zs
		}
		ipConfig["zones"] = zones

		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			ipConfig["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				ipConfig["subnet_id"] = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				ipConfig["private_ip_address"] = *pip
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
			ipConfig["load_balancer_rules"] = schema.NewSet(schema.HashString, loadBalancingRules)

			inboundNatRules := make([]interface{}, 0)
			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, *rule.ID)
				}

			}
			ipConfig["inbound_nat_rules"] = schema.NewSet(schema.HashString, inboundNatRules)

			outboundRules := make([]interface{}, 0)
			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					outboundRules = append(outboundRules, *rule.ID)
				}

			}
			ipConfig["outbound_rules"] = schema.NewSet(schema.HashString, outboundRules)
		}

		result = append(result, ipConfig)
	}
	return result
}
