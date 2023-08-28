// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate4 "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	RebootSettingIfRequired = "IfRequired"
	RebootSettingNever      = "Never"
	RebootSettingAlways     = "Always"
	RebootSettingRebootOnly = "RebootOnly"

	FrequencyOneTime = "OneTime"
	FrequencyDay     = "Day"
	FrequencyHour    = "Hour"
	FrequencyWeek    = "Week"
	FrequencyMonth   = "Month"

	DaysOfWeekMonday    = "Monday"
	DaysOfWeekTuesday   = "Tuesday"
	DaysOfWeekWednesday = "Wednesday"
	DaysOfWeekThursday  = "Thursday"
	DaysOfWeekFriday    = "Friday"
	DaysOfWeekSaturday  = "Saturday"
	DaysOfWeekSunday    = "Sunday"
)

type Tag struct {
	Tag    string   `tfschema:"tag"`
	Values []string `tfschema:"values"`
}

type AzureQuery struct {
	Scope     []string `tfschema:"scope"`
	Locations []string `tfschema:"locations"`
	Tags      []Tag    `tfschema:"tags"`
	TagFilter string   `tfschema:"tag_filter"`
}

type Linux struct {
	Reboot           string   `tfschema:"reboot"`
	Classifications  []string `tfschema:"classifications_included"`
	ExcludedPackages []string `tfschema:"excluded_packages"`
	IncludedPackages []string `tfschema:"included_packages"`

	Classification string `tfschema:"classification_included"` // Deprecated use Classifications instead
}

type MonthlyOccurrence struct {
	Occurrence int    `tfschema:"occurrence"`
	Day        string `tfschema:"day"`
}

type NonAzureQuery struct {
	FunctionAlias string `tfschema:"function_alias"`
	WorkspaceId   string `tfschema:"workspace_id"`
}

type UpdateTask struct {
	Source     string            `tfschema:"source"`
	Parameters map[string]string `tfschema:"parameters"`
}

type Schedule struct {
	Description             string              `tfschema:"description"`
	StartTime               string              `tfschema:"start_time"`
	StartTimeOffsetMinutes  float64             `tfschema:"start_time_offset_minutes"`
	ExpiryTime              string              `tfschema:"expiry_time"`
	ExpiryTimeOffsetMinutes float64             `tfschema:"expiry_time_offset_minutes"`
	IsEnabled               bool                `tfschema:"is_enabled"`
	NextRun                 string              `tfschema:"next_run"`
	NextRunOffsetMinutes    float64             `tfschema:"next_run_offset_minutes"`
	Interval                int                 `tfschema:"interval"`
	Frequency               string              `tfschema:"frequency"`
	CreationTime            string              `tfschema:"creation_time"`
	LastModifiedTime        string              `tfschema:"last_modified_time"`
	TimeZone                string              `tfschema:"time_zone"`
	AdvancedWeekDays        []string            `tfschema:"advanced_week_days"`
	AdvancedMonthDays       []int               `tfschema:"advanced_month_days"`
	MonthlyOccurrence       []MonthlyOccurrence `tfschema:"monthly_occurrence"`
}

type Target struct {
	AzureQueries    []AzureQuery    `tfschema:"azure_query"`
	NonAzureQueries []NonAzureQuery `tfschema:"non_azure_query"`
}

type Windows struct {
	Classifications []string `tfschema:"classifications_included"`
	ExcludedKbs     []string `tfschema:"excluded_knowledge_base_numbers"`
	IncludedKbs     []string `tfschema:"included_knowledge_base_numbers"`
	RebootSetting   string   `tfschema:"reboot"`

	Classification string `tfschema:"classification_included"` // Deprecated use Classifications instead

}

type SoftwareUpdateConfigurationModel struct {
	AutomationAccountID   string       `tfschema:"automation_account_id"`
	Name                  string       `tfschema:"name"`
	ErrorCode             string       `tfschema:"error_code"`
	ErrorMeesage          string       `tfschema:"error_meesage"`
	ErrorMessage          string       `tfschema:"error_message"`
	OperatingSystem       string       `tfschema:"operating_system"`
	Linux                 []Linux      `tfschema:"linux"`
	Windows               []Windows    `tfschema:"windows"`
	Duration              string       `tfschema:"duration"`
	VirtualMachines       []string     `tfschema:"virtual_machine_ids"`
	NonAzureComputerNames []string     `tfschema:"non_azure_computer_names"`
	Targets               []Target     `tfschema:"target"`
	Schedule              []Schedule   `tfschema:"schedule"`
	PreTask               []UpdateTask `tfschema:"pre_task"`
	PostTask              []UpdateTask `tfschema:"post_task"`
}

type SoftwareUpdateConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = SoftwareUpdateConfigurationResource{}

func (m SoftwareUpdateConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	linux := pluginsdk.Resource{}
	windows := pluginsdk.Resource{}
	if !features.FourPointOhBeta() {
		linux = pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"reboot": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  RebootSettingIfRequired,
					ValidateFunc: validation.StringInSlice([]string{
						RebootSettingAlways,
						RebootSettingIfRequired,
						RebootSettingNever,
						RebootSettingRebootOnly,
					}, false),
				},

				"classification_included": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					ConflictsWith: []string{"windows.0.classifications_included"},
					Computed:      true,
					ValidateFunc: validation.StringInSlice(
						softwareupdateconfiguration.PossibleValuesForLinuxUpdateClasses(),
						false),
					Deprecated: "", // TODO
				},

				"classifications_included": {
					Type:          pluginsdk.TypeList,
					Optional:      true,
					Computed:      true,
					ConflictsWith: []string{"windows.0.classification_included"},
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							softwareupdateconfiguration.PossibleValuesForLinuxUpdateClasses(),
							false),
					},
				},

				"excluded_packages": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"included_packages": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		}
		windows = pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"classification_included": {
					Type:          pluginsdk.TypeString,
					Optional:      true,
					Computed:      true,
					ConflictsWith: []string{"windows.0.classifications_included"},
					Deprecated:    "windows classification can be set as a list, use `classifications_included` instead.",
					ValidateFunc: validation.StringInSlice(func() (vs []string) {
						vs = append(vs, softwareupdateconfiguration.PossibleValuesForWindowsUpdateClasses()...)
						return
					}(), false),
				},

				"classifications_included": {
					Type:          pluginsdk.TypeList,
					Optional:      true,
					Computed:      true,
					ConflictsWith: []string{"windows.0.classification_included"},
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(func() (res []string) {
							res = append(res, softwareupdateconfiguration.PossibleValuesForWindowsUpdateClasses()...)
							return
						}(), false),
					},
				},

				"excluded_knowledge_base_numbers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"included_knowledge_base_numbers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"reboot": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  RebootSettingIfRequired,
					ValidateFunc: validation.StringInSlice([]string{
						RebootSettingAlways,
						RebootSettingIfRequired,
						RebootSettingNever,
						RebootSettingRebootOnly,
					}, false),
				},
			},
		}
	} else {
		linux = pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"reboot": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  RebootSettingIfRequired,
					ValidateFunc: validation.StringInSlice([]string{
						RebootSettingAlways,
						RebootSettingIfRequired,
						RebootSettingNever,
						RebootSettingRebootOnly,
					}, false),
				},

				"classifications_included": {
					Type:     pluginsdk.TypeList,
					Required: true, // TODO 4.0 - Update docs to reflect this breaking change.
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(softwareupdateconfiguration.PossibleValuesForLinuxUpdateClasses(),
							false),
					},
				},

				"excluded_packages": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"included_packages": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		}
		windows = pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"classifications_included": {
					Type:     pluginsdk.TypeList,
					Required: true, // TODO 4.0 - Update docs to reflect this breaking change.
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(
							softwareupdateconfiguration.PossibleValuesForWindowsUpdateClasses(),
							false),
					},
				},

				"excluded_knowledge_base_numbers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"included_knowledge_base_numbers": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},

				"reboot": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Default:  RebootSettingIfRequired,
					ValidateFunc: validation.StringInSlice([]string{
						RebootSettingAlways,
						RebootSettingIfRequired,
						RebootSettingNever,
						RebootSettingRebootOnly,
					}, false),
				},
			},
		}
	}
	r := map[string]*pluginsdk.Schema{
		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: softwareupdateconfiguration.ValidateAutomationAccountID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"linux": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     &linux,
			ExactlyOneOf: []string{
				"windows",
				"linux",
			},
		},

		"windows": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem:     &windows,
			ExactlyOneOf: []string{
				"windows",
				"linux",
			},
		},

		"duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT2H",
			ValidateFunc: validate4.ISO8601Duration,
		},

		"virtual_machine_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: computeValidate.VirtualMachineID,
			},
		},

		"non_azure_computer_names": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"target": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"azure_query": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"scope": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										// Subscription or Resource Group ARM Id
										ValidateFunc: func(i interface{}, s string) ([]string, []error) {
											w, e := validate2.ResourceGroupID(i, s)
											if len(e) == 0 {
												return w, e
											}
											w, e = commonids.ValidateSubscriptionID(i, s)
											return w, e
										},
									},
								},

								"locations": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"tags": {
									Type:     schema.TypeList,
									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"tag": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"values": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &schema.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},

								"tag_filter": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Computed: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(softwareupdateconfiguration.TagOperatorsAny),
										string(softwareupdateconfiguration.TagOperatorsAll),
									}, false),
								},
							},
						},
					},

					"non_azure_query": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"function_alias": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"workspace_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"start_time": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						DiffSuppressFunc: suppress.RFC3339MinuteTime,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					"start_time_offset_minutes": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Computed: true,
					},

					"expiry_time": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						DiffSuppressFunc: suppress.RFC3339MinuteTime,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					"expiry_time_offset_minutes": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Computed: true,
					},

					"is_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"next_run": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						DiffSuppressFunc: suppress.RFC3339MinuteTime,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					"next_run_offset_minutes": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Computed: true,
					},

					"interval": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							FrequencyOneTime,
							FrequencyHour,
							FrequencyDay,
							FrequencyWeek,
							FrequencyMonth,
						}, false),
					},

					"creation_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"last_modified_time": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"time_zone": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "Etc/UTC",
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"advanced_week_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								DaysOfWeekMonday,
								DaysOfWeekTuesday,
								DaysOfWeekWednesday,
								DaysOfWeekThursday,
								DaysOfWeekFriday,
								DaysOfWeekSaturday,
								DaysOfWeekSunday,
							}, false),
						},
					},

					"advanced_month_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(1, 31),
						},
					},

					"monthly_occurrence": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"occurrence": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 5),
								},

								"day": {
									Type:     pluginsdk.TypeString,
									Required: true,
									// not hardcode Enum values
									ValidateFunc: func(i interface{}, s string) ([]string, []error) {
										var vs []string
										vs = append(vs, softwareupdateconfiguration.PossibleValuesForScheduleDay()...)
										vf := validation.StringInSlice(vs, false)
										return vf(i, s)
									},
								},
							},
						},
					},
				},
			},
		},

		"pre_task": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: nil,
						},
					},
				},
			},
		},

		"post_task": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: nil,
						},
					},
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		r["operating_system"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(softwareupdateconfiguration.OperatingSystemTypeLinux),
				string(softwareupdateconfiguration.OperatingSystemTypeWindows),
			}, false),
			Deprecated: "This property has been deprecated and will be removed in a future release. The use of either the `linux` or `windows` blocks replaces setting this value directly. This value is ignored by the provider.",
		}
	}

	return r
}

func (m SoftwareUpdateConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"error_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// TODO 4.0 remove & update docs
		"error_meesage": {
			Type:       pluginsdk.TypeString,
			Computed:   true,
			Deprecated: "`error_meesage` will be removed in favour of `error_message` in version 4.0 of the AzureRM Provider",
		},

		"error_message": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (m SoftwareUpdateConfigurationResource) ModelObject() interface{} {
	return &SoftwareUpdateConfigurationModel{}
}

func (m SoftwareUpdateConfigurationResource) ResourceType() string {
	return "azurerm_automation_software_update_configuration"
}

func (m SoftwareUpdateConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Automation.SoftwareUpdateConfigClient

			var model SoftwareUpdateConfigurationModel
			if err := meta.Decode(&model); err != nil {
				return err
			}

			automationID, err := automationaccount.ParseAutomationAccountID(model.AutomationAccountID)
			if err != nil {
				return err
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := softwareupdateconfiguration.NewSoftwareUpdateConfigurationID(subscriptionID, automationID.ResourceGroupName, automationID.AutomationAccountName, model.Name)
			existing, err := client.SoftwareUpdateConfigurationsGetByName(ctx, id, softwareupdateconfiguration.SoftwareUpdateConfigurationsGetByNameOperationOptions{})
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return meta.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			updateConfig := expandUpdateConfig(model)
			if _, err = client.SoftwareUpdateConfigurationsCreate(ctx, id, *updateConfig, softwareupdateconfiguration.SoftwareUpdateConfigurationsCreateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (m SoftwareUpdateConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := softwareupdateconfiguration.ParseSoftwareUpdateConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.SoftwareUpdateConfigClient
			resp, err := client.SoftwareUpdateConfigurationsGetByName(ctx, *id, softwareupdateconfiguration.SoftwareUpdateConfigurationsGetByNameOperationOptions{})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return err
			}

			var state SoftwareUpdateConfigurationModel

			state.AutomationAccountID = softwareupdateconfiguration.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()

			props := resp.Model.Properties
			updateConfiguration := props.UpdateConfiguration
			scheduleConfiguration := props.ScheduleInfo

			state.Name = id.SoftwareUpdateConfigurationName
			state.Duration = pointer.From(updateConfiguration.Duration)
			if linux := updateConfiguration.Linux; linux != nil {
				l := Linux{
					Reboot:           pointer.From(linux.RebootSetting),
					Classifications:  strings.Split(string(pointer.From(linux.IncludedPackageClassifications)), ", "),
					ExcludedPackages: pointer.From(linux.ExcludedPackageNameMasks),
					IncludedPackages: pointer.From(linux.IncludedPackageNameMasks),
				}

				if !features.FourPointOhBeta() && len(strings.Split(string(pointer.From(linux.IncludedPackageClassifications)), ", ")) == 1 {
					l.Classification = string(pointer.From(linux.IncludedPackageClassifications))
				}

				state.Linux = []Linux{l}
				state.OperatingSystem = string(softwareupdateconfiguration.OperatingSystemTypeLinux)
			}
			if windows := updateConfiguration.Windows; windows != nil {
				w := Windows{
					Classifications: strings.Split(strings.ReplaceAll(string(pointer.From(windows.IncludedUpdateClassifications)), " ", ""), ","),
					ExcludedKbs:     pointer.From(windows.ExcludedKbNumbers),
					IncludedKbs:     pointer.From(windows.IncludedKbNumbers),
					RebootSetting:   pointer.From(windows.RebootSetting),
				}

				if !features.FourPointOhBeta() && len(strings.Split(string(pointer.From(windows.IncludedUpdateClassifications)), ", ")) == 1 {
					w.Classification = string(pointer.From(windows.IncludedUpdateClassifications))
				}

				state.Windows = []Windows{w}
				state.OperatingSystem = string(softwareupdateconfiguration.OperatingSystemTypeWindows)
			}
			if targets := updateConfiguration.Targets; targets != nil {
				t := Target{}
				aq := make([]AzureQuery, 0)
				for _, v := range pointer.From(targets.AzureQueries) {
					tagsList := make([]Tag, 0)
					tagFilter := ""
					if tags := v.TagSettings; tags != nil {
						for k, vals := range pointer.From(tags.Tags) {
							tagsList = append(tagsList, Tag{
								Tag:    k,
								Values: vals,
							})
						}
						tagFilter = string(pointer.From(tags.FilterOperator))
					}
					aq = append(aq, AzureQuery{
						Scope:     pointer.From(v.Scope),
						Locations: pointer.From(v.Locations),
						Tags:      tagsList,
						TagFilter: tagFilter,
					})
				}

				t.AzureQueries = aq

				naq := make([]NonAzureQuery, 0)
				for _, v := range pointer.From(targets.NonAzureQueries) {
					naq = append(naq, NonAzureQuery{
						FunctionAlias: pointer.From(v.FunctionAlias),
						WorkspaceId:   pointer.From(v.WorkspaceId),
					})
				}

				t.NonAzureQueries = naq
				state.Targets = []Target{t}
			}

			state.VirtualMachines = pointer.From(updateConfiguration.AzureVirtualMachines)
			state.NonAzureComputerNames = pointer.From(updateConfiguration.NonAzureComputerNames)

			schedule := Schedule{
				Description:             pointer.From(scheduleConfiguration.Description),
				StartTime:               pointer.From(scheduleConfiguration.StartTime),
				StartTimeOffsetMinutes:  pointer.From(scheduleConfiguration.StartTimeOffsetMinutes),
				ExpiryTime:              pointer.From(scheduleConfiguration.ExpiryTime),
				ExpiryTimeOffsetMinutes: pointer.From(scheduleConfiguration.ExpiryTimeOffsetMinutes),
				IsEnabled:               pointer.From(scheduleConfiguration.IsEnabled),
				NextRun:                 pointer.From(scheduleConfiguration.NextRun),
				NextRunOffsetMinutes:    pointer.From(scheduleConfiguration.NextRunOffsetMinutes),
				Interval:                int(pointer.From(scheduleConfiguration.Interval)),
				Frequency:               string(pointer.From(scheduleConfiguration.Frequency)),
				CreationTime:            pointer.From(scheduleConfiguration.CreationTime),
				LastModifiedTime:        pointer.From(scheduleConfiguration.LastModifiedTime),
				TimeZone:                pointer.From(scheduleConfiguration.TimeZone),
			}

			// (@jackofallops) - Advanced Schedule info is never returned so we'll pull it in from Config until the tracked issue is resolved
			// Tracking Issue: https://github.com/Azure/azure-rest-api-specs/issues/24436
			if advSchedule := scheduleConfiguration.AdvancedSchedule; advSchedule != nil {
				schedule.AdvancedWeekDays = pointer.From(advSchedule.WeekDays)
				if monthDays := pointer.From(advSchedule.MonthDays); len(monthDays) > 0 {
					advMonthDays := make([]int, 0)
					for _, v := range monthDays {
						advMonthDays = append(advMonthDays, int(v))
					}
					schedule.AdvancedMonthDays = advMonthDays
				}
				if monthlyOccurrence := pointer.From(advSchedule.MonthlyOccurrences); len(monthlyOccurrence) > 0 {
					mo := make([]MonthlyOccurrence, 0)
					for _, v := range monthlyOccurrence {
						mo = append(mo, MonthlyOccurrence{
							Occurrence: int(pointer.From(v.Occurrence)),
							Day:        string(pointer.From(v.Day)),
						})
					}
					schedule.MonthlyOccurrence = mo
				}
			} else {
				if weekDays, ok := meta.ResourceData.GetOk("schedule.0.advanced_week_days"); ok {
					wd := make([]string, 0)
					for _, v := range weekDays.([]interface{}) {
						wd = append(wd, v.(string))
					}
					schedule.AdvancedWeekDays = wd
				}
				if monthDays, ok := meta.ResourceData.GetOk("schedule.0.advanced_month_days"); ok {
					md := make([]int, 0)
					for _, v := range monthDays.([]interface{}) {
						md = append(md, v.(int))
					}
					schedule.AdvancedMonthDays = md
				}
				if monthlyOccurrence, ok := meta.ResourceData.GetOk("schedule.0.monthly_occurrence"); ok {
					mos := make([]MonthlyOccurrence, 0)
					if moRaw, ok := monthlyOccurrence.([]interface{}); ok {
						for _, v := range moRaw {
							mo := v.(map[string]interface{})
							mos = append(mos, MonthlyOccurrence{
								Occurrence: mo["occurrence"].(int),
								Day:        mo["day"].(string),
							})
						}
					}
					schedule.MonthlyOccurrence = mos
				}
			}

			state.Schedule = []Schedule{schedule}

			if tasks := props.Tasks; tasks != nil {
				if pre := tasks.PreTask; pre != nil {
					state.PreTask = []UpdateTask{{
						Source:     pointer.From(pre.Source),
						Parameters: pointer.From(pre.Parameters),
					}}
				}
				if post := tasks.PostTask; post != nil {
					state.PostTask = []UpdateTask{{
						Source:     pointer.From(post.Source),
						Parameters: pointer.From(post.Parameters),
					}}
				}
			}

			return meta.Encode(&state)
		},
	}
}

func (m SoftwareUpdateConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automation.SoftwareUpdateConfigClient

			id, err := softwareupdateconfiguration.ParseSoftwareUpdateConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SoftwareUpdateConfigurationModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.SoftwareUpdateConfigurationsGetByName(ctx, *id, softwareupdateconfiguration.SoftwareUpdateConfigurationsGetByNameOperationOptions{})
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			existing := resp.Model

			if metadata.ResourceData.HasChange("linux") {
				if len(model.Linux) > 0 {
					existing.Properties.UpdateConfiguration.OperatingSystem = softwareupdateconfiguration.OperatingSystemTypeLinux
					v := model.Linux[0]
					existing.Properties.UpdateConfiguration.Linux = &softwareupdateconfiguration.LinuxProperties{
						ExcludedPackageNameMasks:       pointer.To(v.ExcludedPackages),
						IncludedPackageClassifications: pointer.To(softwareupdateconfiguration.LinuxUpdateClasses(strings.Join(v.Classifications, ","))),
						IncludedPackageNameMasks:       pointer.To(v.IncludedPackages),
						RebootSetting:                  pointer.To(v.Reboot),
					}
				} else {
					existing.Properties.UpdateConfiguration.Linux = &softwareupdateconfiguration.LinuxProperties{}
				}
			}

			if metadata.ResourceData.HasChange("windows") {
				if len(model.Windows) > 0 {
					existing.Properties.UpdateConfiguration.OperatingSystem = softwareupdateconfiguration.OperatingSystemTypeWindows
					v := model.Windows[0]
					existing.Properties.UpdateConfiguration.Windows = &softwareupdateconfiguration.WindowsProperties{
						ExcludedKbNumbers:             pointer.To(v.ExcludedKbs),
						IncludedKbNumbers:             pointer.To(v.IncludedKbs),
						IncludedUpdateClassifications: pointer.To(softwareupdateconfiguration.WindowsUpdateClasses(strings.Join(v.Classifications, ","))),
						RebootSetting:                 pointer.To(v.RebootSetting),
					}
				} else {
					existing.Properties.UpdateConfiguration.Windows = &softwareupdateconfiguration.WindowsProperties{}
				}
			}

			if metadata.ResourceData.HasChange("duration") {
				existing.Properties.UpdateConfiguration.Duration = pointer.To(model.Duration)
			}

			if metadata.ResourceData.HasChange("virtual_machine_ids") {
				existing.Properties.UpdateConfiguration.AzureVirtualMachines = pointer.To(model.VirtualMachines)
			}

			if metadata.ResourceData.HasChange("non_azure_computer_names") {
				existing.Properties.UpdateConfiguration.NonAzureComputerNames = pointer.To(model.NonAzureComputerNames)
			}

			if metadata.ResourceData.HasChange("target") {
				target := softwareupdateconfiguration.TargetProperties{}
				if len(model.Targets) > 0 {
					t := model.Targets[0]
					if len(t.AzureQueries) > 0 {
						aq := make([]softwareupdateconfiguration.AzureQueryProperties, 0)
						for _, v := range t.AzureQueries {
							q := softwareupdateconfiguration.AzureQueryProperties{}
							if len(v.Locations) > 0 {
								q.Locations = pointer.To(v.Locations)
							}
							if len(v.Scope) > 0 {
								q.Scope = pointer.To(v.Scope)
							}
							if len(v.Tags) > 0 || v.TagFilter != "" {
								q.TagSettings = &softwareupdateconfiguration.TagSettingsProperties{
									FilterOperator: pointer.To(softwareupdateconfiguration.TagOperators(v.TagFilter)),
								}
								tags := make(map[string][]string)
								for _, tag := range v.Tags {
									tags[tag.Tag] = tag.Values
								}
								q.TagSettings.Tags = pointer.To(tags)
							}

							aq = append(aq, q)
						}

						target.AzureQueries = pointer.To(aq)
					} else {
						target.AzureQueries = &[]softwareupdateconfiguration.AzureQueryProperties{}
					}

					if len(t.NonAzureQueries) > 0 {
						naqs := make([]softwareupdateconfiguration.NonAzureQueryProperties, 0)
						for _, v := range t.NonAzureQueries {
							naq := softwareupdateconfiguration.NonAzureQueryProperties{}
							if v.FunctionAlias != "" {
								naq.FunctionAlias = pointer.To(v.FunctionAlias)
							}
							if v.WorkspaceId != "" {
								naq.WorkspaceId = pointer.To(v.WorkspaceId)
							}
							naqs = append(naqs, naq)
						}

						target.NonAzureQueries = pointer.To(naqs)
					} else {
						target.NonAzureQueries = &[]softwareupdateconfiguration.NonAzureQueryProperties{}
					}

				} else {
					target.AzureQueries = &[]softwareupdateconfiguration.AzureQueryProperties{}
					target.NonAzureQueries = &[]softwareupdateconfiguration.NonAzureQueryProperties{}
				}
				existing.Properties.UpdateConfiguration.Targets = pointer.To(target)
			}

			if metadata.ResourceData.HasChange("schedule") {
				if len(model.Schedule) == 1 {
					v := model.Schedule[0]
					scheduleConfig := softwareupdateconfiguration.SUCScheduleProperties{
						Description:             pointer.To(v.Description),
						ExpiryTime:              pointer.To(v.ExpiryTime),
						ExpiryTimeOffsetMinutes: pointer.To(v.ExpiryTimeOffsetMinutes),
						Frequency:               pointer.To(softwareupdateconfiguration.ScheduleFrequency(v.Frequency)),
						Interval:                pointer.To(int64(v.Interval)),
						IsEnabled:               pointer.To(v.IsEnabled),
						NextRun:                 pointer.To(v.NextRun),
						NextRunOffsetMinutes:    pointer.To(v.NextRunOffsetMinutes),
						StartTime:               pointer.To(v.StartTime),
						StartTimeOffsetMinutes:  pointer.To(v.StartTimeOffsetMinutes),
						TimeZone:                pointer.To(v.TimeZone),
					}

					if len(v.AdvancedWeekDays) > 0 || len(v.AdvancedMonthDays) > 0 || len(v.MonthlyOccurrence) > 0 {
						advSchedule := &softwareupdateconfiguration.AdvancedSchedule{}
						if len(v.AdvancedWeekDays) > 0 {
							advSchedule.WeekDays = pointer.To(v.AdvancedWeekDays)
						}

						if len(v.AdvancedMonthDays) > 0 {
							i := make([]int64, 0)
							for _, v := range v.AdvancedMonthDays {
								i = append(i, int64(v))
							}
							advSchedule.MonthDays = pointer.To(i)
						}

						if len(v.MonthlyOccurrence) > 0 {
							monthlyOccurrences := make([]softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence, 0)
							for _, mo := range v.MonthlyOccurrence {
								monthlyOccurrences = append(monthlyOccurrences, softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence{
									Day:        pointer.To(softwareupdateconfiguration.ScheduleDay(mo.Day)),
									Occurrence: pointer.To(int64(mo.Occurrence)),
								})
							}

							advSchedule.MonthlyOccurrences = pointer.To(monthlyOccurrences)
						}

						scheduleConfig.AdvancedSchedule = advSchedule
					}

					existing.Properties.ScheduleInfo = scheduleConfig
				} else {
					existing.Properties.ScheduleInfo = softwareupdateconfiguration.SUCScheduleProperties{}
				}
			}

			if metadata.ResourceData.HasChanges("pre_task", "post_task") {
				tasksConfig := &softwareupdateconfiguration.SoftwareUpdateConfigurationTasks{}
				if existing.Properties.Tasks != nil {
					tasksConfig = existing.Properties.Tasks
				}
				if len(model.PreTask) > 0 {
					v := model.PreTask[0]
					task := &softwareupdateconfiguration.TaskProperties{
						Parameters: pointer.To(v.Parameters),
						Source:     pointer.To(v.Source),
					}

					tasksConfig.PreTask = task
				} else {
					tasksConfig.PreTask = &softwareupdateconfiguration.TaskProperties{}
				}
				if len(model.PostTask) > 0 {
					v := model.PostTask[0]
					task := &softwareupdateconfiguration.TaskProperties{
						Parameters: pointer.To(v.Parameters),
						Source:     pointer.To(v.Source),
					}

					tasksConfig.PostTask = task
				} else {
					tasksConfig.PostTask = &softwareupdateconfiguration.TaskProperties{}
				}
			}

			if _, err = client.SoftwareUpdateConfigurationsCreate(ctx, *id, *existing, softwareupdateconfiguration.SoftwareUpdateConfigurationsCreateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (m SoftwareUpdateConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := softwareupdateconfiguration.ParseSoftwareUpdateConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", *id)
			client := meta.Client.Automation.SoftwareUpdateConfigClient
			if _, err = client.SoftwareUpdateConfigurationsDelete(ctx, *id, softwareupdateconfiguration.SoftwareUpdateConfigurationsDeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}
			return nil
		},
	}
}

func (m SoftwareUpdateConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return softwareupdateconfiguration.ValidateSoftwareUpdateConfigurationID
}

func expandUpdateConfig(input SoftwareUpdateConfigurationModel) *softwareupdateconfiguration.SoftwareUpdateConfiguration {
	result := &softwareupdateconfiguration.SoftwareUpdateConfiguration{
		Properties: softwareupdateconfiguration.SoftwareUpdateConfigurationProperties{
			ScheduleInfo: softwareupdateconfiguration.SUCScheduleProperties{},
		},
	}

	if len(input.Schedule) == 1 {
		v := input.Schedule[0]
		scheduleConfig := softwareupdateconfiguration.SUCScheduleProperties{
			Description:             pointer.To(v.Description),
			ExpiryTime:              pointer.To(v.ExpiryTime),
			ExpiryTimeOffsetMinutes: pointer.To(v.ExpiryTimeOffsetMinutes),
			Frequency:               pointer.To(softwareupdateconfiguration.ScheduleFrequency(v.Frequency)),
			Interval:                pointer.To(int64(v.Interval)),
			IsEnabled:               pointer.To(v.IsEnabled),
			NextRun:                 pointer.To(v.NextRun),
			NextRunOffsetMinutes:    pointer.To(v.NextRunOffsetMinutes),
			StartTime:               pointer.To(v.StartTime),
			StartTimeOffsetMinutes:  pointer.To(v.StartTimeOffsetMinutes),
			TimeZone:                pointer.To(v.TimeZone),
		}

		if len(v.AdvancedWeekDays) > 0 || len(v.AdvancedMonthDays) > 0 || len(v.MonthlyOccurrence) > 0 {
			advSchedule := &softwareupdateconfiguration.AdvancedSchedule{}
			if len(v.AdvancedWeekDays) > 0 {
				advSchedule.WeekDays = pointer.To(v.AdvancedWeekDays)
			}

			if len(v.AdvancedMonthDays) > 0 {
				i := make([]int64, 0)
				for _, v := range v.AdvancedMonthDays {
					i = append(i, int64(v))
				}
				advSchedule.MonthDays = pointer.To(i)
			}

			if len(v.MonthlyOccurrence) > 0 {
				monthlyOccurrences := make([]softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence, 0)
				for _, mo := range v.MonthlyOccurrence {
					monthlyOccurrences = append(monthlyOccurrences, softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence{
						Day:        pointer.To(softwareupdateconfiguration.ScheduleDay(mo.Day)),
						Occurrence: pointer.To(int64(mo.Occurrence)),
					})
				}

				advSchedule.MonthlyOccurrences = pointer.To(monthlyOccurrences)
			}

			scheduleConfig.AdvancedSchedule = advSchedule
		}

		result.Properties.ScheduleInfo = scheduleConfig
	}

	if len(input.PreTask) > 0 || len(input.PostTask) > 0 {
		tasksConfig := &softwareupdateconfiguration.SoftwareUpdateConfigurationTasks{}

		if len(input.PreTask) > 0 {
			v := input.PreTask[0]
			task := &softwareupdateconfiguration.TaskProperties{
				Parameters: pointer.To(v.Parameters),
				Source:     pointer.To(v.Source),
			}

			tasksConfig.PreTask = task
		}

		if len(input.PostTask) > 0 {
			v := input.PostTask[0]
			task := &softwareupdateconfiguration.TaskProperties{
				Parameters: pointer.To(v.Parameters),
				Source:     pointer.To(v.Source),
			}

			tasksConfig.PostTask = task
		}

		result.Properties.Tasks = tasksConfig
	}

	updateConfig := softwareupdateconfiguration.UpdateConfiguration{}

	if len(input.VirtualMachines) > 0 {
		updateConfig.AzureVirtualMachines = pointer.To(input.VirtualMachines)
	}
	if input.Duration != "" {
		updateConfig.Duration = pointer.To(input.Duration)
	}

	if len(input.NonAzureComputerNames) > 0 {
		updateConfig.NonAzureComputerNames = pointer.To(input.NonAzureComputerNames)
	}

	if len(input.Targets) == 1 {
		t := input.Targets[0]
		target := softwareupdateconfiguration.TargetProperties{}
		if len(t.AzureQueries) > 0 {
			aq := make([]softwareupdateconfiguration.AzureQueryProperties, 0)
			for _, v := range t.AzureQueries {
				q := softwareupdateconfiguration.AzureQueryProperties{}
				if len(v.Locations) > 0 {
					q.Locations = pointer.To(v.Locations)
				}
				if len(v.Scope) > 0 {
					q.Scope = pointer.To(v.Scope)
				}
				if len(v.Tags) > 0 || v.TagFilter != "" {
					q.TagSettings = &softwareupdateconfiguration.TagSettingsProperties{
						FilterOperator: pointer.To(softwareupdateconfiguration.TagOperators(v.TagFilter)),
					}
					tags := make(map[string][]string)
					for _, tag := range v.Tags {
						tags[tag.Tag] = tag.Values
					}
					q.TagSettings.Tags = pointer.To(tags)
				}

				aq = append(aq, q)
			}

			target.AzureQueries = pointer.To(aq)
		}

		if len(t.NonAzureQueries) > 0 {
			naqs := make([]softwareupdateconfiguration.NonAzureQueryProperties, 0)
			for _, v := range t.NonAzureQueries {
				naq := softwareupdateconfiguration.NonAzureQueryProperties{}
				if v.FunctionAlias != "" {
					naq.FunctionAlias = pointer.To(v.FunctionAlias)
				}
				if v.WorkspaceId != "" {
					naq.WorkspaceId = pointer.To(v.WorkspaceId)
				}
				naqs = append(naqs, naq)
			}

			target.NonAzureQueries = pointer.To(naqs)
		}
		updateConfig.Targets = pointer.To(target)
	}

	if len(input.Linux) == 1 {
		v := input.Linux[0]
		updateConfig.OperatingSystem = softwareupdateconfiguration.OperatingSystemTypeLinux
		updateConfig.Linux = &softwareupdateconfiguration.LinuxProperties{
			ExcludedPackageNameMasks: pointer.To(v.ExcludedPackages),
			IncludedPackageNameMasks: pointer.To(v.IncludedPackages),
		}
		if v.Reboot != "" {
			updateConfig.Linux.RebootSetting = pointer.To(v.Reboot)
		}
		if !features.FourPointOhBeta() {
			if v.Classification != "" {
				updateConfig.Linux.IncludedPackageClassifications = pointer.To(softwareupdateconfiguration.LinuxUpdateClasses(v.Classification))
			}
		}
		if len(v.Classifications) > 0 {
			updateConfig.Linux.IncludedPackageClassifications = pointer.To(softwareupdateconfiguration.LinuxUpdateClasses(strings.Join(v.Classifications, ",")))
		}
	}

	if len(input.Windows) == 1 {
		v := input.Windows[0]
		updateConfig.OperatingSystem = softwareupdateconfiguration.OperatingSystemTypeWindows
		w := &softwareupdateconfiguration.WindowsProperties{}
		if v.RebootSetting != "" {
			w.RebootSetting = pointer.To(v.RebootSetting)
		}

		if len(v.ExcludedKbs) > 0 {
			w.ExcludedKbNumbers = pointer.To(v.ExcludedKbs)
		}

		if len(v.IncludedKbs) > 0 {
			w.IncludedKbNumbers = pointer.To(v.IncludedKbs)
		}

		if !features.FourPointOhBeta() && len(v.Classification) == 1 {
			w.IncludedUpdateClassifications = pointer.To(softwareupdateconfiguration.WindowsUpdateClasses(strings.Join(v.Classifications, ",")))
		}
		if len(v.Classifications) > 0 {
			w.IncludedUpdateClassifications = pointer.To(softwareupdateconfiguration.WindowsUpdateClasses(strings.Join(v.Classifications, ",")))
		}

		updateConfig.Windows = w
	}

	result.Properties.UpdateConfiguration = updateConfig

	return result
}
