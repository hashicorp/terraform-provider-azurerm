package aadmgmt

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/aad/mgmt/2017-04-01/aad"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"net/url"
	"strings"
	"time"
)

func azureADDiagnosticSettingsResource() *schema.Resource {
	return &schema.Resource{
		Create: diagnosticsSettingsCreateOrUpdate,
		Read:   diagnosticsSettingsRead,
		Update: diagnosticsSettingsCreateOrUpdate,
		Delete: diagnosticsSettingsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Type:     schema.TypeString,
				Optional: true,
			},

			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"event_hub_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"event_hub_auth_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"logs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},

						//NonInteractiveUserSignInLogs, ServicePrincipalSignInLogs and ManagedIdentitySignInLogs are not supported in azure
						// SDK-for-go enum even though it is supported in REST API v2017-04-01
						"category": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{string(aad.AuditLogs), string(aad.SignInLogs)},
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

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func diagnosticsSettingsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	storageAccount := d.Get("storage_account_id").(string)
	workspace := d.Get("workspace_id").(string)
	eventHubAuthRule := d.Get("event_hub_auth_rule_id").(string)
	eventHub := d.Get("event_hub_name").(string)

	if storageAccount == "" && workspace == "" && (eventHubAuthRule == "" || eventHub == "") {
		return fmt.Errorf("Please specify atleast one data sink for diagnostic settings %q", client.BaseURI)
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

	if _, ok := d.GetOk("logs"); ok {
		res, err := expandDiagnosticLogSettings(d)
		if err != nil {
			return fmt.Errorf("state validation failed for diagnostic settings %q: %+v", name, err)
		}
		diagnosticSettingsResource.DiagnosticSettings.Logs = res
	}

	if _, err := client.CreateOrUpdate(ctx, diagnosticSettingsResource, name); err != nil {
		return fmt.Errorf("error creating diagnostic settings for AAD %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to retrieve Azure Active Directory Diagnostic settings %q: %+v", name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read ID for Azure Active Directory Diagnostic settings %q", name)
	}

	d.SetId(*resp.ID)
	return diagnosticsSettingsRead(d, meta)
}

func expandDiagnosticLogSettings(d *schema.ResourceData) (*[]aad.LogSettings, error) {
	diagnosticLogSettings := d.Get("logs").([]interface{})
	result := make([]aad.LogSettings, 0)

	for _, raw := range diagnosticLogSettings {
		logSettings := raw.(map[string]interface{})
		logSettingCategory := logSettings["category"].(string)
		logSettingEnabled := logSettings["enabled"].(bool)
		var retentionPolicyObject *aad.RetentionPolicy
		retentionPolicy := logSettings["retention_policy"].([]interface{})
		if retentionPolicy != nil && len(retentionPolicy) > 0 {
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

func diagnosticsSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resMap, err := parseDiagnosticSettingResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("error while parsing Resource Id:%+v", err)
	}
	name := resMap["diagnosticSettings"]

	if _, err := client.Delete(ctx, name); err != nil {
		return fmt.Errorf("Error while deleting AAD Diagnostic settings %s: %+v", name, err)
	}

	return nil
}

func diagnosticsSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AADManagement.DiagnosticSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resMap, err := parseDiagnosticSettingResourceId(d.Id())
	if err != nil {
		return fmt.Errorf("error while parsing Resource Id:%+v", err)
	}
	name := resMap["diagnosticSettings"]
	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error while fetching AAD Diagnostic settings %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("type", resp.Type)

	if settings := resp.DiagnosticSettings; settings != nil {
		d.Set("event_hub_auth_rule_id", settings.EventHubAuthorizationRuleID)
		d.Set("event_hub_name", settings.EventHubName)
		d.Set("storage_account_id", settings.StorageAccountID)
		d.Set("workspace_id", settings.WorkspaceID)
		d.Set("logs", flattenDiagnosticSettingLogs(settings.Logs))
	}
	return nil
}

func flattenDiagnosticSettingLogs(in *[]aad.LogSettings) []map[string]interface{} {
	if in == nil {
		return []map[string]interface{}{}
	}
	result := make([]map[string]interface{}, 0, len(*in))
	for _, logSetting := range *in {
		resource := make(map[string]interface{})
		retentionPolicy := make([]interface{}, 0, 1)
		if aad.AuditLogs == logSetting.Category || aad.SignInLogs == logSetting.Category {
			resource["category"] = string(logSetting.Category)
			resource["enabled"] = logSetting.Enabled
			if logSetting.RetentionPolicy != nil {
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
func parseDiagnosticSettingResourceId(id string) (map[string]string, error) {
	idURL, err := url.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse AAD Diagnostic settings ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	componentMap := make(map[string]string, len(components)/2)
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}
		componentMap[key] = value
	}

	return componentMap, nil

}
