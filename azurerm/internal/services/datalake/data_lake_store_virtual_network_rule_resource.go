package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataLakeStoreVirtualNetworkRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataLakeStoreVirtualNetworkRuleCreateUpdate,
		Read:   resourceDataLakeStoreVirtualNetworkRuleRead,
		Update: resourceDataLakeStoreVirtualNetworkRuleCreateUpdate,
		Delete: resourceDataLakeStoreVirtualNetworkRuleDelete,

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
				ValidateFunc: validate.ValidateDataLakeStoreVirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDataLakeStoreVirtualNetworkRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountName := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	virtualNetworkSubnetId := d.Get("subnet_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Lake Store Virtual Network Rule %q (Account: %q, Resource Group: %q): %+v", name, accountName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_lake_store_virtual_network_rule", *existing.ID)
		}
	}

	parameters := account.CreateOrUpdateVirtualNetworkRuleParameters{
		CreateOrUpdateVirtualNetworkRuleProperties: &account.CreateOrUpdateVirtualNetworkRuleProperties{
			SubnetID: utils.String(virtualNetworkSubnetId),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Data Lake Store Virtual Network Rule %q (Account: %q, Resource Group: %q): %+v", name, accountName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Data Lake Store Virtual Network Rule %q (Account: %q, Resource Group: %q): %+v", name, accountName, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read Data Lake Store Virtual Network Rule %q (Account: %q, Resource Group: %q)", name, accountName, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceDataLakeStoreVirtualNetworkRuleRead(d, meta)
}

func resourceDataLakeStoreVirtualNetworkRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Data Lake Store Virtual Network Rule %q (Account: %q / Resource Group %q) was not found - removing from state", id.Name, id.AccountName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Network Rule %q (Account: %q / Resource Group: %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("account_name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.SubnetID)
	}

	return nil
}

func resourceDataLakeStoreVirtualNetworkRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("deleting Data Lake Store Virtual Network Rule %q (Account: %q / Resource Group: %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	return nil
}
