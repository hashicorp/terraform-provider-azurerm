package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	logicParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
	logicValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const OperationInsightsRPName = "Microsoft.OperationalInsights"

func resourceArmSentinelAlertRuleAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSentinelAlertRuleActionCreate,
		Read:   resourceArmSentinelAlertRuleActionRead,
		Delete: resourceArmSentinelAlertRuleActionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.SentinelAlertRuleActionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"rule_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SentinelAlertRuleID,
			},

			"logic_app_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: logicValidate.LogicAppWorkflowID,
			},

			"logic_app_trigger_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmSentinelAlertRuleActionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	logicTriggerClient := meta.(*clients.Client).Logic.WorkflowTriggersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	ruleID, err := parse.SentinelAlertRuleID(d.Get("rule_id").(string))
	if err != nil {
		return err
	}
	lappId, err := logicParse.LogicAppWorkflowID(d.Get("logic_app_id").(string))
	if err != nil {
		return err
	}

	// Ensure no existed resources
	resp, err := client.GetAction(ctx, ruleID.ResourceGroup, OperationInsightsRPName, ruleID.Workspace, ruleID.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q): %+v", name, ruleID.ResourceGroup, ruleID.Workspace, ruleID.Name, err)
		}
	}

	if resp.ID != nil && *resp.ID != "" {
		return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_action", *resp.ID)
	}

	// List callback URL for sentinel alert specific trigger from the workspace containing specified alert rule.
	triggerName := d.Get("logic_app_trigger_name").(string)
	tresp, err := logicTriggerClient.ListCallbackURL(ctx, ruleID.ResourceGroup, lappId.Name, triggerName)
	if err != nil {
		return fmt.Errorf("listing callback URL for Logic App Trigger %q (Resource Group %q / Workspace %q): %+v", triggerName, ruleID.ResourceGroup, ruleID.Workspace, err)
	}

	param := securityinsight.ActionRequest{
		ActionRequestProperties: &securityinsight.ActionRequestProperties{
			TriggerURI:         tresp.Value,
			LogicAppResourceID: utils.String(lappId.String()),
		},
	}

	if _, err := client.CreateOrUpdateAction(ctx, ruleID.ResourceGroup, OperationInsightsRPName, ruleID.Workspace, ruleID.Name, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q): %+v", name, ruleID.ResourceGroup, ruleID.Workspace, ruleID.Name, err)
	}

	resp, err = client.GetAction(ctx, ruleID.ResourceGroup, OperationInsightsRPName, ruleID.Workspace, ruleID.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q): %+v", name, ruleID.ResourceGroup, ruleID.Workspace, ruleID.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q)", name, ruleID.ResourceGroup, ruleID.Workspace, ruleID.Name)
	}
	d.SetId(*resp.ID)

	return resourceArmSentinelAlertRuleActionRead(d, meta)
}

func resourceArmSentinelAlertRuleActionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleActionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAction(ctx, id.ResourceGroup, OperationInsightsRPName, id.Workspace, id.Rule, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Action %q was not found in Rule %q in Workspace %q in Resource Group %q - removing from state!", id.Name, id.Rule, id.Workspace, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q): %+v", id.Name, id.ResourceGroup, id.Workspace, id.Rule, err)
	}

	d.Set("name", id.Name)
	d.Set("rule_id", id.FormatSentinelAlertRuleId().String())
	if prop := resp.ActionResponseProperties; prop != nil {
		d.Set("logic_app_id", prop.LogicAppResourceID)
		// TODO: Uncomment below line once https://github.com/Azure/azure-rest-api-specs/issues/9424 is addressed.
		//       Also, remove the ignore import step in acctest for `logic_app_trigger_name`.
		//d.Set("logic_app_trigger_name", prop.TriggerUrl)
	}

	return nil
}

func resourceArmSentinelAlertRuleActionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SentinelAlertRuleActionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.DeleteAction(ctx, id.ResourceGroup, OperationInsightsRPName, id.Workspace, id.Rule, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Action %q (Resource Group %q / Workspace %q / Rule %q): %+v", id.Name, id.ResourceGroup, id.Workspace, id.Rule, err)
	}

	return nil
}
