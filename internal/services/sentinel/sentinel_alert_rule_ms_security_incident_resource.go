// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSentinelAlertRuleMsSecurityIncident() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Read:   resourceSentinelAlertRuleMsSecurityIncidentRead,
		Update: resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Delete: resourceSentinelAlertRuleMsSecurityIncidentDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := alertrules.ParseAlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation)),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: alertrules.ValidateWorkspaceID,
			},

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"product_filter": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(alertrules.MicrosoftSecurityProductNameMicrosoftCloudAppSecurity),
					string(alertrules.MicrosoftSecurityProductNameAzureSecurityCenter),
					string(alertrules.MicrosoftSecurityProductNameAzureActiveDirectoryIdentityProtection),
					string(alertrules.MicrosoftSecurityProductNameAzureSecurityCenterForIoT),
					string(alertrules.MicrosoftSecurityProductNameAzureAdvancedThreatProtection),
					string(alertrules.MicrosoftSecurityProductNameMicrosoftDefenderAdvancedThreatProtection),
					string(alertrules.MicrosoftSecurityProductNameOfficeThreeSixFiveAdvancedThreatProtection),
				}, false),
			},

			"severity_filter": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(alertrules.AlertSeverityHigh),
						string(alertrules.AlertSeverityMedium),
						string(alertrules.AlertSeverityLow),
						string(alertrules.AlertSeverityInformational),
					}, false),
				},
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"display_name_filter": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true, // remove in 3.0
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"display_name_exclude_filter": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_ms_security_incident", id.ID())
		}
	}

	param := alertrules.MicrosoftSecurityIncidentCreationAlertRule{
		Properties: &alertrules.MicrosoftSecurityIncidentCreationAlertRuleProperties{
			ProductFilter:    alertrules.MicrosoftSecurityProductName(d.Get("product_filter").(string)),
			DisplayName:      d.Get("display_name").(string),
			Description:      utils.String(d.Get("description").(string)),
			Enabled:          d.Get("enabled").(bool),
			SeveritiesFilter: expandAlertRuleMsSecurityIncidentSeverityFilter(d.Get("severity_filter").(*pluginsdk.Set).List()),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.Properties.AlertRuleTemplateName = utils.String(v.(string))
	}

	if dnf, ok := d.GetOk("display_name_filter"); ok {
		param.Properties.DisplayNamesFilter = utils.ExpandStringSlice(dnf.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("display_name_exclude_filter"); ok {
		param.Properties.DisplayNamesExcludeFilter = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleMsSecurityIncidentRead(d, meta)
}

func resourceSentinelAlertRuleMsSecurityIncidentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Sentinel Alert Rule Ms Security Incident %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if err = assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}

		modelPtr := *model
		rule := modelPtr.(alertrules.MicrosoftSecurityIncidentCreationAlertRule)

		d.Set("name", id.RuleId)

		workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
		d.Set("log_analytics_workspace_id", workspaceId.ID())
		if prop := rule.Properties; prop != nil {
			d.Set("product_filter", string(prop.ProductFilter))
			d.Set("display_name", prop.DisplayName)
			d.Set("description", prop.Description)
			d.Set("enabled", prop.Enabled)
			d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)

			if err := d.Set("display_name_filter", utils.FlattenStringSlice(prop.DisplayNamesFilter)); err != nil {
				return fmt.Errorf(`setting "display_name_filter": %+v`, err)
			}
			if err := d.Set("display_name_exclude_filter", utils.FlattenStringSlice(prop.DisplayNamesExcludeFilter)); err != nil {
				return fmt.Errorf(`setting "display_name_exclude_filter": %+v`, err)
			}
			if err := d.Set("severity_filter", flattenAlertRuleMsSecurityIncidentSeverityFilter(prop.SeveritiesFilter)); err != nil {
				return fmt.Errorf(`setting "severity_filter": %+v`, err)
			}
		}
	}

	return nil
}

func resourceSentinelAlertRuleMsSecurityIncidentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	return nil
}

func expandAlertRuleMsSecurityIncidentSeverityFilter(input []interface{}) *[]alertrules.AlertSeverity {
	result := make([]alertrules.AlertSeverity, 0)

	for _, e := range input {
		result = append(result, alertrules.AlertSeverity(e.(string)))
	}

	return &result
}

func flattenAlertRuleMsSecurityIncidentSeverityFilter(input *[]alertrules.AlertSeverity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
