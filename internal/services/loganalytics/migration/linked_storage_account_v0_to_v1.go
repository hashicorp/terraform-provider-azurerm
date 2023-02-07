package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/linkedstorageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = LinkedStorageAccountV0ToV1{}

type LinkedStorageAccountV0ToV1 struct{}

func (LinkedStorageAccountV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := linkedstorageaccounts.ParseDataSourceTypeIDInsensitively(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}

func (LinkedStorageAccountV0ToV1) Schema() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"data_source_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(linkedstorageaccounts.DataSourceTypeCustomLogs),
				string(linkedstorageaccounts.DataSourceTypeAzureWatson),
				string(linkedstorageaccounts.DataSourceTypeQuery),
				string(linkedstorageaccounts.DataSourceTypeAlerts),
				string(linkedstorageaccounts.DataSourceTypeIngestion),
			}, !features.FourPointOhBeta()),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: linkedstorageaccounts.ValidateWorkspaceID,
		},

		"storage_account_ids": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}

	if !features.FourPointOh() {
		schema["data_source_type"].DiffSuppressFunc = suppress.CaseDifference
	}
	return schema
}
