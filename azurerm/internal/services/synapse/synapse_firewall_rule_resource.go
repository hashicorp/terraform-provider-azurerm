package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/2019-06-01-preview/synapse"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/synapse/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSynapseFirewallRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSynapseFirewallRuleCreateUpdate,
		Read:   resourceArmSynapseFirewallRuleRead,
		Update: resourceArmSynapseFirewallRuleCreateUpdate,
		Delete: resourceArmSynapseFirewallRuleDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SynapseFirewallRuleID(id)
			return err
		}),

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
				ValidateFunc: validate.SynapseFirewallRuleName,
			},

			"synapse_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SynapseWorkspaceID,
			},

			"start_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},

			"end_ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceArmSynapseFirewallRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, _ := parse.SynapseWorkspaceID(d.Get("synapse_workspace_id").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", name, workspaceId.ResourceGroup, workspaceId.Name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_synapse_firewall_rule", *existing.ID)
		}
	}

	parameters := synapse.IPFirewallRuleInfo{
		IPFirewallRuleProperties: &synapse.IPFirewallRuleProperties{
			StartIPAddress: utils.String(d.Get("start_ip_address").(string)),
			EndIPAddress:   utils.String(d.Get("end_ip_address").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, workspaceId.ResourceGroup, workspaceId.Name, name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", name, workspaceId.ResourceGroup, workspaceId.Name, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation/updation for Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", name, workspaceId.ResourceGroup, workspaceId.Name, err)
	}

	resp, err := client.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", name, workspaceId.ResourceGroup, workspaceId.Name, err)
	}

	d.SetId(*resp.ID)

	return resourceArmSynapseFirewallRuleRead(d, meta)
}

func resourceArmSynapseFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Synapse Firewall Rule %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", id.Name, id.Workspace.ResourceGroup, id.Workspace.Name, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", id.Workspace.String())
	if resp.IPFirewallRuleProperties != nil {
		d.Set("start_ip_address", resp.IPFirewallRuleProperties.StartIPAddress)
		d.Set("end_ip_address", resp.IPFirewallRuleProperties.EndIPAddress)
	}

	return nil
}

func resourceArmSynapseFirewallRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Synapse.FirewallRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SynapseFirewallRuleID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", id.Name, id.Workspace.ResourceGroup, id.Workspace.Name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting Synapse Firewall Rule %q (Resource Group %q / Workspace %q): %+v", id.Name, id.Workspace.ResourceGroup, id.Workspace.Name, err)
	}

	return nil
}
