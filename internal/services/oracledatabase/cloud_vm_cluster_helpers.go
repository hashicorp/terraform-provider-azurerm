package oracledatabase

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func ConvertDataCollectionOptionsToInternal(dataCollectionOptions *cloudvmclusters.DataCollectionOptions) []DataCollectionOptionsModel {
	if dataCollectionOptions != nil {
		return []DataCollectionOptionsModel{
			{
				IsDiagnosticsEventsEnabled: pointer.From(dataCollectionOptions.IsDiagnosticsEventsEnabled),
				IsHealthMonitoringEnabled:  pointer.From(dataCollectionOptions.IsHealthMonitoringEnabled),
				IsIncidentLogsEnabled:      pointer.From(dataCollectionOptions.IsIncidentLogsEnabled),
			},
		}
	}
	return nil
}

func ConvertExadataIormConfigToInternal(exadataIormConfig *cloudvmclusters.ExadataIormConfig) []ExadataIormConfigModel {
	if exadataIormConfig != nil {
		var dbIormConfigModel []DbIormConfigModel
		if exadataIormConfig.DbPlans != nil {
			dbPlans := *exadataIormConfig.DbPlans
			for _, dbPlan := range dbPlans {
				dbIormConfigModel = append(dbIormConfigModel, DbIormConfigModel{
					DbName:          pointer.From(dbPlan.DbName),
					FlashCacheLimit: pointer.From(dbPlan.FlashCacheLimit),
					Share:           pointer.From(dbPlan.Share),
				})
			}
		}
		return []ExadataIormConfigModel{
			{
				DbPlans:          dbIormConfigModel,
				LifecycleDetails: pointer.From(exadataIormConfig.LifecycleDetails),
				LifecycleState:   string(pointer.From(exadataIormConfig.LifecycleState)),
				Objective:        string(pointer.From(exadataIormConfig.Objective)),
			},
		}
	}
	return nil
}

func GiVersionDiffSuppress(key string, old string, new string, d *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}
	oldVersion := strings.Split(old, ".")
	newVersion := strings.Split(new, ".")

	if oldVersion[0] == newVersion[0] {
		return true
	}
	return false
}

func DbSystemHostnameDiffSuppress(key string, old string, new string, d *schema.ResourceData) bool {
	return EqualIgnoreCaseSuppressDiff(key, old, new, d) || NewIsPrefixOfOldDiffSuppress(key, old, new, d)
}

func NewIsPrefixOfOldDiffSuppress(key string, old string, new string, d *schema.ResourceData) bool {
	return strings.HasPrefix(strings.ToLower(old), strings.ToLower(new))
}

func EqualIgnoreCaseSuppressDiff(key string, old string, new string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}
