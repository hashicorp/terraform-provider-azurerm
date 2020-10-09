package aadmgmt

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func azureADDiagnosticSettingsResource() *schema.Resource {
	return &schema.Resource{
		Create: diagnosticsSettingsCreateOrUpdate,
		Read:   diagnosticsSettingsRead,
		Update: diagnosticsSettingsCreateOrUpdate,
		Delete: diagnosticsSettingsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateDiagnosticSettingsName,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				AtLeastOneOf: []string{"storage_account_id", "workspace_id", "event_hub_auth_rule_id", "event_hub_name"},
			},

			"workspace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"storage_account_id", "workspace_id", "event_hub_auth_rule_id", "event_hub_name"},
				ValidateFunc: azure.ValidateResourceID,
			},

			"event_hub_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"storage_account_id", "workspace_id", "event_hub_auth_rule_id", "event_hub_name"},
				RequiredWith: []string{"event_hub_name", "event_hub_auth_rule_id"},
				ValidateFunc: azure.ValidateEventHubName(),
			},

			"event_hub_auth_rule_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"storage_account_id", "workspace_id", "event_hub_auth_rule_id", "event_hub_name"},
				RequiredWith: []string{"event_hub_name", "event_hub_auth_rule_id"},
				ValidateFunc: azure.ValidateResourceID,
			},

			"logs": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						//NonInteractiveUserSignInLogs, ServicePrincipalSignInLogs and ManagedIdentitySignInLogs are not supported in azure
						// SDK-for-go enum
						"category": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{string(aad.AuditLogs),
									string(aad.SignInLogs),
									"ManagedIdentitySignInLogs",
									"NonInteractiveUserSignInLogs",
									"ProvisioningLogs",
									"ServicePrincipalSignInLogs",
								},
								false,
							),
						},

						"retention_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retention_policy_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"retention_policy_days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(0, 365),
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

func diagnosticsSettingsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Azure Active Directory Diagnostic setting %q: %+v", name, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_active_directory_diagnostic_setting", *existing.ID)
		}
	}

	diagnosticSettingsResource := aad.DiagnosticSettingsResource{
		DiagnosticSettings: &aad.DiagnosticSettings{},
	}

	if v, ok := d.GetOk("storage_account_id"); ok {
		storageAccountId := v.(string)
		diagnosticSettingsResource.DiagnosticSettings.StorageAccountID = utils.String(storageAccountId)
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		workspaceId := v.(string)
		diagnosticSettingsResource.DiagnosticSettings.WorkspaceID = utils.String(workspaceId)
	}

	if v, ok := d.GetOk("event_hub_auth_rule_id"); ok {
		eventHubAuthorizationRuleId := v.(string)
		diagnosticSettingsResource.DiagnosticSettings.EventHubAuthorizationRuleID = utils.String(eventHubAuthorizationRuleId)
	}

	if v, ok := d.GetOk("event_hub_name"); ok {
		eventHubName := v.(string)
		diagnosticSettingsResource.DiagnosticSettings.EventHubName = utils.String(eventHubName)
	}

	if v, ok := d.GetOk("logs"); ok {
		res, err := expandDiagnosticLogSettings(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("parsing `logs`: %+v", err)
		}
		diagnosticSettingsResource.DiagnosticSettings.Logs = res
	}

	if _, err := client.CreateOrUpdate(ctx, diagnosticSettingsResource, name); err != nil {
		return fmt.Errorf("error creating/updating Active Directory Diagnostic Setting %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Azure Active Directory Diagnostic setting %q: %+v", name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("null ID returned for Azure Active Directory Diagnostic settings %q", name)
	}

	d.SetId(*resp.ID)
	return diagnosticsSettingsRead(d, meta)
}

func diagnosticsSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := parseDiagnosticSettingResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Resource ID: %+v", err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Active Directory Diagnostic Setting %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("while retrieving Active Directory Diagnostic Setting %q: %+v", name, err)
	}

	d.Set("name", name)

	if settings := resp.DiagnosticSettings; settings != nil {
		d.Set("event_hub_auth_rule_id", settings.EventHubAuthorizationRuleID)
		d.Set("event_hub_name", settings.EventHubName)
		d.Set("storage_account_id", settings.StorageAccountID)
		d.Set("workspace_id", settings.WorkspaceID)
		d.Set("logs", flattenDiagnosticSettingLogs(settings.Logs))
	}
	return nil
}

func diagnosticsSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := parseDiagnosticSettingResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Resource ID: %+v", err)
	}

	if _, err := client.Delete(ctx, name); err != nil {
		return fmt.Errorf("Error while deleting AAD Diagnostic settings %s: %+v", name, err)
	}

	return nil
}

func expandDiagnosticLogSettings(logs []interface{}) (*[]aad.LogSettings, error) {
	result := make([]aad.LogSettings, 0, len(logs))

	for _, raw := range logs {
		logSettings := raw.(map[string]interface{})
		logSettingCategory := logSettings["category"].(string)
		logSettingEnabled := logSettings["enabled"].(bool)
		var retentionPolicyObject *aad.RetentionPolicy
		retentionPolicy := logSettings["retention_policy"].([]interface{})
		if len(retentionPolicy) > 0 {
			rVal, ok := retentionPolicy[0].(map[string]interface{})
			if ok {
				retentionPolicyEnabled := rVal["retention_policy_enabled"].(bool)
				retentionPolicyDays := rVal["retention_policy_days"].(int)
				if retentionPolicyEnabled && !logSettingEnabled {
					return nil, fmt.Errorf("The log setting %s should be enabled for retention policy to be applied.", logSettingCategory)
				}
				retentionPolicyObject = &aad.RetentionPolicy{
					Enabled: utils.Bool(retentionPolicyEnabled),
					Days:    utils.Int32(int32(retentionPolicyDays)),
				}
			}
		}

		result = append(result,
			aad.LogSettings{
				Category:        aad.Category(logSettingCategory),
				Enabled:         utils.Bool(logSettingEnabled),
				RetentionPolicy: retentionPolicyObject,
			},
		)
	}

	return &result, nil
}

func flattenDiagnosticSettingLogs(in *[]aad.LogSettings) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(*in))
	for _, logSetting := range *in {
		resource := make(map[string]interface{})
		retentionPolicy := make([]interface{}, 0, 1)
		resource["category"] = string(logSetting.Category)
		resource["enabled"] = logSetting.Enabled
		if *logSetting.Enabled {
			if logSetting.RetentionPolicy != nil && *logSetting.RetentionPolicy.Enabled {
				v := make(map[string]interface{})
				v["retention_policy_enabled"] = logSetting.RetentionPolicy.Enabled
				v["retention_policy_days"] = logSetting.RetentionPolicy.Days
				retentionPolicy = append(retentionPolicy, v)
				resource["retention_policy"] = retentionPolicy
			}
			result = append(result, resource)
		}
	}

	return result
}

// AAD Diagnostic setting resource ID does not comform to any subscriptions or resource groups since it is set at Directory level.
// Hence the parseResourceId function is not valid for AAD diagnostic setting resource id
func parseDiagnosticSettingResourceId(id string) (string, error) {
	idURL, err := url.Parse(id)
	if err != nil {
		return "", fmt.Errorf("Cannot parse AAD Diagnostic settings ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	if len(components)%2 != 0 {
		return "", fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for diagnosticSettings.
		if strings.EqualFold(key, "diagnosticSettings") && value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("Could not parse AAD Diagnostic setting name")
}
