package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	eventhubParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	eventhubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	logAnalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	logAnalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/validate"
	storageParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	storageValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorAADDiagnosticSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorAADDiagnosticSettingCreateUpdate,
		Read:   resourceMonitorAADDiagnosticSettingRead,
		Update: resourceMonitorAADDiagnosticSettingCreateUpdate,
		Delete: resourceMonitorAADDiagnosticSettingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MonitorAADDiagnosticSettingID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MonitorDiagnosticSettingName,
			},

			// When absent, will use the default eventhub, whilst the Diagnostic Setting API will return this property as an empty string. Therefore, it is useless to make this property as Computed.
			"eventhub_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_authorization_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: eventhubValidate.NamespaceAuthorizationRuleID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: storageValidate.StorageAccountID,
				AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
			},

			"log": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.AuditLogs),
								string(aad.SignInLogs),
								"ADFSSignInLogs",
								"ManagedIdentitySignInLogs",
								"NonInteractiveUserSignInLogs",
								"ProvisioningLogs",
								"ServicePrincipalSignInLogs",
							}, false),
						},

						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						"retention_policy": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
										Default:      0,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceMonitorAADDiagnosticSettingCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Azure ARM AAD Diagnostic Setting.")

	name := d.Get("name").(string)
	id := parse.NewMonitorAADDiagnosticSettingID(name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_monitor_aad_diagnostic_setting", id.ID())
		}
	}

	logs := expandMonitorAADDiagnosticsSettingsLogs(d.Get("log").(*schema.Set).List())

	// If there is no `enabled` log entry, the PUT will succeed while the next GET will return a 404.
	// Therefore, ensure users has at least one enabled log entry.
	valid := false
	for _, log := range logs {
		if log.Enabled != nil && *log.Enabled {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("At least one of the `log` of the %s should be enabled", id)
	}

	properties := aad.DiagnosticSettingsResource{
		DiagnosticSettings: &aad.DiagnosticSettings{
			Logs: &logs,
		},
	}

	eventHubAuthorizationRuleId := d.Get("eventhub_authorization_rule_id").(string)
	eventHubName := d.Get("eventhub_name").(string)
	if eventHubAuthorizationRuleId != "" {
		properties.DiagnosticSettings.EventHubAuthorizationRuleID = utils.String(eventHubAuthorizationRuleId)
		properties.DiagnosticSettings.EventHubName = utils.String(eventHubName)
	}

	workspaceId := d.Get("log_analytics_workspace_id").(string)
	if workspaceId != "" {
		properties.DiagnosticSettings.WorkspaceID = utils.String(workspaceId)
	}

	storageAccountId := d.Get("storage_account_id").(string)
	if storageAccountId != "" {
		properties.DiagnosticSettings.StorageAccountID = utils.String(storageAccountId)
	}

	if _, err := client.CreateOrUpdate(ctx, properties, name); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorAADDiagnosticSettingRead(d, meta)
}

func resourceMonitorAADDiagnosticSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MonitorAADDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)

	d.Set("eventhub_name", resp.EventHubName)
	eventhubAuthorizationRuleId := ""
	if resp.EventHubAuthorizationRuleID != nil && *resp.EventHubAuthorizationRuleID != "" {
		parsedId, err := eventhubParse.NamespaceAuthorizationRuleIDInsensitively(*resp.EventHubAuthorizationRuleID)
		if err != nil {
			return err
		}

		eventhubAuthorizationRuleId = parsedId.ID()
	}
	d.Set("eventhub_authorization_rule_id", eventhubAuthorizationRuleId)

	workspaceId := ""
	if resp.WorkspaceID != nil && *resp.WorkspaceID != "" {
		parsedId, err := logAnalyticsParse.LogAnalyticsWorkspaceID(*resp.WorkspaceID)
		if err != nil {
			return err
		}

		workspaceId = parsedId.ID()
	}
	d.Set("log_analytics_workspace_id", workspaceId)

	storageAccountId := ""
	if resp.StorageAccountID != nil && *resp.StorageAccountID != "" {
		parsedId, err := storageParse.StorageAccountID(*resp.StorageAccountID)
		if err != nil {
			return err
		}

		storageAccountId = parsedId.ID()
	}
	d.Set("storage_account_id", storageAccountId)

	if err := d.Set("log", flattenMonitorAADDiagnosticLogs(resp.Logs)); err != nil {
		return fmt.Errorf("setting `log`: %+v", err)
	}

	return nil
}

func resourceMonitorAADDiagnosticSettingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AADDiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MonitorAADDiagnosticSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// API appears to be eventually consistent (identified during tainting this resource)
	log.Printf("[DEBUG] Waiting for %s to disappear", id)
	timeout, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   monitorAADDiagnosticSettingDeletedRefreshFunc(ctx, client, id.Name),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	return nil
}

func monitorAADDiagnosticSettingDeletedRefreshFunc(ctx context.Context, client *aad.DiagnosticSettingsClient, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}
			return nil, "", fmt.Errorf("issuing read request in monitorAADDiagnosticSettingDeletedRefreshFunc: %s", err)
		}

		return res, "Exists", nil
	}
}

func expandMonitorAADDiagnosticsSettingsLogs(input []interface{}) []aad.LogSettings {
	results := make([]aad.LogSettings, 0)

	for _, raw := range input {
		v := raw.(map[string]interface{})

		category := v["category"].(string)
		enabled := v["enabled"].(bool)

		policyRaw := v["retention_policy"].([]interface{})[0].(map[string]interface{})
		retentionDays := policyRaw["days"].(int)
		retentionEnabled := policyRaw["enabled"].(bool)

		output := aad.LogSettings{
			Category: aad.Category(category),
			Enabled:  utils.Bool(enabled),
			RetentionPolicy: &aad.RetentionPolicy{
				Days:    utils.Int32(int32(retentionDays)),
				Enabled: utils.Bool(retentionEnabled),
			},
		}

		results = append(results, output)
	}

	return results
}

func flattenMonitorAADDiagnosticLogs(input *[]aad.LogSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		category := string(v.Category)

		enabled := false
		if v.Enabled != nil {
			enabled = *v.Enabled
		}

		policies := make([]interface{}, 0)
		if inputPolicy := v.RetentionPolicy; inputPolicy != nil {
			days := 0
			if inputPolicy.Days != nil {
				days = int(*inputPolicy.Days)
			}

			enabled := false
			if inputPolicy.Enabled != nil {
				enabled = *inputPolicy.Enabled
			}

			policies = append(policies, map[string]interface{}{
				"days":    days,
				"enabled": enabled,
			})
		}

		results = append(results, map[string]interface{}{
			"category":         category,
			"enabled":          enabled,
			"retention_policy": policies,
		})
	}

	return results
}
