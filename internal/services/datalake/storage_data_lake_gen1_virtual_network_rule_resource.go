package datalake

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datalake/store/mgmt/2016-11-01/account"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataLakeStoreVirtualNetworkRule() *pluginsdk.Resource {
	return _resourceStorageDataLakeGen1VirtualNetworkRule(true)
}

func resourceStorageDataLakeGen1VirtualNetworkRule() *pluginsdk.Resource {
	return _resourceStorageDataLakeGen1VirtualNetworkRule(false)
}

func _resourceStorageDataLakeGen1VirtualNetworkRule(showDeprecationMessage bool) *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen1VirtualNetworkRuleCreateUpdate,
		Read:   resourceStorageDataLakeGen1VirtualNetworkRuleRead,
		Update: resourceStorageDataLakeGen1VirtualNetworkRuleCreateUpdate,
		Delete: resourceStorageDataLakeGen1VirtualNetworkRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkRuleID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateStorageDataLakeGen1VirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName(),
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: networkValidate.SubnetID,
			},
		},
	}

	if showDeprecationMessage {
		resource.DeprecationMessage = "This resrouces has been renamed to `azurerm_storage_data_lake_gen1_virtual_network_rule` and it will be removed in version 3.0, you can follow the renaming guide https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/migrating-between-renamed-resources#migrating-to-a-renamed-resource "

	}

	return resource
}

func resourceStorageDataLakeGen1VirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Datalake.VirtualNetworkRulesClient
	subscriptionId := meta.(*clients.Client).Datalake.VirtualNetworkRulesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualNetworkRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	virtualNetworkSubnetId := d.Get("subnet_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Storage Data Lake Gen1 Virtual Network Rule %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_storage_data_lake_gen1_virtual_network_rule", *existing.ID)
		}
	}

	parameters := account.CreateOrUpdateVirtualNetworkRuleParameters{
		CreateOrUpdateVirtualNetworkRuleProperties: &account.CreateOrUpdateVirtualNetworkRuleProperties{
			SubnetID: utils.String(virtualNetworkSubnetId),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AccountName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating Storage Data Lake Gen1 Virtual Network Rule %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageDataLakeGen1VirtualNetworkRuleRead(d, meta)
}

func resourceStorageDataLakeGen1VirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] Storage Data Lake Gen1 Virtual Network Rule %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Network Rule %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("account_name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.SubnetID)
	}

	return nil
}

func resourceStorageDataLakeGen1VirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("deleting Storage Data Lake Gen1 Virtual Network Rule %s: %+v", id, err)
	}

	return nil
}
