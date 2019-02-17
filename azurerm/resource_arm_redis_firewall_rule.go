package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRedisFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRedisFirewallRuleCreateUpdate,
		Read:   resourceArmRedisFirewallRuleRead,
		Update: resourceArmRedisFirewallRuleCreateUpdate,
		Delete: resourceArmRedisFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateRedisFirewallRuleName,
			},

			"redis_cache_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"start_ip": {
				Type:     schema.TypeString,
				Required: true,
			},

			"end_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceArmRedisFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisFirewallClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Redis Firewall Rule creation.")

	name := d.Get("name").(string)
	cacheName := d.Get("redis_cache_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	startIP := d.Get("start_ip").(string)
	endIP := d.Get("end_ip").(string)

	if requireResourcesToBeImported {
		existing, err := client.Get(ctx, resourceGroup, cacheName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Redis Firewall Rule %q (cache %q / resource group %q) ID", name, cacheName, resourceGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_redis_firewall_rule", *existing.ID)
		}
	}

	parameters := redis.FirewallRuleCreateParameters{
		FirewallRuleProperties: &redis.FirewallRuleProperties{
			StartIP: utils.String(startIP),
			EndIP:   utils.String(endIP),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, cacheName, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, cacheName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Redis Firewall Rule %q (cache %q / resource group %q) ID", name, cacheName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmRedisFirewallRuleRead(d, meta)
}

func resourceArmRedisFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisFirewallClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	cacheName := id.Path["Redis"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(ctx, resourceGroup, cacheName, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Redis Firewall Rule %q was not found in Cache %q / Resource Group %q - removing from state", name, cacheName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Redis Firewall Rule %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("redis_cache_name", cacheName)
	d.Set("resource_group_name", resourceGroup)
	if props := resp.FirewallRuleProperties; props != nil {
		d.Set("start_ip", props.StartIP)
		d.Set("end_ip", props.EndIP)
	}

	return nil
}

func resourceArmRedisFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).redisFirewallClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	cacheName := id.Path["Redis"]
	name := id.Path["firewallRules"]

	resp, err := client.Delete(ctx, resourceGroup, cacheName, name)

	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing AzureRM delete request of Redis Firewall Rule %q (cache %q / resource group %q): %+v", name, cacheName, resourceGroup, err)
		}
	}

	return nil
}

func validateRedisFirewallRuleName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if matched := regexp.MustCompile(`^[0-9a-zA-Z_]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters", k))
	}

	return warnings, errors
}
