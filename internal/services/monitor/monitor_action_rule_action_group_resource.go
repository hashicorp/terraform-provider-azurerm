// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-06-01-preview/alertsmanagement" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

		DeprecationMessage: `This resource has been deprecated in favour of the 'azurerm_monitor_alert_processing_rule_action_group' resource and will be removed in v4.0 of the AzureRM Provider`,

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

			"resource_group_name": commonschema.ResourceGroupName(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewActionRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetByName(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_monitor_action_rule_action_group", id.ID())
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

	if _, err := client.CreateUpdate(ctx, id.ResourceGroup, id.Name, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
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
		return fmt.Errorf("retrieving Monitor %s: %+v", *id, err)
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
		return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
	}
	return nil
}
