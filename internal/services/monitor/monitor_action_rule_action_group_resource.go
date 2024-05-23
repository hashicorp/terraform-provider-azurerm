// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/actionrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
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
			_, err := actionrules.ParseActionRuleID(id)
			return err
		}, importMonitorActionRule(actionrules.ActionRuleTypeActionGroup)),

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
								string(actionrules.ScopeTypeResourceGroup),
								string(actionrules.ScopeTypeResource),
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMonitorActionRuleActionGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := actionrules.NewActionRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetByName(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_action_rule_action_group", id.ID())
		}
	}

	actionRuleStatus := actionrules.ActionRuleStatusEnabled
	if !d.Get("enabled").(bool) {
		actionRuleStatus = actionrules.ActionRuleStatusDisabled
	}

	actionRule := actionrules.ActionRule{
		// the location is always global from the portal
		Location: location.Normalize("Global"),
		Properties: &actionrules.ActionGroup{
			ActionGroupId: utils.String(d.Get("action_group_id").(string)),
			Scope:         expandActionRuleScope(d.Get("scope").([]interface{})),
			Conditions:    expandActionRuleConditions(d.Get("condition").([]interface{})),
			Description:   utils.String(d.Get("description").(string)),
			Status:        pointer.To(actionRuleStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateUpdate(ctx, id, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorActionRuleActionGroupRead(d, meta)
}

func resourceMonitorActionRuleActionGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actionrules.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByName(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ActionRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		props, ok := model.Properties.(actionrules.ActionGroup)
		if !ok {
			return fmt.Errorf("%s is not of type `ActionGroup`", id)
		}
		d.Set("description", props.Description)
		d.Set("action_group_id", props.ActionGroupId)
		if props.Status != nil {
			d.Set("enabled", *props.Status == actionrules.ActionRuleStatusEnabled)
		}
		if err := d.Set("scope", flattenActionRuleScope(props.Scope)); err != nil {
			return fmt.Errorf("setting scope: %+v", err)
		}
		if err := d.Set("condition", flattenActionRuleConditions(props.Conditions)); err != nil {
			return fmt.Errorf("setting condition: %+v", err)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceMonitorActionRuleActionGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actionrules.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}
