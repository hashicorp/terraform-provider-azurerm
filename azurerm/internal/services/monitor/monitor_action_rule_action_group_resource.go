package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorActionRuleActionGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorActionRuleActionGroupCreateUpdate,
		Read:   resourceMonitorActionRuleActionGroupRead,
		Update: resourceMonitorActionRuleActionGroupCreateUpdate,
		Delete: resourceMonitorActionRuleActionGroupDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.ActionRuleID(id)
			return err
		}, importMonitorActionRule(alertsmanagement.TypeActionGroup)),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ActionRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"action_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ActionGroupID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"condition": schemaActionRuleConditions(),

			"scope": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(alertsmanagement.ScopeTypeResourceGroup),
								string(alertsmanagement.ScopeTypeResource),
							}, false),
						},

						"resource_ids": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorActionRuleActionGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.GetByName(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_action_rule_action_group", *existing.ID)
		}
	}

	actionRuleStatus := alertsmanagement.Enabled
	if !d.Get("enabled").(bool) {
		actionRuleStatus = alertsmanagement.Disabled
	}

	actionRule := alertsmanagement.ActionRule{
		// the location is always global from the portal
		Location: utils.String(location.Normalize("Global")),
		Properties: &alertsmanagement.ActionGroup{
			ActionGroupID: utils.String(d.Get("action_group_id").(string)),
			Scope:         expandActionRuleScope(d.Get("scope").([]interface{})),
			Conditions:    expandActionRuleConditions(d.Get("condition").([]interface{})),
			Description:   utils.String(d.Get("description").(string)),
			Status:        actionRuleStatus,
			Type:          alertsmanagement.TypeActionGroup,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateUpdate(ctx, resourceGroup, name, actionRule); err != nil {
		return fmt.Errorf("creating/updatinge Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetByName(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceMonitorActionRuleActionGroupRead(d, meta)
}

func resourceMonitorActionRuleActionGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ActionRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Action Rule %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Monitor ActionRule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if resp.Properties != nil {
		props, _ := resp.Properties.AsActionGroup()
		d.Set("description", props.Description)
		d.Set("action_group_id", props.ActionGroupID)
		d.Set("enabled", props.Status == alertsmanagement.Enabled)
		if err := d.Set("scope", flattenActionRuleScope(props.Scope)); err != nil {
			return fmt.Errorf("setting scope: %+v", err)
		}
		if err := d.Set("condition", flattenActionRuleConditions(props.Conditions)); err != nil {
			return fmt.Errorf("setting condition: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMonitorActionRuleActionGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ActionRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting monitor ActionRule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}
