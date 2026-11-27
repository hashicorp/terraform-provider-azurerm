// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/topics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2025-06-01/defenderforstorage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan StorageDefenderModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := securitycenterClient.From(metadata.Client).DefenderForStorageClient

			id := commonids.NewScopeID(plan.StorageAccountId)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				resp, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(resp.HttpResponse) &&
					resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.IsEnabled != nil && *resp.Model.Properties.IsEnabled {
					return tf.ImportAsExistsError(s.ResourceType(), id.ID())
				}
			}

			input := defenderforstorage.DefenderForStorageSetting{
				Properties: &defenderforstorage.DefenderForStorageSettingProperties{
					IsEnabled:                         pointer.To(true),
					OverrideSubscriptionLevelSettings: pointer.To(plan.OverrideSubscriptionSettings),
					MalwareScanning: &defenderforstorage.MalwareScanningProperties{
						OnUpload: &defenderforstorage.OnUploadProperties{
							IsEnabled:     pointer.To(plan.MalwareScanningOnUploadEnabled),
							CapGBPerMonth: pointer.To(plan.MalwareScanningOnUploadCapPerMon),
						},
					},
					SensitiveDataDiscovery: &defenderforstorage.SensitiveDataDiscoveryProperties{
						IsEnabled: pointer.To(plan.SensitiveDataDiscoveryEnabled),
					},
				},
			}

			if plan.ScanResultsEventGridTopicId != "" {
				topicId, err := topics.ParseTopicID(plan.ScanResultsEventGridTopicId)
				if err != nil {
					return err
				}
				input.Properties.MalwareScanning.ScanResultsEventGridTopicResourceId = pointer.To(topicId.ID())
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
