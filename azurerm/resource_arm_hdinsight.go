package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

//hdinsight/mgmt/2015-03-01-preview/hdinsight

func resourceArmHDInsight() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHDInsightCreate,
		Read:   resourceArmHDInsightRead,
		Update: resourceArmHDInsightCreate,
		Delete: resourceArmHDInsightDelete,
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

			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "~3.6",
				ValidateFunc: validation.StringInSlice([]string{
					"~3.6",
					"3.5",
					"3.4",
					"3.3",
					"3.2",
					"3.1",
					"3.0",
				}, false),
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Linux",
				ValidateFunc: validation.StringInSlice([]string{
					"Windows",
					"Linux",
				}, false),
			},

			"tier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Standard",
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, false),
			},

			"cluster_definition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"component_version": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
					},
				},
			},

			"security_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"directory_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"organizational_unit_dn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"lads_urls": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"domain_username": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain_password": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_user_group_dns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_accounts": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"is_default": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"container": {
										Type:     schema.TypeString,
										Required: true,
									},
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmHDInsightCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).hdInsightClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM LoadBalancer creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	sku := network.LoadBalancerSku{
		Name: network.LoadBalancerSkuName(d.Get("sku").(string)),
	}
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
		Sku:      &sku,
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

	if sku := loadBalancer.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := loadBalancer.LoadBalancerPropertiesFormat; props != nil {
		if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
			d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs))

			privateIpAddress := ""
			privateIpAddresses := make([]string, 0, len(*feipConfigs))
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
