package mssql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlVirtualNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlVirtualNetworkRuleCreateUpdate,
		Read:   resourceMsSqlVirtualNetworkRuleRead,
		Update: resourceMsSqlVirtualNetworkRuleCreateUpdate,
		Delete: resourceMsSqlVirtualNetworkRuleDelete,

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
				ValidateFunc: validate.VirtualNetworkRuleName,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerID,
			},

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"ignore_missing_vnet_service_endpoint": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false, // When not provided, Azure defaults to false
			},
		},
	}
}

func resourceMsSqlVirtualNetworkRuleCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID %q: %+v", d.Get("server_id"), err)
	}

	id := parse.NewVirtualNetworkRuleID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, d.Get("name").(string))

	subnetId, err := networkParse.SubnetID(d.Get("subnet_id").(string))
	if err != nil {
		return fmt.Errorf("parsing subnet ID %q: %+v", d.Get("subnet_id"), err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing MSSQL %s: %+v", id.String(), err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_virtual_network_rule", id.ID())
		}
	}

	parameters := sql.VirtualNetworkRule{
		VirtualNetworkRuleProperties: &sql.VirtualNetworkRuleProperties{
			VirtualNetworkSubnetID:           utils.String(subnetId.ID()),
			IgnoreMissingVnetServiceEndpoint: utils.Bool(d.Get("ignore_missing_vnet_service_endpoint").(bool)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, parameters); err != nil {
		return fmt.Errorf("creating MSSQL %s: %+v", id.String(), err)
	}

	// Wait for the provisioning state to become ready
	log.Printf("[DEBUG] Waiting for MSSQL %s to become ready", id.String())
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Initializing", "InProgress", "Unknown", "ResponseNotFound"},
		Target:                    []string{"Ready"},
		Refresh:                   mssqlVirtualNetworkStateStatusCodeRefreshFunc(ctx, client, id),
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for MSSQL %s to be created or updated: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceMsSqlVirtualNetworkRuleRead(d, meta)
}

func resourceMsSqlVirtualNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] MSSQL %s was not found - removing from state", id.String())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving MSSQL %s: %+v", id.String(), err)
	}

	d.Set("name", id.Name)

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)
	d.Set("server_id", serverId.ID())

	if props := resp.VirtualNetworkRuleProperties; props != nil {
		d.Set("ignore_missing_vnet_service_endpoint", props.IgnoreMissingVnetServiceEndpoint)

		subnetId := ""
		if sid := props.VirtualNetworkSubnetID; sid != nil {
			id, err := networkParse.SubnetID(*props.VirtualNetworkSubnetID)
			if err != nil {
				return fmt.Errorf("parsing subnet ID returned by API %q: %+v", *sid, err)
			}
			subnetId = id.ID()
		}
		d.Set("subnet_id", subnetId)
	}

	return nil
}

func resourceMsSqlVirtualNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.VirtualNetworkRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkRuleID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing ID %q: %+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("deleting MSSQL %s: %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("waiting for deletion of MSSQL %s: %+v", id.String(), err)
	}

	return nil
}

// mssqlVirtualNetworkStateStatusCodeRefreshFunc refreshes and checks the state of the SQL Virtual Network Rule.
// Response will contain a VirtualNetworkRuleProperties struct with a State property.
// The state property contain one of the following states (except ResponseNotFound).
//	* Deleting
//	* Initializing
//	* InProgress
//	* Unknown
//	* Ready
//	* ResponseNotFound (custom state in case of 404)
func mssqlVirtualNetworkStateStatusCodeRefreshFunc(ctx context.Context, client *sql.VirtualNetworkRulesClient, id parse.VirtualNetworkRuleId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving MSSQL %s", id.String())
				return nil, "ResponseNotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of MSSQL %s: %+v", id.String(), err)
		}

		if props := resp.VirtualNetworkRuleProperties; props != nil {
			log.Printf("[DEBUG] Retrieving MSSQL %s returned Status %s", id.String(), props.State)
			return resp, string(props.State), nil
		}

		// Valid response was returned but VirtualNetworkRuleProperties was nil. Basically the rule exists, but with no properties for some reason. Assume Unknown instead of returning error.
		log.Printf("[DEBUG] Retrieving MSSQL %s returned empty VirtualNetworkRuleProperties", id.String())
		return resp, "Unknown", nil
	}
}
