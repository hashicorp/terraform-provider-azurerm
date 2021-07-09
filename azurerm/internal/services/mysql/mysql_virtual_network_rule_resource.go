package mysql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMySQLVirtualNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySQLVirtualNetworkRuleCreateUpdate,
		Read:   resourceMySQLVirtualNetworkRuleRead,
		Update: resourceMySQLVirtualNetworkRuleCreateUpdate,
		Delete: resourceMySQLVirtualNetworkRuleDelete,

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
				ValidateFunc: networkValidate.VirtualNetworkRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceMySQLVirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewVirtualNetworkRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SubscriptionId, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mysql_virtual_network_rule", id.ID())
		}
	}

	parameters := mysql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &mysql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(d.Get("subnet_id").(string)),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(false),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   mySQLVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, id),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
	}
	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created or updated: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMySQLVirtualNetworkRuleRead(d, meta)
}

func resourceMySQLVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("server_name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("subnet_id", props.VirtualNetworkSubnetID)
	}

	return nil
}

func resourceMySQLVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

func mySQLVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *mysql.VirtualNetworkRulesClient, id parse.VirtualNetworkRuleId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving %s returned 404.", id)
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the %s: %+v", id, err)
		}

		if props := resp.VirtualNetworkRuleProperties; props != nil {
			log.Printf("[DEBUG] Retrieving %s returned Status %s", id, props.State)
			return resp, string(props.State), nil
		}

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving %s returned empty VirtualNetworkRuleProperties", id)
		return resp, "Unknown", nil
	}
}
