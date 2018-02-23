package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var networkSecurityGroupResourceName = "azurerm_network_security_group"

func resourceArmNetworkSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkSecurityGroupCreate,
		Read:   resourceArmNetworkSecurityGroupRead,
		Update: resourceArmNetworkSecurityGroupCreate,
		Delete: resourceArmNetworkSecurityGroupDelete,
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

			"security_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateStringLength(140),
						},

						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateNetworkSecurityRuleProtocol,
							StateFunc:        ignoreCaseStateFunc,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"source_port_range": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"security_rule.source_port_ranges"},
						},

						"source_port_ranges": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"security_rule.source_port_range"},
						},

						"destination_port_range": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"security_rule.destination_port_ranges"},
						},

						"destination_port_ranges": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"security_rule.destination_port_range"},
						},

						"source_address_prefix": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"security_rule.source_address_prefixes"},
						},

						"source_address_prefixes": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"security_rule.source_address_prefix"},
						},

						"destination_address_prefix": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"security_rule.destination_address_prefixes"},
						},

						"destination_address_prefixes": {
							Type:          schema.TypeSet,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Set:           schema.HashString,
							ConflictsWith: []string{"security_rule.destination_address_prefix"},
						},

						"access": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.SecurityRuleAccessAllow),
								string(network.SecurityRuleAccessDeny),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"priority": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 4096),
						},

						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.SecurityRuleDirectionInbound),
								string(network.SecurityRuleDirectionOutbound),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	sgRules, sgErr := expandAzureRmSecurityRules(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Network Security Group Rules: %+v", sgErr)
	}

	azureRMLockByName(name, networkSecurityGroupResourceName)
	defer azureRMUnlockByName(name, networkSecurityGroupResourceName)

	sg := network.SecurityGroup{
		Name:     &name,
		Location: &location,
		SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
			SecurityRules: &sgRules,
		},
		Tags: expandTags(tags),
	}

	_, createErr := client.CreateOrUpdate(resGroup, name, sg, make(chan struct{}))
	err := <-createErr
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Network %q (resource group %q) ID", name, resGroup)
	}

	log.Printf("[DEBUG] Waiting for NSG (%q) to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    networkSecurityGroupStateRefreshFunc(client, resGroup, name),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for NSG (%q) to become available: %+v", name, err)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkSecurityGroupRead(d, meta)
}

func resourceArmNetworkSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkSecurityGroups"]

	resp, err := client.Get(resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Security Group %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	if props := resp.SecurityGroupPropertiesFormat; props != nil {
		d.Set("security_rule", flattenNetworkSecurityRules(props.SecurityRules))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmNetworkSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).secGroupClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkSecurityGroups"]

	_, deleteErr := client.Delete(resGroup, name, make(chan struct{}))
	err = <-deleteErr

	return err
}

func flattenNetworkSecurityRules(rules *[]network.SecurityRule) []interface{} {
	result := make([]interface{}, 0)

	if rules != nil {
		for _, rule := range *rules {
			sgRule := make(map[string]interface{})
			sgRule["name"] = *rule.Name

			if props := rule.SecurityRulePropertiesFormat; props != nil {
				if props.DestinationAddressPrefix != nil {
					sgRule["destination_address_prefix"] = *props.DestinationAddressPrefix
				}
				if props.DestinationAddressPrefixes != nil {
					sgRule["destination_address_prefixes"] = *props.DestinationAddressPrefixes
				}
				if props.DestinationPortRange != nil {
					sgRule["destination_port_range"] = *props.DestinationPortRange
				}
				if props.DestinationPortRanges != nil {
					sgRule["destination_port_ranges"] = *props.DestinationPortRanges
				}
				if props.SourceAddressPrefix != nil {
					sgRule["source_address_prefix"] = *props.SourceAddressPrefix
				}
				if props.SourceAddressPrefixes != nil {
					sgRule["source_address_prefixes"] = *props.SourceAddressPrefixes
				}
				if props.SourcePortRange != nil {
					sgRule["source_port_range"] = *props.SourcePortRange
				}
				if props.SourcePortRanges != nil {
					sgRule["source_port_ranges"] = *props.SourcePortRanges
				}
				sgRule["priority"] = int(*props.Priority)
				sgRule["access"] = string(props.Access)
				sgRule["direction"] = string(props.Direction)
				sgRule["protocol"] = string(props.Protocol)

				if props.Description != nil {
					sgRule["description"] = *props.Description
				}
			}

			result = append(result, sgRule)
		}
	}

	return result
}

func expandAzureRmSecurityRules(d *schema.ResourceData) ([]network.SecurityRule, error) {
	sgRules := d.Get("security_rule").([]interface{})
	rules := make([]network.SecurityRule, 0)

	for _, sgRaw := range sgRules {
		data := sgRaw.(map[string]interface{})

		name := data["name"].(string)
		source_port_range := data["source_port_range"].(string)
		source_port_ranges := data["source_port_ranges"].(*schema.Set).List()

		source_port_ranges_outputs := make([]string, 0)
		for _, sourcePortRange := range source_port_ranges {
			source_port_ranges_outputs = append(source_port_ranges_outputs, sourcePortRange.(string))
		}

		destination_port_range := data["destination_port_range"].(string)
		destination_port_ranges := data["destination_port_ranges"].(*schema.Set).List()

		destination_port_ranges_outputs := make([]string, 0)
		for _, destinationPortRange := range destination_port_ranges {
			destination_port_ranges_outputs = append(destination_port_ranges_outputs, destinationPortRange.(string))
		}

		source_address_prefix := data["source_address_prefix"].(string)
		source_address_prefixes := data["source_address_prefixes"].(*schema.Set).List()

		source_address_prefixes_outputs := make([]string, 0)
		for _, sourceAddressPrefix := range source_address_prefixes {
			source_address_prefixes_outputs = append(source_address_prefixes_outputs, sourceAddressPrefix.(string))
		}

		destination_address_prefix := data["destination_address_prefix"].(string)
		destination_address_prefixes := data["destination_address_prefixes"].(*schema.Set).List()

		destination_address_prefixes_outputs := make([]string, 0)
		for _, destinationAddressPrefix := range destination_address_prefixes {
			destination_address_prefixes_outputs = append(destination_address_prefixes_outputs, destinationAddressPrefix.(string))
		}

		priority := int32(data["priority"].(int))
		access := data["access"].(string)
		direction := data["direction"].(string)
		protocol := data["protocol"].(string)

		properties := network.SecurityRulePropertiesFormat{
			SourcePortRange:            &source_port_range,
			SourcePortRanges:           &source_port_ranges_outputs,
			DestinationPortRange:       &destination_port_range,
			DestinationPortRanges:      &destination_port_ranges_outputs,
			SourceAddressPrefix:        &source_address_prefix,
			SourceAddressPrefixes:      &source_address_prefixes_outputs,
			DestinationAddressPrefix:   &destination_address_prefix,
			DestinationAddressPrefixes: &destination_address_prefixes_outputs,
			Priority:                   &priority,
			Access:                     network.SecurityRuleAccess(access),
			Direction:                  network.SecurityRuleDirection(direction),
			Protocol:                   network.SecurityRuleProtocol(protocol),
		}

		if v := data["description"].(string); v != "" {
			properties.Description = &v
		}

		rule := network.SecurityRule{
			Name: &name,
			SecurityRulePropertiesFormat: &properties,
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func networkSecurityGroupStateRefreshFunc(client network.SecurityGroupsClient, resourceGroupName string, sgName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(resourceGroupName, sgName, "")
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in networkSecurityGroupStateRefreshFunc for NSG '%s' (RG: '%s'): %+v", sgName, resourceGroupName, err)
		}

		return res, *res.SecurityGroupPropertiesFormat.ProvisioningState, nil
	}
}
