// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_feature

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appconfiguration/2024-05-01/configurationstores"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appconfiguration/1.0/appconfiguration"
)

func (k Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model FeatureResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			configurationStoreEndpoint, err := metadata.Client.AppConfiguration.EndpointForConfigurationStore(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("retrieving Endpoint for feature %q in %q: %s", model.Name, *configurationStoreId, err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(*configurationStoreEndpoint)
			if err != nil {
				return err
			}

			// users can customize the key, but if they don't we use the name
			rawKey := model.Name
			if model.Key != "" {
				rawKey = model.Key
			}
			featureKey := fmt.Sprintf("%s/%s", FeatureKeyPrefix, rawKey)

			nestedItemId, err := parse.NewNestedItemID(client.Endpoint, featureKey, model.Label)
			if err != nil {
				return err
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return errors.New("internal-error: context had no deadline")
			}

			// from https://learn.microsoft.com/azure/azure-app-configuration/concept-enable-rbac#azure-built-in-roles-for-azure-app-configuration
			// allow some time for role permission to be propagated
			stateConf := &pluginsdk.StateChangeConf{
				Pending:                   []string{"Forbidden"},
				Target:                    []string{"Error", "Exists", "NotFound"},
				Refresh:                   helpers.AppConfigurationGetKeyRefreshFunc(ctx, client, featureKey, model.Label),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Feature %q read permission to be propagated: %+v", featureKey, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				kv, err := client.GetKeyValue(ctx, featureKey, model.Label, "", "", "", []appconfiguration.KeyValueFields{})
				if err != nil {
					if v, ok := err.(autorest.DetailedError); ok {
						if !response.WasNotFound(v.Response) {
							return fmt.Errorf("got http status code %d while checking for key's %q existence: %+v", v.Response.StatusCode, featureKey, v.Error())
						}
					} else {
						return fmt.Errorf("while checking for key's %q existence: %+v", featureKey, err)
					}
				} else if kv.StatusCode == 200 {
					return tf.ImportAsExistsError(k.ResourceType(), nestedItemId.ID())
				}
			}

			entity := appconfiguration.KeyValue{
				Key:         pointer.To(featureKey),
				Label:       pointer.To(model.Label),
				Tags:        tags.Expand(model.Tags),
				ContentType: pointer.To(FeatureKeyContentType),
				Locked:      pointer.To(model.Locked),
			}

			value := FeatureValue{
				ID:          model.Name,
				Description: model.Description,
				Enabled:     model.Enabled,
			}

			value.Conditions.ClientFilters.Filters = make([]any, 0)

			if model.PercentageFilter > 0 {
				value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, PercentageFeatureFilter{
					Name:       PercentageFilterName,
					Parameters: PercentageFilterParameters{Value: model.PercentageFilter},
				})
			}

			if len(model.TargetingFilters) > 0 {
				for _, tgtf := range model.TargetingFilters {
					value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, TargetingFeatureFilter{
						Name:       TargetingFilterName,
						Parameters: TargetingFilterParameters{Audience: tgtf},
					})
				}
			}

			if len(model.TimewindowFilters) > 0 {
				for _, twf := range model.TimewindowFilters {
					value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, TimewindowFeatureFilter{
						Name:       TimewindowFilterName,
						Parameters: twf,
					})
				}
			}

			if len(model.CustomFilters) > 0 {
				for _, cf := range model.CustomFilters {
					value.Conditions.ClientFilters.Filters = append(value.Conditions.ClientFilters.Filters, CustomFilter{
						Name:       cf.Name,
						Parameters: cf.Parameters,
					})
				}
			}

			valueBytes, err := json.Marshal(value)
			if err != nil {
				return fmt.Errorf("while marshalling FeatureValue struct: %+v", err)
			}
			entity.Value = pointer.To(string(valueBytes))
			if _, err = client.PutKeyValue(ctx, featureKey, model.Label, &entity, "", ""); err != nil {
				return err
			}

			if model.Locked {
				if _, err = client.PutLock(ctx, featureKey, model.Label, "", ""); err != nil {
					return fmt.Errorf("while locking key/label pair %s/%s: %+v", model.Name, model.Label, err)
				}
			} else {
				if _, err = client.DeleteLock(ctx, featureKey, model.Label, "", ""); err != nil {
					return fmt.Errorf("while unlocking key/label pair %s/%s: %+v", model.Name, model.Label, err)
				}
			}

			// https://github.com/Azure/AppConfiguration/issues/763
			stateConf = &pluginsdk.StateChangeConf{
				Pending:                   []string{"NotFound", "Forbidden"},
				Target:                    []string{"Exists"},
				Refresh:                   helpers.AppConfigurationGetKeyRefreshFunc(ctx, client, featureKey, model.Label),
				PollInterval:              5 * time.Second,
				ContinuousTargetOccurence: 4,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Feature %q to be provisioned: %+v", featureKey, err)
			}

			metadata.SetID(nestedItemId)
			return nil
		},
		Timeout: 45 * time.Minute,
	}
}
