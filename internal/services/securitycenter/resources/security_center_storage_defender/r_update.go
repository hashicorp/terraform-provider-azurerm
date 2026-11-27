// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_storage_defender

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2025-06-01/defenderforstorage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan StorageDefenderModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			prop := model.Properties
			if prop == nil {
				return fmt.Errorf("retrieving %s: properties was nil", *id)
			}

			if metadata.ResourceData.HasChange("override_subscription_settings_enabled") {
				prop.OverrideSubscriptionLevelSettings = pointer.To(plan.OverrideSubscriptionSettings)
			}

			if prop.MalwareScanning == nil {
				prop.MalwareScanning = &defenderforstorage.MalwareScanningProperties{}
			}

			if prop.MalwareScanning.OnUpload == nil {
				prop.MalwareScanning.OnUpload = &defenderforstorage.OnUploadProperties{}
			}

			if metadata.ResourceData.HasChange("malware_scanning_on_upload_enabled") {
				prop.MalwareScanning.OnUpload.IsEnabled = pointer.To(plan.MalwareScanningOnUploadEnabled)
			}

			if metadata.ResourceData.HasChange("malware_scanning_on_upload_cap_gb_per_month") {
				prop.MalwareScanning.OnUpload.CapGBPerMonth = pointer.To(plan.MalwareScanningOnUploadCapPerMon)
			}

			if metadata.ResourceData.HasChange("scan_results_event_grid_topic_id") {
				prop.MalwareScanning.ScanResultsEventGridTopicResourceId = pointer.To(plan.ScanResultsEventGridTopicId)
			}

			if prop.SensitiveDataDiscovery == nil {
				prop.SensitiveDataDiscovery = &defenderforstorage.SensitiveDataDiscoveryProperties{}
			}

			if metadata.ResourceData.HasChange("sensitive_data_discovery_enabled") {
				prop.SensitiveDataDiscovery.IsEnabled = pointer.To(plan.SensitiveDataDiscoveryEnabled)
			}

			input := defenderforstorage.DefenderForStorageSetting{
				Properties: prop,
			}

			if _, err = client.Create(ctx, *id, input); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
