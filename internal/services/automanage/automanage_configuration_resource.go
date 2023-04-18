package automanage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

type ConfigurationModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`

	Antimalware               *AntimalwareConfiguration `tfschema:"antimalware"`
	AutomationAccountEnabled  bool                      `tfschema:"automation_account_enabled"`
	BootDiagnosticsEnabled    bool                      `tfschema:"boot_diagnostics_enabled"`
	ChangeTrackingEnabled     bool                      `tfschema:"change_tracking_enabled"`
	DefenderForCloudEnabled   bool                      `tfschema:"defender_for_cloud_enabled"`
	GuestConfigurationEnabled bool                      `tfschema:"guest_configuration_enabled"`
	StatusChangeAlertEnabled  bool                      `tfschema:"status_change_alert_enabled"`

	Location string            `tfschema:"location"`
	Tags     map[string]string `tfschema:"tags"`
}

type AntimalwareConfiguration struct {
	Enabled                   bool                   `tfschema:"enabled"`
	Exclusions                *AntimalwareExclusions `tfschema:"exclusions"`
	RealTimeProtectionEnabled bool                   `tfschema:"real_time_protection_enabled"`
	ScheduledScanEnabled      bool                   `tfschema:"scheduled_scan_enabled"`
	ScanType                  string                 `tfschema:"scheduled_scan_type"`
	ScanDay                   int                    `tfschema:"scheduled_scan_day"`
	ScanTimeInMinutes         int                    `tfschema:"scheduled_scan_time_in_minutes"`
}

type AntimalwareExclusions struct {
	Extensions string `tfschema:"extensions"`
	Paths      string `tfschema:"paths"`
	Processes  string `tfschema:"processes"`
}

type AutoManageConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = AutoManageConfigurationResource{}

func (r AutoManageConfigurationResource) ResourceType() string {
	return "azurerm_automanage_configuration"
}

func (r AutoManageConfigurationResource) ModelObject() interface{} {
	return &ConfigurationModel{}
}

func (r AutoManageConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.AutomanageConfigurationID
}

func (r AutoManageConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		//"Antimalware/Enable": boolean,
		//"Antimalware/EnableRealTimeProtection": boolean,
		//"Antimalware/RunScheduledScan": boolean,
		//"Antimalware/ScanType": string ("Quick", "Full"),
		//"Antimalware/ScanDay": int (0-8) Ex: 0 - daily, 1 - Sunday, 2 - Monday, .... 7 - Saturday, 8 - Disabled,
		//"Antimalware/ScanTimeInMinutes": int (0 - 1440),
		//"Antimalware/Exclusions/Extensions": string (extensions separated by semicolon. Ex: ".ext1;.ext2"),
		//"Antimalware/Exclusions/Paths": string (Paths separated by semicolon. Ex: "c:\excluded-path-1;c:\excluded-path-2"),
		//"Antimalware/Exclusions/Processes": string (Processes separated by semicolon. Ex: "proc1.exe;proc2.exe"),
		"antimalware": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"real_time_protection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"scheduled_scan_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"scheduled_scan_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Quick",
						ValidateFunc: validation.StringInSlice([]string{
							"Quick",
							"Full",
						}, false),
					},
					"scheduled_scan_day": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
						ValidateFunc: validation.IntInSlice([]int{
							0, 1, 2, 3, 4, 5, 6, 7, 8,
						}),
					},
					"scheduled_scan_time_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 1440),
					},
					"exclusions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"extensions": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"paths": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"processes": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		//"AutomationAccount/Enable": boolean,
		"automation_account_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		//"BootDiagnostics/Enable": boolean,
		"boot_diagnostics_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		//"ChangeTrackingAndInventory/Enable": boolean,
		"change_tracking_and_inventory_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "DefenderForCloud/Enable": boolean,
		"defender_for_cloud_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		//"GuestConfiguration/Enable": boolean,
		"guest_configuration_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		//"Alerts/AutomanageStatusChanges/Enable": boolean,
		"status_change_alert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AutoManageConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AutoManageConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Automanage.ConfigurationClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewAutomanageConfigurationID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := automanage.ConfigurationProfile{
				Location:   utils.String(location.Normalize(model.Location)),
				Properties: &automanage.ConfigurationProfileProperties{},
				Tags:       tags.FromTypedObject(model.Tags),
			}

			// Convert all to a map[string]interface{} and convert to a json property
			jsonConfig := make(map[string]interface{})

			if model.Antimalware != nil {
				jsonConfig["Antimalware/Enable"] = model.Antimalware.Enabled
				jsonConfig["Antimalware/RealTimeProtectionEnabled"] = model.Antimalware.RealTimeProtectionEnabled
				jsonConfig["Antimalware/RunScheduledScan"] = model.Antimalware.ScheduledScanEnabled
				jsonConfig["Antimalware/ScanType"] = model.Antimalware.ScanType
				jsonConfig["Antimalware/ScanDay"] = model.Antimalware.ScanDay
				jsonConfig["Antimalware/ScanTimeInMinutes"] = model.Antimalware.ScanTimeInMinutes
				if model.Antimalware.Exclusions != nil {
					jsonConfig["Antimalware/Exclusions/Extensions"] = model.Antimalware.Exclusions.Extensions
					jsonConfig["Antimalware/Exclusions/Paths"] = model.Antimalware.Exclusions.Paths
					jsonConfig["Antimalware/Exclusions/Processes"] = model.Antimalware.Exclusions.Processes
				}
			}

			if model.AutomationAccountEnabled {
				jsonConfig["AutomationAccount/Enable"] = model.AutomationAccountEnabled
			}

			if model.BootDiagnosticsEnabled {
				jsonConfig["BootDiagnostics/Enable"] = model.BootDiagnosticsEnabled
			}

			if model.ChangeTrackingEnabled {
				jsonConfig["ChangeTracking/Enable"] = model.ChangeTrackingEnabled
			}

			if model.DefenderForCloudEnabled {
				jsonConfig["DefenderForCloud/Enable"] = model.DefenderForCloudEnabled
			}

			if model.GuestConfigurationEnabled {
				jsonConfig["GuestConfiguration/Enable"] = model.GuestConfigurationEnabled
			}

			if model.StatusChangeAlertEnabled {
				jsonConfig["Alerts/AutomanageStatusChanges/Enable"] = model.StatusChangeAlertEnabled
			}

			properties.Properties.Configuration = &jsonConfig

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.ResourceGroup, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AutoManageConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationClient

			id, err := parse.AutomanageConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			jsonConfig := make(map[string]interface{})

			if resp.Properties != nil && resp.Properties.Configuration != nil {
				err := json.Unmarshal([]byte(resp.Properties.Configuration.(string)), &jsonConfig)
				if err != nil {
					return fmt.Errorf("unmarshalling %s: %+v", *id, err)
				}
			}

			if model.Antimalware != nil {
				jsonConfig["Antimalware/Enable"] = model.Antimalware.Enabled
				jsonConfig["Antimalware/RealTimeProtectionEnabled"] = model.Antimalware.RealTimeProtectionEnabled
				jsonConfig["Antimalware/RunScheduledScan"] = model.Antimalware.ScheduledScanEnabled
				jsonConfig["Antimalware/ScanType"] = model.Antimalware.ScanType
				jsonConfig["Antimalware/ScanDay"] = model.Antimalware.ScanDay
				jsonConfig["Antimalware/ScanTimeInMinutes"] = model.Antimalware.ScanTimeInMinutes
				if model.Antimalware.Exclusions != nil {
					jsonConfig["Antimalware/Exclusions/Extensions"] = model.Antimalware.Exclusions.Extensions
					jsonConfig["Antimalware/Exclusions/Paths"] = model.Antimalware.Exclusions.Paths
					jsonConfig["Antimalware/Exclusions/Processes"] = model.Antimalware.Exclusions.Processes
				}
			}

			if metadata.ResourceData.HasChange("automation_account_enabled") {
				jsonConfig["AutomationAccount/Enable"] = model.AutomationAccountEnabled
			}

			if metadata.ResourceData.HasChange("boot_diagnostics_enabled") {
				jsonConfig["BootDiagnostics/Enable"] = model.BootDiagnosticsEnabled
			}

			if metadata.ResourceData.HasChange("change_tracking_enabled") {
				jsonConfig["ChangeTracking/Enable"] = model.ChangeTrackingEnabled
			}

			if metadata.ResourceData.HasChange("defender_for_cloud_enabled") {
				jsonConfig["DefenderForCloud/Enable"] = model.DefenderForCloudEnabled
			}

			if metadata.ResourceData.HasChange("guest_configuration_enabled") {
				jsonConfig["GuestConfiguration/Enable"] = model.GuestConfigurationEnabled
			}

			if metadata.ResourceData.HasChange("status_change_alert_enabled") {
				jsonConfig["Alerts/AutomanageStatusChanges/Enable"] = model.StatusChangeAlertEnabled
			}

			if metadata.ResourceData.HasChange("tags") {
				resp.Tags = tags.FromTypedObject(model.Tags)
			}

			configBytes, err := json.Marshal(jsonConfig)
			if err != nil {
				return fmt.Errorf("marshalling %s: %+v", *id, err)
			}

			properties := automanage.ConfigurationProfile{
				Location: utils.String(metadata.ResourceData.Get("location").(string)),
				Properties: &automanage.ConfigurationProfileProperties{
					Configuration: configBytes,
				},
				Tags: resp.Tags,
			}

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.ResourceGroup, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r AutoManageConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationClient

			id, err := parse.AutomanageConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ConfigurationModel{
				Name:              id.ConfigurationProfileName,
				ResourceGroupName: id.ResourceGroup,
				Location:          location.NormalizeNilable(resp.Location),
			}

			if resp.Properties != nil && resp.Properties.Configuration != nil {
				configMap := make(map[string]interface{})
				err := json.Unmarshal(resp.Properties.Configuration.([]byte), &configMap)
				if err != nil {
					return fmt.Errorf("unmarshalling %s: %+v", *id, err)
				}

				if configMap["Antimalware/Enable"] != nil {
					state.Antimalware = &AntimalwareConfiguration{
						Enabled:                   configMap["Antimalware/Enable"].(bool),
						RealTimeProtectionEnabled: configMap["Antimalware/RealTimeProtectionEnabled"].(bool),
						ScheduledScanEnabled:      configMap["Antimalware/RunScheduledScan"].(bool),
						ScanType:                  configMap["Antimalware/ScanType"].(string),
						ScanDay:                   configMap["Antimalware/ScanDay"].(int),
						ScanTimeInMinutes:         configMap["Antimalware/ScanTimeInMinutes"].(int),
					}

					if configMap["Antimalware/Exclusions/Extensions"] != nil {
						state.Antimalware.Exclusions = &AntimalwareExclusions{
							Extensions: configMap["Antimalware/Exclusions/Extensions"].(string),
							Paths:      configMap["Antimalware/Exclusions/Paths"].(string),
							Processes:  configMap["Antimalware/Exclusions/Processes"].(string),
						}
					}
				}

				if configMap["AutomationAccount/Enable"] != nil {
					state.AutomationAccountEnabled = configMap["AutomationAccount/Enable"].(bool)
				}

				if configMap["BootDiagnostics/Enable"] != nil {
					state.BootDiagnosticsEnabled = configMap["BootDiagnostics/Enable"].(bool)
				}

				if configMap["ChangeTracking/Enable"] != nil {
					state.ChangeTrackingEnabled = configMap["ChangeTracking/Enable"].(bool)
				}

				if configMap["DefenderForCloud/Enable"] != nil {
					state.DefenderForCloudEnabled = configMap["DefenderForServers/Enable"].(bool)
				}

				if configMap["GuestConfiguration/Enable"] != nil {
					state.GuestConfigurationEnabled = configMap["GuestConfiguration/Enable"].(bool)
				}

				if configMap["Alerts/AutomanageStatusChanges/Enable"] != nil {
					state.StatusChangeAlertEnabled = configMap["Alerts/AutomanageStatusChanges/Enable"].(bool)
				}
			}

			if resp.Tags != nil {
				state.Tags = tags.ToTypedObject(resp.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AutoManageConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationClient

			id, err := parse.AutomanageConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileName); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
