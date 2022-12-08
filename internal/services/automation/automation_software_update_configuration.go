package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation" // nolint: staticcheck
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate4 "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	validate2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func (a *AzureQuery) LoadSDKTags(tags map[string][]string) {
	if tags == nil {
		return
	}
	for k, vs := range tags {
		t := Tag{}
		t.Tag = k
		t.Values = append(t.Values, vs...)
		a.Tags = append(a.Tags, t)
	}
}
func (a *AzureQuery) ToSDKTags() map[string][]string {
	m := map[string][]string{}
	if len(a.Tags) == 0 {
		// return an empty map instead of nil until issue fixed: https://github.com/Azure/azure-rest-api-specs/issues/21719
		return m
	}
	for _, tag := range a.Tags {
		m[tag.Tag] = tag.Values
	}
	return m
}

type Linux struct {
	Reboot           string   `tfschema:"reboot"`
	Classification   string   `tfschema:"classification_included"`
	ExcludedPackages []string `tfschema:"excluded_packages"`
	IncludedPackages []string `tfschema:"included_packages"`
}

type MonthlyOccurrence struct {
	Occurrence int32  `tfschema:"occurrence"`
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

func updateTaskFromSDK(prop *automation.TaskProperties) (res []UpdateTask) {
	if prop == nil {
		return
	}
	res = append(res, UpdateTask{
		Source: utils.NormalizeNilableString(prop.Source),
	})
	for k, v := range prop.Parameters {
		res[0].Parameters[k] = *v
	}
	return
}

func (u *UpdateTask) ToSDKModel() *automation.TaskProperties {
	if u == nil {
		return nil
	}
	res := &automation.TaskProperties{
		Parameters: map[string]*string{},
		Source:     utils.String(u.Source),
	}
	for k, v := range u.Parameters {
		vCopy := v
		res.Parameters[k] = &vCopy
	}
	return res
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

func (s *Schedule) LoadSDKModel(info *automation.SUCScheduleProperties) {
	timeString := func(t *date.Time) string {
		if t == nil {
			return ""
		}
		return t.Time.Format(time.RFC3339)
	}

	s.StartTime = timeString(info.StartTime)
	s.StartTimeOffsetMinutes = pointer.ToFloat64(info.StartTimeOffsetMinutes)
	s.ExpiryTime = timeString(info.ExpiryTime)
	s.ExpiryTimeOffsetMinutes = pointer.ToFloat64(info.ExpiryTimeOffsetMinutes)
	s.IsEnabled = utils.NormaliseNilableBool(info.IsEnabled)
	s.NextRun = timeString(info.NextRun)
	s.NextRunOffsetMinutes = pointer.ToFloat64(info.NextRunOffsetMinutes)
	if info.Interval != nil {
		s.Interval = int(*info.Interval)
	}
	s.Frequency = string(info.Frequency)
	s.TimeZone = utils.NormalizeNilableString(info.TimeZone)
	s.CreationTime = timeString(info.CreationTime)
	s.LastModifiedTime = timeString(info.LastModifiedTime)
	s.Description = utils.NormalizeNilableString(info.Description)

	if setting := info.AdvancedSchedule; setting != nil {
		s.AdvancedWeekDays = pointer.ToSliceOfStrings(setting.WeekDays)
		if setting.MonthDays != nil {
			for _, v := range *(setting.MonthDays) {
				s.AdvancedMonthDays = append(s.AdvancedMonthDays, int(v))
			}
		}

		if setting.MonthlyOccurrences != nil {
			for _, occ := range *setting.MonthlyOccurrences {
				s.MonthlyOccurrence = append(s.MonthlyOccurrence, MonthlyOccurrence{
					Occurrence: utils.NormaliseNilableInt32(occ.Occurrence),
					Day:        string(occ.Day),
				})
			}
		}
	}
}

// will keep old values encode from config
func scheduleFromSDK(info *automation.SUCScheduleProperties, old []Schedule) []Schedule {
	if info == nil {
		return old
	}
	if len(old) == 0 {
		old = append(old, Schedule{})
	}
	old[0].LoadSDKModel(info)

	return old
}

func (s *Schedule) ToSDKModel() *automation.SUCScheduleProperties {
	if s == nil {
		return nil
	}

	parseTime := func(s string) *date.Time {
		t, _ := time.Parse(time.RFC3339, s)
		return &date.Time{Time: t}
	}

	res := automation.SUCScheduleProperties{
		StartTime:               parseTime(s.StartTime),
		StartTimeOffsetMinutes:  utils.Float(s.StartTimeOffsetMinutes),
		ExpiryTime:              parseTime(s.ExpiryTime),
		ExpiryTimeOffsetMinutes: utils.Float(s.ExpiryTimeOffsetMinutes),
		IsEnabled:               utils.Bool(s.IsEnabled),
		NextRun:                 parseTime(s.NextRun),
		NextRunOffsetMinutes:    utils.Float(s.NextRunOffsetMinutes),
		Interval:                utils.Int64(int64(s.Interval)),
		TimeZone:                utils.String(s.TimeZone),
		AdvancedSchedule:        &automation.AdvancedSchedule{},
		Description:             utils.String(s.Description),
		Frequency:               automation.ScheduleFrequency(s.Frequency),
	}

	if len(s.AdvancedWeekDays) > 0 {
		res.AdvancedSchedule.WeekDays = &s.AdvancedWeekDays
	}

	if len(s.AdvancedMonthDays) > 0 {
		var is []int32
		for _, v := range s.AdvancedMonthDays {
			is = append(is, int32(v))
		}
		res.AdvancedSchedule.MonthDays = &is
	}

	var occ []automation.AdvancedScheduleMonthlyOccurrence
	for _, m := range s.MonthlyOccurrence {
		occ = append(occ, automation.AdvancedScheduleMonthlyOccurrence{
			Occurrence: utils.Int32(m.Occurrence),
			Day:        automation.ScheduleDay(m.Day),
		})
	}

	if len(occ) > 0 {
		res.AdvancedSchedule.MonthlyOccurrences = &occ
	}
	return &res
}

type Target struct {
	AzureQueries    []AzureQuery    `tfschema:"azure_query"`
	NonAzureQueries []NonAzureQuery `tfschema:"non_azure_query"`
}

func targetsFromSDK(prop *automation.TargetProperties) []Target {
	if prop == nil {
		return nil
	}

	var t Target
	if prop.AzureQueries != nil {
		for _, az := range *prop.AzureQueries {
			q := AzureQuery{
				Scope:     pointer.ToSliceOfStrings(az.Scope),
				Locations: pointer.ToSliceOfStrings(az.Locations),
			}
			if setting := az.TagSettings; setting != nil {
				q.LoadSDKTags(setting.Tags)
				q.TagFilter = string(setting.FilterOperator)
			}
			t.AzureQueries = append(t.AzureQueries, q)
		}
	}

	if prop.NonAzureQueries != nil {
		for _, az := range *prop.NonAzureQueries {
			q := NonAzureQuery{
				FunctionAlias: utils.NormalizeNilableString(az.FunctionAlias),
				WorkspaceId:   utils.NormalizeNilableString(az.WorkspaceID),
			}
			t.NonAzureQueries = append(t.NonAzureQueries, q)
		}
	}

	return []Target{t}
}

type Windows struct {
	// Classification Deprecated, use Classifications instead
	Classification string `tfschema:"classification_included"`

	Classifications []string `tfschema:"classifications_included"`
	ExcludedKbs     []string `tfschema:"excluded_knowledge_base_numbers"`
	IncludedKbs     []string `tfschema:"included_knowledge_base_numbers"`
	RebootSetting   string   `tfschema:"reboot"`
}

type SoftwareUpdateConfigurationModel struct {
	AutomationAccountID   string       `tfschema:"automation_account_id"`
	Name                  string       `tfschema:"name"`
	ErrorCode             string       `tfschema:"error_code"`
	ErrorMeesage          string       `tfschema:"error_meesage"`
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

func (s *SoftwareUpdateConfigurationModel) ToSDKModel() automation.SoftwareUpdateConfiguration {
	var param automation.SoftwareUpdateConfiguration
	param.Name = utils.String(s.Name)
	param.SoftwareUpdateConfigurationProperties = &automation.SoftwareUpdateConfigurationProperties{}
	prop := param.SoftwareUpdateConfigurationProperties
	prop.UpdateConfiguration = &automation.UpdateConfiguration{}
	upd := prop.UpdateConfiguration
	upd.OperatingSystem = automation.OperatingSystemType(s.OperatingSystem)

	if len(s.Linux) > 0 {
		l := s.Linux[0]
		upd.Linux = &automation.LinuxProperties{
			IncludedPackageClassifications: automation.LinuxUpdateClasses(l.Classification),
		}
		if l.Reboot != "" {
			upd.Linux.RebootSetting = utils.String(l.Reboot)
		}

		upd.Linux.IncludedPackageNameMasks = utils.StringSlice(l.IncludedPackages)
		upd.Linux.ExcludedPackageNameMasks = utils.StringSlice(l.ExcludedPackages)
	}

	if len(s.Windows) > 0 {
		w := s.Windows[0]

		upd.Windows = &automation.WindowsProperties{
			IncludedUpdateClassifications: automation.WindowsUpdateClasses(strings.Join(w.Classifications, ",")),
		}
		if len(w.Classifications) == 0 && w.Classification != "" {
			upd.Windows.IncludedUpdateClassifications = automation.WindowsUpdateClasses(w.Classification)
		}

		if w.RebootSetting != "" {
			upd.Windows.RebootSetting = utils.String(w.RebootSetting)
		}
		upd.Windows.IncludedKbNumbers = utils.StringSlice(w.IncludedKbs)
		upd.Windows.ExcludedKbNumbers = utils.StringSlice(w.ExcludedKbs)
	}

	upd.Duration = utils.String(s.Duration)
	upd.AzureVirtualMachines = utils.StringSlice(s.VirtualMachines)
	upd.NonAzureComputerNames = utils.StringSlice(s.NonAzureComputerNames)

	if len(s.Targets) > 0 {
		upd.Targets = &automation.TargetProperties{}
		var azureQueries []automation.AzureQueryProperties
		t := s.Targets[0]
		for _, az := range t.AzureQueries {
			q := automation.AzureQueryProperties{
				Scope:       utils.StringSlice(az.Scope),
				Locations:   utils.StringSlice(az.Locations),
				TagSettings: nil,
			}
			tag := automation.TagSettingsProperties{}
			tag.Tags = az.ToSDKTags()
			// always set filterOperator until issue fixed: https://github.com/Azure/azure-rest-api-specs/issues/21719
			tag.FilterOperator = automation.TagOperatorsAll
			if az.TagFilter != "" {
				tag.FilterOperator = automation.TagOperators(az.TagFilter)
			}
			q.TagSettings = &tag
			azureQueries = append(azureQueries, q)
		}

		if azureQueries != nil {
			upd.Targets.AzureQueries = &azureQueries
		}

		var nonAzureQueries []automation.NonAzureQueryProperties
		for _, az := range t.NonAzureQueries {
			q := automation.NonAzureQueryProperties{
				FunctionAlias: utils.String(az.FunctionAlias),
				WorkspaceID:   utils.String(az.WorkspaceId),
			}
			nonAzureQueries = append(nonAzureQueries, q)
		}

		if nonAzureQueries != nil {
			upd.Targets.NonAzureQueries = &nonAzureQueries
		}
	}

	if len(s.Schedule) > 0 {
		prop.ScheduleInfo = s.Schedule[0].ToSDKModel()
	}

	prop.Tasks = &automation.SoftwareUpdateConfigurationTasks{}
	if len(s.PreTask) > 0 {
		prop.Tasks.PreTask = s.PreTask[0].ToSDKModel()
	}
	if len(s.PostTask) > 0 {
		prop.Tasks.PostTask = s.PostTask[0].ToSDKModel()
	}
	return param
}

func (s *SoftwareUpdateConfigurationModel) LoadSDKModel(prop *automation.SoftwareUpdateConfigurationProperties) {
	if prop == nil {
		return
	}

	if prop.Error != nil {
		s.ErrorCode = utils.NormalizeNilableString(prop.Error.Code)
		s.ErrorMeesage = utils.NormalizeNilableString(prop.Error.Message)
	}

	if conf := prop.UpdateConfiguration; conf != nil {
		s.OperatingSystem = string(conf.OperatingSystem)

		if l := conf.Linux; l != nil {
			s.Linux = []Linux{{
				Reboot:           utils.NormalizeNilableString(l.RebootSetting),
				Classification:   string(l.IncludedPackageClassifications),
				ExcludedPackages: pointer.ToSliceOfStrings(l.ExcludedPackageNameMasks),
				IncludedPackages: pointer.ToSliceOfStrings(l.IncludedPackageNameMasks),
			}}
		}

		if w := conf.Windows; w != nil {
			s.Windows = []Windows{
				{
					Classification: string(w.IncludedUpdateClassifications),
					ExcludedKbs:    pointer.ToSliceOfStrings(w.ExcludedKbNumbers),
					IncludedKbs:    pointer.ToSliceOfStrings(w.IncludedKbNumbers),
					RebootSetting:  utils.NormalizeNilableString(w.RebootSetting),
				},
			}

			for _, v := range strings.Split(string(w.IncludedUpdateClassifications), ",") {
				s.Windows[0].Classifications = append(s.Windows[0].Classifications, strings.TrimSpace(v))
			}
		}

		s.Duration = utils.NormalizeNilableString(conf.Duration)
		s.VirtualMachines = pointer.ToSliceOfStrings(conf.AzureVirtualMachines)
		s.NonAzureComputerNames = pointer.ToSliceOfStrings(conf.NonAzureComputerNames)
		s.Targets = targetsFromSDK(conf.Targets)
	}

	// service api response scheduleInfo.advancedSchedule as null, which cause import lost it
	s.Schedule = scheduleFromSDK(prop.ScheduleInfo, s.Schedule)
	if tasks := prop.Tasks; tasks != nil {
		s.PreTask = updateTaskFromSDK(tasks.PreTask)
		s.PostTask = updateTaskFromSDK(tasks.PostTask)
	}
}

type SoftwareUpdateConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = (*SoftwareUpdateConfigurationResource)(nil)

func (m SoftwareUpdateConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{

		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AutomationAccountID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"operating_system": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(automation.OperatingSystemTypeLinux),
				string(automation.OperatingSystemTypeWindows),
			}, false),
		},

		"linux": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{

					"reboot": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"classification_included": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(func() (vs []string) {
							for _, v := range automation.PossibleLinuxUpdateClassesValues() {
								vs = append(vs, string(v))
							}
							return
						}(), false),
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
			},
		},

		"windows": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{

					"classification_included": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"windows.0.classifications_included"},
						AtLeastOneOf:  []string{"windows.0.classification_included", "windows.0.classifications_included"},
						Deprecated:    "windows classification can be set as a list, use `classifications_included` instead.",
						ValidateFunc: validation.StringInSlice(func() (vs []string) {
							for _, v := range automation.PossibleWindowsUpdateClassesValues() {
								vs = append(vs, string(v))
							}
							return
						}(), false),
					},

					"classifications_included": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						Computed:      true,
						ConflictsWith: []string{"windows.0.classification_included"},
						AtLeastOneOf:  []string{"windows.0.classification_included", "windows.0.classifications_included"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(func() (res []string) {
								for _, v := range automation.PossibleWindowsUpdateClassesValues() {
									res = append(res, string(v))
								}
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
					},
				},
			},
		},

		"duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
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
										string(automation.TagOperatorsAny),
										string(automation.TagOperatorsAll),
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
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"start_time": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
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
						Optional: true,
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
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"advanced_week_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: nil,
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
										for _, v := range automation.PossibleScheduleDayValues() {
											vs = append(vs, string(v))
										}
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
}

func (m SoftwareUpdateConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"error_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"error_meesage": {
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
			automationID, _ := parse.AutomationAccountID(model.AutomationAccountID)

			subscriptionID := meta.Client.Account.SubscriptionId
			id := parse.NewSoftwareUpdateConfigurationID(subscriptionID, automationID.ResourceGroup, automationID.Name, model.Name)
			existing, err := client.GetByName(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, "")
			if !utils.ResponseWasNotFound(existing.Response) {
				if err != nil {
					return fmt.Errorf("retreiving %s: %v", id, err)
				}
				if meta.ResourceData.IsNewResource() {
					return meta.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			param := model.ToSDKModel()
			future, err := client.Create(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, param, "")
			_ = future

			if err != nil {
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
			id, err := parse.SoftwareUpdateConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			client := meta.Client.Automation.SoftwareUpdateConfigClient
			result, err := client.GetByName(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, "")
			if err != nil {
				return err
			}
			var output SoftwareUpdateConfigurationModel
			if err := meta.Decode(&output); err != nil {
				return err
			}

			output.Name = id.Name
			output.AutomationAccountID = parse.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName).ID()
			output.LoadSDKModel(result.SoftwareUpdateConfigurationProperties)

			return meta.Encode(&output)
		},
	}
}

func (m SoftwareUpdateConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func:    m.Create().Func,
		Timeout: 10 * time.Minute,
	}
}

func (m SoftwareUpdateConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			id, err := parse.SoftwareUpdateConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			meta.Logger.Infof("deleting %s", id)
			client := meta.Client.Automation.SoftwareUpdateConfigClient
			if _, err = client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, id.Name, ""); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (m SoftwareUpdateConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SoftwareUpdateConfigurationID
}
