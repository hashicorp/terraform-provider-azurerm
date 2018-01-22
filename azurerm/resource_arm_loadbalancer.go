package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLoadBalancerCreate,
		Read:   resourecArmLoadBalancerRead,
		Update: resourceArmLoadBalancerCreate,
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"frontend_ip_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validateLoadBalancerPrivateIpAddressAllocation,
							StateFunc:        ignoreCaseStateFunc,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"load_balancer_rules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"inbound_nat_rules": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
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

			"tags": tagsSchema(),
		},
	}
}

func resourceArmLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM LoadBalancer creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	properties := network.LoadBalancerPropertiesFormat{}

	if _, ok := d.GetOk("frontend_ip_configuration"); ok {
		properties.FrontendIPConfigurations = expandAzureRmLoadBalancerFrontendIpConfigurations(d)
	}

	loadBalancer := network.LoadBalancer{
		Name:     utils.String(name),
		Location: utils.String(location),
		Tags:     expandedTags,
		LoadBalancerPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, loadBalancer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating LoadBalancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error Retrieving LoadBalancer %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read LoadBalancer %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	// TODO: is this still needed?
	log.Printf("[DEBUG] Waiting for LoadBalancer (%q) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: loadbalancerStateRefreshFunc(ctx, client, resGroup, name),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for LoadBalancer (%q - Resource Group %q) to become available: %s", name, resGroup, err)
	}

	return resourecArmLoadBalancerRead(d, meta)
}

func resourecArmLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	loadBalancer, exists, err := retrieveLoadBalancerById(d.Id(), meta)
	if err != nil {
		return fmt.Errorf("Error retrieving Load Balancer by ID %q: %+v", d.Id(), err)
	}
	if !exists {
		d.SetId("")
		log.Printf("[INFO] LoadBalancer %q not found. Removing from state", d.Id())
		return nil
	}

	d.Set("name", loadBalancer.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := loadBalancer.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs))

			privateIpAddress := ""
			privateIpAddresses := make([]string, 0, len(*feipConfigs))
			for _, config := range *feipConfigs {
				if feipProps := config.FrontendIPConfigurationPropertiesFormat; feipProps != nil {
					if privateIpAddress == "" {
						privateIpAddress = *feipProps.PrivateIPAddress
					}

					privateIpAddresses = append(privateIpAddresses, *feipProps.PrivateIPAddress)
				}
			}

			d.Set("private_ip_address", privateIpAddress)
			d.Set("private_ip_addresses", privateIpAddresses)
		}
	}

	flattenAndSetTags(d, loadBalancer.Tags)

	return nil
}

func resourceArmLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).loadBalancerClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return errwrap.Wrapf("Error Parsing Azure Resource ID {{err}}", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["loadBalancers"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Load Balancer %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
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

		if v := data["subnet_id"].(string); v != "" {
			properties.Subnet = &network.Subnet{
				ID: &v,
			}
		}

		name := data["name"].(string)
		frontEndConfig := network.FrontendIPConfiguration{
			Name: &name,
			FrontendIPConfigurationPropertiesFormat: &properties,
		}

		frontEndConfigs = append(frontEndConfigs, frontEndConfig)
	}

	return &frontEndConfigs
}

func flattenLoadBalancerFrontendIpConfiguration(ipConfigs *[]network.FrontendIPConfiguration) []interface{} {
	result := make([]interface{}, 0, len(*ipConfigs))
	for _, config := range *ipConfigs {
		ipConfig := make(map[string]interface{})
		ipConfig["name"] = *config.Name

		if props := config.FrontendIPConfigurationPropertiesFormat; props != nil {
			ipConfig["private_ip_address_allocation"] = props.PrivateIPAllocationMethod

			if subnet := props.Subnet; subnet != nil {
				ipConfig["subnet_id"] = *subnet.ID
			}

			if pip := props.PrivateIPAddress; pip != nil {
				ipConfig["private_ip_address"] = *pip
			}

			if pip := props.PublicIPAddress; pip != nil {
				ipConfig["public_ip_address_id"] = *pip.ID
			}

			if rules := props.LoadBalancingRules; rules != nil {
				loadBalancingRules := make([]string, 0, len(*rules))
				for _, rule := range *rules {
					loadBalancingRules = append(loadBalancingRules, *rule.ID)
				}

				ipConfig["load_balancer_rules"] = loadBalancingRules
			}

			if rules := props.InboundNatRules; rules != nil {
				inboundNatRules := make([]string, 0, len(*rules))
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, *rule.ID)
				}

				ipConfig["inbound_nat_rules"] = inboundNatRules

			}
		}

		result = append(result, ipConfig)
	}
	return result
}
