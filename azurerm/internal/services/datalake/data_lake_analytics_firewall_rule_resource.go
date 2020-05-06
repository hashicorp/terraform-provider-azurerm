package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/analytics/mgmt/2016-11-01/account"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataLakeAnalyticsFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDateLakeAnalyticsFirewallRuleCreateUpdate,
		Read:   resourceArmDateLakeAnalyticsFirewallRuleRead,
		Update: resourceArmDateLakeAnalyticsFirewallRuleCreateUpdate,
		Delete: resourceArmDateLakeAnalyticsFirewallRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateDataLakeFirewallRuleName(),
			},

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateDataLakeAccountName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"start_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IPv4Address,
			},

			"end_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.IPv4Address,
			},
		},
	}
}

func resourceArmDateLakeAnalyticsFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsFirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Lake Analytics Firewall Rule %q (Account %q / Resource Group %q): %s", name, accountName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_lake_analytics_firewall_rule", *existing.ID)
		}
	}

	startIPAddress := d.Get("start_ip_address").(string)
	endIPAddress := d.Get("end_ip_address").(string)

	log.Printf("[INFO] preparing arguments for Date Lake Analytics Firewall Rule creation %q (Account %q / Resource Group %q)", name, accountName, resourceGroup)

	dateLakeStore := account.CreateOrUpdateFirewallRuleParameters{
		CreateOrUpdateFirewallRuleProperties: &account.CreateOrUpdateFirewallRuleProperties{
			StartIPAddress: utils.String(startIPAddress),
			EndIPAddress:   utils.String(endIPAddress),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, dateLakeStore); err != nil {
		return fmt.Errorf("Error issuing create request for Data Lake Analytics %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Lake Analytics Firewall Rule %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Lake Analytics Firewall Rule %q (Account %q / Resource Group %q) ID", name, accountName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDateLakeAnalyticsFirewallRuleRead(d, meta)
}

func resourceArmDateLakeAnalyticsFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsFirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	accountName := id.Path["accounts"]
	name := id.Path["firewallRules"]

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Data Lake Analytics Firewall Rule %q was not found (Account %q / Resource Group %q)", name, accountName, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Data Lake Analytics Firewall Rule %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.FirewallRuleProperties; props != nil {
		d.Set("start_ip_address", props.StartIPAddress)
		d.Set("end_ip_address", props.EndIPAddress)
	}

	return nil
}

func resourceArmDateLakeAnalyticsFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.AnalyticsFirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	accountName := id.Path["accounts"]
	name := id.Path["firewallRules"]

	resp, err := client.Delete(ctx, resourceGroup, accountName, name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing delete request for Data Lake Analytics Firewall Rule %q (Account %q / Resource Group %q): %+v", name, accountName, resourceGroup, err)
	}

	return nil
}
