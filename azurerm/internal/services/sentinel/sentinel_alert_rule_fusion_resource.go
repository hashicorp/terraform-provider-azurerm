package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSentinelAlertRuleFusion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleFusionCreateUpdate,
		Read:   resourceSentinelAlertRuleFusionRead,
		Update: resourceSentinelAlertRuleFusionCreateUpdate,
		Delete: resourceSentinelAlertRuleFusionDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindFusion)),

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
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelAlertRuleFusionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroup, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Fusion %q: %+v", id, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_fusion", *id)
		}
	}

	params := securityinsight.FusionAlertRule{
		Kind: securityinsight.KindBasicAlertRuleKindFusion,
		FusionAlertRuleProperties: &securityinsight.FusionAlertRuleProperties{
			AlertRuleTemplateName: utils.String(d.Get("alert_rule_template_guid").(string)),
			Enabled:               utils.Bool(d.Get("enabled").(bool)),
		},
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindFusion); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
		params.Etag = resp.Value.(securityinsight.FusionAlertRule).Etag
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroup, OperationalInsightsResourceProvider, workspaceID.WorkspaceName, name, params); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleFusionRead(d, meta)
}

func resourceSentinelAlertRuleFusionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Fusion %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindFusion); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := resp.Value.(securityinsight.FusionAlertRule)

	d.Set("name", id.Name)

	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if prop := rule.FusionAlertRuleProperties; prop != nil {
		d.Set("enabled", prop.Enabled)
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
	}

	return nil
}

func resourceSentinelAlertRuleFusionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, OperationalInsightsResourceProvider, id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Fusion %q: %+v", id, err)
	}

	return nil
}
