package automanage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"regexp"
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

type AutoManageConfigurationModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Configuration     string            `tfschema:"configuration_json"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
}

type AutoManageConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = AutoManageConfigurationResource{}

func (r AutoManageConfigurationResource) ResourceType() string {
	return "azurerm_automanage_configuration"
}

func (r AutoManageConfigurationResource) ModelObject() interface{} {
	return &AutoManageConfigurationModel{}
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

		//"Alerts/AutomanageStatusChanges/Enable": boolean,
		"status_change_alert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

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
					"scan_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "Quick",
						ValidateFunc: validation.StringInSlice([]string{
							"Quick",
							"Full",
						}, false),
					},
					"scan_day": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
						ValidateFunc: validation.IntInSlice([]int{
							0, 1, 2, 3, 4, 5, 6, 7, 8,
						}),
					},
					"scan_time_in_minutes": {
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

		//"AzureSecurityBaseline/Enable": boolean,
		//"AzureSecurityBaseline/AssignmentType": string ("ApplyAndAutoCorrect", "ApplyAndMonitor", "Audit", "DeployAndAutoCorrect"),
		"azure_security_baseline": {
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
					"assignment_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "ApplyAndAutoCorrect",
						ValidateFunc: validation.StringInSlice([]string{
							"ApplyAndAutoCorrect",
							"ApplyAndMonitor",
							"Audit",
							"DeployAndAutoCorrect",
						}, false),
					},
				},
			},
		},

		//"Backup/Enable": boolean,
		//"Backup/PolicyName": string (length 3 - 150, begin with alphanumeric char, only contain alphanumeric chars and hyphens),
		//"Backup/TimeZone": timezone,
		//"Backup/InstantRpRetentionRangeInDays": int (1 - 5 if ScheduleRunFrequency is Daily, 5 if ScheduleRunFrequency is Weekly),
		//"Backup/SchedulePolicy/ScheduleRunFrequency": string ("Daily", "Weekly"),
		//"Backup/SchedulePolicy/ScheduleRunTimes": list of DateTime,
		//"Backup/SchedulePolicy/ScheduleRunDays": list of strings (["Sunday", "Monday", "Wednesday", "Thursday", "Friday", "Saturday"]),
		//"Backup/SchedulePolicy/SchedulePolicyType": string ("SimpleSchedulePolicy"),
		//"Backup/RetentionPolicy/RetentionPolicyType": string ("LongTermRetentionPolicy"),
		//"Backup/RetentionPolicy/DailySchedule/RetentionTimes": list of DateTime,
		//"Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count": int (7 - 9999),
		//"Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType": string ("Days"),
		//"Backup/RetentionPolicy/WeeklySchedule/RetentionTimes":, list of DateTime
		//"Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count":, int (1 - 5163)
		//"Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType": string ("Weeks"),
		"backup": {
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
					"policy_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,149}$`), "Policy name must be 3 - 150 characters long, begin with an alphanumeric character, and only contain alphanumeric characters and hyphens."),
					},
					"time_zone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},
					"instant_rp_retention_range_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      5,
						ValidateFunc: validation.IntBetween(1, 5),
					},
					"schedule_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"schedule_run_frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Daily",
									ValidateFunc: validation.StringInSlice([]string{
										"Daily",
										"Weekly",
									}, false),
								},
								"schedule_run_times": {
									Type:     pluginsdk.TypeList,
									Optional: true,

									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},
								"schedule_run_days": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Sunday",
											"Monday",
											"Tuesday",
											"Wednesday",
											"Thursday",
											"Friday",
											"Saturday",
										}, false),
									},
								},
								"schedule_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "SimpleSchedulePolicy",
									ValidateFunc: validation.StringInSlice([]string{
										"SimpleSchedulePolicy",
									}, false),
								},
							},
						},
					},
					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"retention_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "LongTermRetentionPolicy",
									ValidateFunc: validation.StringInSlice([]string{
										"LongTermRetentionPolicy",
									}, false),
								},
								"daily_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.IsRFC3339Time,
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															Default:      7,
															ValidateFunc: validation.IntBetween(7, 9999),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Days",
															ValidateFunc: validation.StringInSlice([]string{
																"Days",
															}, false),
														},
													},
												},
											},
										},
									},
								},
								"weekly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.IsRFC3339Time,
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															Default:      4,
															ValidateFunc: validation.IntBetween(4, 9999),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Weeks",
															ValidateFunc: validation.StringInSlice([]string{
																"Weeks",
															}, false),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
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

		//"LogAnalytics/Enable": boolean,
		//"LogAnalytics/Reprovision": boolean,
		//"LogAnalytics/Workspace": resource ID (Log analytics workspace ID),
		"log_analytics": {
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
					"reprovision": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"workspace_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},

		//"TrustedLaunchVM/Backup/Enable": boolean,
		//"TrustedLaunchVM/Backup/PolicyName": string (length 3 - 150, begin with alphanumeric char, only contain alphanumeric chars and hyphens),
		//"TrustedLaunchVM/Backup/TimeZone": timezone,
		//"TrustedLaunchVM/Backup/InstantRpRetentionRangeInDays": int (1 - 30),
		//"TrustedLaunchVM/Backup/SchedulePolicy/ScheduleRunFrequency": string ("Hourly", "Daily", "Weekly"),
		//"TrustedLaunchVM/Backup/SchedulePolicy/SchedulePolicyType": string ("SimpleSchedulePolicyV2"),
		//"TrustedLaunchVM/Backup/RetentionPolicy/RetentionPolicyType":, string ("LongTermRetentionPolicy")
		//"TrustedLaunchVM/Backup/SchedulePolicy/HourlySchedule/Interval": int (4, 6, 8, 12),
		//"TrustedLaunchVM/Backup/SchedulePolicy/HourlySchedule/ScheduleWindowStartTime": DateTime,
		//"TrustedLaunchVM/Backup/SchedulePolicy/HourlySchedule/ScheduleWindowDuration": int (4, 8, 12, 16, 20, 24),
		//"TrustedLaunchVM/Backup/RetentionPolicy/DailySchedule/RetentionTimes": list of DateTime,
		//"TrustedLaunchVM/Backup/RetentionPolicy/DailySchedule/RetentionDuration/Count": int (7 - 9999),
		//"TrustedLaunchVM/Backup/RetentionPolicy/DailySchedule/RetentionDuration/DurationType": string ("Hours", "Days"),
		//"TrustedLaunchVM/Backup/SchedulePolicy/DailySchedule/ScheduleWindowStartTime": DateTime,
		//"TrustedLaunchVM/Backup/SchedulePolicy/WeeklySchedule/ScheduleRunDays": list of strings (["Sunday", "Monday", "Wednesday", "Thursday", "Friday", "Saturday"]),
		//"TrustedLaunchVM/Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/Count": int (1 - 1563),
		//"TrustedLaunchVM/Backup/RetentionPolicy/WeeklySchedule/RetentionDuration/DurationType": string ("Weeks"),
		"trusted_launch_vm_backup": {
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
					"policy_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{2,149}$`), "policy name must be 3 - 150 characters long, begin with an alphanumeric character, and only contain alphanumeric characters and hyphens"),
					},
					"time_zone": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "UTC",
					},
					"instant_rp_retention_range_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(1, 30),
					},
					"schedule_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"schedule_run_frequency": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Daily",
									ValidateFunc: validation.StringInSlice([]string{
										"Hourly",
										"Daily",
										"Weekly",
									}, false),
								},
								"schedule_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "SimpleSchedulePolicyV2",
									ValidateFunc: validation.StringInSlice([]string{
										"SimpleSchedulePolicyV2",
									}, false),
								},
								"hourly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"interval": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      4,
												ValidateFunc: validation.IntInSlice([]int{4, 6, 8, 12}),
											},
											"schedule_window_start_time": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.IsRFC3339Time,
											},
											"schedule_window_duration": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												Default:      4,
												ValidateFunc: validation.IntInSlice([]int{4, 8, 12, 16, 20, 24}),
											},
										},
									},
								},
								"daily_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"schedule_window_start_time": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.IsRFC3339Time,
											},
										},
									},
								},
								"weekly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"schedule_run_days": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice([]string{
														"Sunday",
														"Monday",
														"Tuesday",
														"Wednesday",
														"Thursday",
														"Friday",
														"Saturday",
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},
					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"retention_policy_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "LongTermRetentionPolicy",
									ValidateFunc: validation.StringInSlice([]string{
										"LongTermRetentionPolicy",
									}, false),
								},
								"daily_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_times": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.IsRFC3339Time,
												},
											},
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															Default:      7,
															ValidateFunc: validation.IntBetween(7, 9999),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Days",
															ValidateFunc: validation.StringInSlice([]string{
																"Days",
															}, false),
														},
													},
												},
											},
										},
									},
								},
								"weekly_schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"retention_duration": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												MaxItems: 1,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"count": {
															Type:         pluginsdk.TypeInt,
															Optional:     true,
															Default:      4,
															ValidateFunc: validation.IntBetween(4, 9999),
														},
														"duration_type": {
															Type:     pluginsdk.TypeString,
															Optional: true,
															Default:  "Weeks",
															ValidateFunc: validation.StringInSlice([]string{
																"Weeks",
															}, false),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		//"UpdateManagement/Enable": boolean,
		"update_management_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		//"VMInsights/Enable": boolean,
		"vm_insights_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		//"WindowsAdminCenter/Enable": boolean,
		"windows_admin_center_enabled": {
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
			var model AutoManageConfigurationModel
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

			if model.Configuration != "" {
				var configurationValue interface{}
				err = json.Unmarshal([]byte(model.Configuration), &configurationValue)
				if err != nil {
					return err
				}
				properties.Properties.Configuration = &configurationValue
			}

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

			var model AutoManageConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("configuration_json") {
				var configurationValue interface{}
				err := json.Unmarshal([]byte(model.Configuration), &configurationValue)
				if err != nil {
					return err
				}

				resp.Properties.Configuration = &configurationValue
			}

			if metadata.ResourceData.HasChange("tags") {
				resp.Tags = tags.FromTypedObject(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, id.ConfigurationProfileName, id.ResourceGroup, resp); err != nil {
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

			state := AutoManageConfigurationModel{
				Name:              id.ConfigurationProfileName,
				ResourceGroupName: id.ResourceGroup,
				Location:          location.NormalizeNilable(resp.Location),
			}

			if properties := resp.Properties; properties != nil {
				if properties.Configuration != nil {
					configurationValue, err := json.Marshal(properties.Configuration)
					if err != nil {
						return err
					}

					state.Configuration = string(configurationValue)
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

func backupResourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"backup_retention_days": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(7, 9999),
				},

				"storage_account_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
		},
	}
}
