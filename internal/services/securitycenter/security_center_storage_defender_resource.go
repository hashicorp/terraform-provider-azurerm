package securitycenter

import (
	"context"
	"fmt"
	"time"

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
	Enabled                          bool   `tfschema:"enabled"`
	OverrideSubscriptionSettings     bool   `tfschema:"override_subscription_settings_enabled"`
	MalwareScanningOnUploadEnabled   bool   `tfschema:"malware_scanning_on_upload_enabled"`
	MalwareScanningOnUploadCapPerMon int64  `tfschema:"malware_scanning_on_upload_cap_gb_per_month"`
	SensitiveDataDiscoveryEnabled    bool   `tfschema:"sensitive_data_discovery_enabled"`
}

var _ sdk.ResourceWithUpdate = StorageDefenderResource{}

func (s StorageDefenderResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateScopeID
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

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
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
			ValidateFunc: func(i interface{}, s string) (warnings []string, errors []error) {
				// it requires -1 or greater than 0
				v, ok := i.(int)
				if !ok {
					errors = append(errors, fmt.Errorf("expected type of %s to be integer", s))
					return warnings, errors
				}

				if v == -1 {
					return warnings, errors
				}

				return validation.IntAtLeast(-1)(i, s)
			},
		},

		"sensitive_data_discovery_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (s StorageDefenderResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s StorageDefenderResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
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
					IsEnabled:                         pointer.To(plan.Enabled),
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

			_, err = client.Create(ctx, id, input)
			if err != nil {
				return fmt.Errorf("creating: %+v", err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s StorageDefenderResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
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

			if metadata.ResourceData.HasChange("enabled") {
				prop.IsEnabled = pointer.To(plan.Enabled)
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
				return fmt.Errorf("updating: %+v", err)
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
				return fmt.Errorf("reading %+v", err)
			}

			state := StorageDefenderModel{
				StorageAccountId: id.ID(),
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					state.Enabled = pointer.From(prop.IsEnabled)
					state.OverrideSubscriptionSettings = pointer.From(prop.OverrideSubscriptionLevelSettings)

					if ms := prop.MalwareScanning; ms != nil {
						if onUpload := ms.OnUpload; onUpload != nil {
							state.MalwareScanningOnUploadEnabled = pointer.From(onUpload.IsEnabled)
							state.MalwareScanningOnUploadCapPerMon = pointer.From(onUpload.CapGBPerMonth)
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
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SecurityCenter.DefenderForStorageClient

			id, err := commonids.ParseScopeID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("reading %+v", err)
				}
			}
			// if the resource has never been created, it returns 404.
			// once created, it could only be set to disable.
			if response.WasNotFound(resp.HttpResponse) {
				return nil
			}

			input := defenderforstorage.DefenderForStorageSetting{
				Properties: &defenderforstorage.DefenderForStorageSettingProperties{
					IsEnabled: pointer.To(false),
				},
			}

			_, err = client.Create(ctx, *id, input)
			if err != nil {
				return fmt.Errorf("deleting %+v", err)
			}

			return nil
		},
	}
}
