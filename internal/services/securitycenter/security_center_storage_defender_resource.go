// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/topics"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2022-12-01-preview/defenderforstorage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StorageDefenderResource struct{}

type StorageDefenderModel struct {
	StorageAccountId                 string `tfschema:"storage_account_id"`
	OverrideSubscriptionSettings     bool   `tfschema:"override_subscription_settings_enabled"`
	MalwareScanningOnUploadEnabled   bool   `tfschema:"malware_scanning_on_upload_enabled"`
	MalwareScanningOnUploadCapPerMon int64  `tfschema:"malware_scanning_on_upload_cap_gb_per_month"`
	SensitiveDataDiscoveryEnabled    bool   `tfschema:"sensitive_data_discovery_enabled"`
	ScanResultsEventGridTopicId      string `tfschema:"scan_results_event_grid_topic_id"`
}

var _ sdk.ResourceWithUpdate = StorageDefenderResource{}

func (s StorageDefenderResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateStorageAccountID
}

func (s StorageDefenderResource) ModelObject() interface{} {
	return &StorageDefenderModel{}
}

func (s StorageDefenderResource) ResourceType() string {
	return "azurerm_security_center_storage_defender"
}

func (s StorageDefenderResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"override_subscription_settings_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"malware_scanning_on_upload_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"malware_scanning_on_upload_cap_gb_per_month": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  -1,
			ValidateFunc: validation.Any(
				validation.IntAtLeast(1),
				validation.IntInSlice([]int{-1}),
			),
		},

		"sensitive_data_discovery_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"scan_results_event_grid_topic_id": commonschema.ResourceIDReferenceOptional(&topics.TopicId{}),
	}
}

func (s StorageDefenderResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s StorageDefenderResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan StorageDefenderModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id := commonids.NewScopeID(plan.StorageAccountId)

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

			_, err = client.Create(ctx, id, input)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s StorageDefenderResource) Update() sdk.ResourceFunc {
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
				return fmt.Errorf("decoding %+v", err)
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

			_, err = client.Create(ctx, *id, input)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (s StorageDefenderResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			storageAccountId, err := commonids.ParseStorageAccountID(id.Scope)
			if err != nil {
				return err
			}
			state := StorageDefenderModel{
				StorageAccountId: storageAccountId.ID(),
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					if !pointer.From(prop.IsEnabled) {
						return metadata.MarkAsGone(id)
					}

					state.OverrideSubscriptionSettings = pointer.From(prop.OverrideSubscriptionLevelSettings)

					if ms := prop.MalwareScanning; ms != nil {
						if onUpload := ms.OnUpload; onUpload != nil {
							state.MalwareScanningOnUploadEnabled = pointer.From(onUpload.IsEnabled)
							state.MalwareScanningOnUploadCapPerMon = pointer.From(onUpload.CapGBPerMonth)
						}
						if ms.ScanResultsEventGridTopicResourceId != nil {
							topicId, err := topics.ParseTopicID(*ms.ScanResultsEventGridTopicResourceId)
							if err != nil {
								return err
							}
							state.ScanResultsEventGridTopicId = topicId.ID()
						}
					}

					if sdd := prop.SensitiveDataDiscovery; sdd != nil {
						state.SensitiveDataDiscoveryEnabled = pointer.From(sdd.IsEnabled)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s StorageDefenderResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			_, err = client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			input := defenderforstorage.DefenderForStorageSetting{
				Properties: &defenderforstorage.DefenderForStorageSettingProperties{
					IsEnabled: pointer.To(false),
				},
			}

			_, err = client.Create(ctx, *id, input)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
