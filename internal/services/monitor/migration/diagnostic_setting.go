package migration

import (
	"context"
	"log"
	"strings"

	authRuleParse "github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = DiagnosticSettingUpgradeV0ToV1{}

type DiagnosticSettingUpgradeV0ToV1 struct{}

func (DiagnosticSettingUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return diagnosticSettingSchemaForV0AndV1()
}

func (DiagnosticSettingUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.KeyVault/vaults/vault1|logMonitoring1
		oldId := rawState["id"].(string)
		newId := CorrectDiagnosticSettingIdResourceGroup(oldId)
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func CorrectDiagnosticSettingIdResourceGroup(oldId string) string {
	idSegments := strings.Split(oldId, "/")
	if len(idSegments) > 4 && idSegments[3] == "resourcegroups" {
		idSegments[3] = "resourceGroups"
		return strings.Join(idSegments, "/")
	}
	return oldId
}

func diagnosticSettingSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MonitorDiagnosticSettingName,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"eventhub_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_authorization_rule_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: authRuleParse.ValidateAuthorizationRuleID,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: storageValidate.StorageAccountID,
		},

		"log_analytics_destination_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: false,
			ValidateFunc: validation.StringInSlice([]string{
				"Dedicated",
				"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
			}, false),
		},

		"log": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"category_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},
							},
						},
					},
				},
			},
			Set: validate.ResourceMonitorDiagnosticLogSettingHash,
		},

		"metric": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},
							},
						},
					},
				},
			},
			Set: validate.ResourceMonitorDiagnosticMetricsSettingHash,
		},
	}
}
