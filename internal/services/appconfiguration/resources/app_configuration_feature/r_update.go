// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func (k Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			nestedItemId, err := parse.ParseNestedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(nestedItemId.ConfigurationStoreEndpoint)
			if err != nil {
				return err
			}

			kv, err := client.GetKeyValue(ctx, nestedItemId.Key, nestedItemId.Label, "", "", "", []appconfiguration.KeyValueFields{})
			if err != nil {
				return fmt.Errorf("while checking for key %q existence: %+v", *nestedItemId, err)
			}

			var fv FeatureValue
			if err = json.Unmarshal([]byte(pointer.From(kv.Value)), &fv); err != nil {
				return fmt.Errorf("while unmarshalling underlying key's value: %+v", err)
			}

			var model FeatureResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			metadata.Client.AppConfiguration.AddToCache(*configurationStoreId, nestedItemId.ConfigurationStoreEndpoint)

			// Remove the lock, if any. We will put it back again if the model says so.
			if _, err = client.DeleteLock(ctx, nestedItemId.Key, nestedItemId.Label, "", ""); err != nil {
				return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", nestedItemId.Key, nestedItemId.Label, err)
			}

			if metadata.ResourceData.HasChange("tags") {
				kv.Tags = tags.Expand(model.Tags)
			}

			if metadata.ResourceData.HasChange("locked") {
				kv.Locked = pointer.To(model.Locked)
			}

			if metadata.ResourceData.HasChange("enabled") {
				fv.Enabled = model.Enabled
			}

			if metadata.ResourceData.HasChange("description") {
				fv.Description = model.Description
			}

			filters := make([]any, 0)
			filterChanged := false
			timewindowFilters := make([]any, 0)
			targetingFilters := make([]any, 0)
			percentageFilter := PercentageFeatureFilter{}
			customFilters := make([]any, 0)
			if len(fv.Conditions.ClientFilters.Filters) > 0 {
				for _, f := range fv.Conditions.ClientFilters.Filters {
					switch f := f.(type) {
					case TimewindowFeatureFilter:
						twfp := f
						timewindowFilters = append(timewindowFilters, twfp)
					case TargetingFeatureFilter:
						tfp := f
						targetingFilters = append(targetingFilters, tfp)
					case PercentageFeatureFilter:
						percentageFilter = f
					case CustomFilter:
						cf := f
						customFilters = append(customFilters, cf)
					default:
						return fmt.Errorf("while unmarshaling feature payload: unknown filter type %+v", f)
					}
				}
			}

			if metadata.ResourceData.HasChange("percentage_filter_value") {
				filters = append(filters, PercentageFeatureFilter{
					Name:       PercentageFilterName,
					Parameters: PercentageFilterParameters{Value: model.PercentageFilter},
				})
				filterChanged = true
			} else if percentageFilter.Name != "" {
				filters = append(filters, percentageFilter)
			}

			if metadata.ResourceData.HasChange("targeting_filter") {
				for _, tgtf := range model.TargetingFilters {
					filters = append(filters, TargetingFeatureFilter{
						Name:       TargetingFilterName,
						Parameters: TargetingFilterParameters{Audience: tgtf},
					})
				}
				filterChanged = true
			} else {
				filters = append(filters, targetingFilters...)
			}

			if metadata.ResourceData.HasChange("timewindow_filter") {
				for _, twf := range model.TimewindowFilters {
					filters = append(filters, TimewindowFeatureFilter{
						Name:       TimewindowFilterName,
						Parameters: twf,
					})
				}
				filterChanged = true
			} else {
				filters = append(filters, timewindowFilters...)
			}

			if metadata.ResourceData.HasChange("custom_filter") {
				for _, cf := range model.CustomFilters {
					filters = append(filters, CustomFilter{
						Name:       cf.Name,
						Parameters: cf.Parameters,
					})
				}
				filterChanged = true
			} else {
				filters = append(filters, customFilters...)
			}

			if filterChanged {
				fv.Conditions.ClientFilters.Filters = filters
			}

			valueBytes, err := json.Marshal(fv)
			if err != nil {
				return fmt.Errorf("while marshalling FeatureValue struct: %+v", err)
			}
			kv.Value = pointer.To(string(valueBytes))
			if _, err = client.PutKeyValue(ctx, nestedItemId.Key, model.Label, &kv, "", ""); err != nil {
				return err
			}

			if model.Locked {
				if _, err = client.PutLock(ctx, nestedItemId.Key, model.Label, "", ""); err != nil {
					return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Name, model.Label, err)
				}
			} else {
				if _, err = client.DeleteLock(ctx, nestedItemId.Key, model.Label, "", ""); err != nil {
					return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Name, model.Label, err)
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
