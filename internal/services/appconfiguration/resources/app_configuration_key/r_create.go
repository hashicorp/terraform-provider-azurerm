// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package app_configuration_key

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
			var model KeyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			configurationStoreId, err := configurationstores.ParseConfigurationStoreID(model.ConfigurationStoreId)
			if err != nil {
				return err
			}

			configurationStoreEndpoint, err := metadata.Client.AppConfiguration.EndpointForConfigurationStore(ctx, *configurationStoreId)
			if err != nil {
				return fmt.Errorf("retrieving Endpoint for feature %q in %q: %s", model.Key, *configurationStoreId, err)
			}

			client, err := metadata.Client.AppConfiguration.DataPlaneClientWithEndpoint(*configurationStoreEndpoint)
			if err != nil {
				return err
			}

			nestedItemId, err := parse.NewNestedItemID(client.Endpoint, model.Key, model.Label)
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
				Refresh:                   helpers.AppConfigurationGetKeyRefreshFunc(ctx, client, model.Key, model.Label),
				PollInterval:              10 * time.Second,
				ContinuousTargetOccurence: 3,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q read permission to be propagated: %+v", model.Key, err)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				kv, err := client.GetKeyValue(ctx, model.Key, model.Label, "", "", "", []appconfiguration.KeyValueFields{})
				if err != nil {
					if v, ok := err.(autorest.DetailedError); ok {
						if !response.WasNotFound(v.Response) {
							return fmt.Errorf("checking for presence of existing %s: %+v", nestedItemId, err)
						}
					} else {
						return fmt.Errorf("while checking for key's %q existence: %+v", model.Key, err)
					}
				} else if kv.StatusCode == 200 {
					return tf.ImportAsExistsError(k.ResourceType(), nestedItemId.ID())
				}
			}

			entity := appconfiguration.KeyValue{
				Key:   pointer.To(model.Key),
				Label: pointer.To(model.Label),
				Tags:  tags.Expand(model.Tags),
			}

			switch model.Type {
			case KeyTypeKV:
				entity.ContentType = pointer.To(model.ContentType)
				entity.Value = pointer.To(model.Value)
			case KeyTypeVault:
				entity.ContentType = pointer.To(VaultKeyContentType)
				ref, err := json.Marshal(VaultKeyReference{URI: model.VaultKeyReference})
				if err != nil {
					return fmt.Errorf("while encoding vault key reference: %+v", err)
				}
				entity.Value = pointer.To(string(ref))
			}

			if _, err = client.PutKeyValue(ctx, model.Key, model.Label, &entity, "", ""); err != nil {
				return err
			}

			if model.Locked {
				if _, err = client.PutLock(ctx, model.Key, model.Label, "", ""); err != nil {
					return fmt.Errorf("while locking key/label pair %q/%q: %+v", model.Key, model.Label, err)
				}
			}

			// https://github.com/Azure/AppConfiguration/issues/763
			stateConf = &pluginsdk.StateChangeConf{
				Pending:                   []string{"NotFound", "Forbidden"},
				Target:                    []string{"Exists"},
				Refresh:                   helpers.AppConfigurationGetKeyRefreshFunc(ctx, client, model.Key, model.Label),
				PollInterval:              5 * time.Second,
				ContinuousTargetOccurence: 4,
				Timeout:                   time.Until(deadline),
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for App Configuration Key %q to be provisioned: %+v", model.Key, err)
			}

			metadata.SetID(nestedItemId)
			return nil
		},
		Timeout: 45 * time.Minute,
	}
}
