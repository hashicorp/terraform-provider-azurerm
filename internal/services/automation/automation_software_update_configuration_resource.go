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
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/automationaccount"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validate4 "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
func (a *AzureQuery) ToSDKTags() *map[string][]string {
	m := map[string][]string{}
	if len(a.Tags) == 0 {
		// return an empty map instead of nil until issue fixed: https://github.com/Azure/azure-rest-api-specs/issues/21719
		return &m
	}
	for _, tag := range a.Tags {
		m[tag.Tag] = tag.Values
	}
	return &m
}

type Linux struct {
	Reboot           string   `tfschema:"reboot"`
	Classification   string   `tfschema:"classification_included"`
	ExcludedPackages []string `tfschema:"excluded_packages"`
	IncludedPackages []string `tfschema:"included_packages"`
}

type MonthlyOccurrence struct {
	Occurrence int64  `tfschema:"occurrence"`
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

func updateTaskFromSDK(prop *softwareupdateconfiguration.TaskProperties) (res []UpdateTask) {
	if prop == nil {
		return
	}
	res = append(res, UpdateTask{
		Source: utils.NormalizeNilableString(prop.Source),
	})

	if prop.Parameters != nil {
		res[0].Parameters = map[string]string{}
		for k, v := range *prop.Parameters {
			res[0].Parameters[k] = v
		}
	}
	return
}

func (u *UpdateTask) ToSDKModel() *softwareupdateconfiguration.TaskProperties {
	if u == nil {
		return nil
	}
	res := &softwareupdateconfiguration.TaskProperties{
		Source: utils.String(u.Source),
	}

	params := make(map[string]string)
	for k, v := range u.Parameters {
		vCopy := v
		params[k] = vCopy
	}
	res.Parameters = &params
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

func (s *Schedule) LoadSDKModel(info *softwareupdateconfiguration.SUCScheduleProperties) {

	startTime := ""
	if info.StartTime != nil {
		startTime = *info.StartTime
	}

	expiryTime := ""
	if info.ExpiryTime != nil {
		expiryTime = *info.ExpiryTime
	}

	nextRun := ""
	if info.NextRun != nil {
		nextRun = *info.NextRun
	}

	creationTime := ""
	if info.CreationTime != nil {
		creationTime = *info.CreationTime
	}

	lastModifiedTime := ""
	if info.LastModifiedTime != nil {
		lastModifiedTime = *info.LastModifiedTime
	}

	s.StartTime = startTime
	s.StartTimeOffsetMinutes = pointer.ToFloat64(info.StartTimeOffsetMinutes)
	s.ExpiryTime = expiryTime
	s.ExpiryTimeOffsetMinutes = pointer.ToFloat64(info.ExpiryTimeOffsetMinutes)
	s.IsEnabled = utils.NormaliseNilableBool(info.IsEnabled)
	s.NextRun = nextRun
	s.NextRunOffsetMinutes = pointer.ToFloat64(info.NextRunOffsetMinutes)
	if info.Interval != nil {
		s.Interval = int(*info.Interval)
	}
	if info.Frequency != nil {
		s.Frequency = string(*info.Frequency)
	}

	s.TimeZone = utils.NormalizeNilableString(info.TimeZone)
	s.CreationTime = creationTime
	s.LastModifiedTime = lastModifiedTime
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

				day := ""
				if occ.Day != nil {
					day = string(*occ.Day)
				}

				s.MonthlyOccurrence = append(s.MonthlyOccurrence, MonthlyOccurrence{
					Occurrence: utils.NormaliseNilableInt64(occ.Occurrence),
					Day:        day,
				})
			}
		}
	}
}

// will keep old values encode from config
func scheduleFromSDK(info softwareupdateconfiguration.SUCScheduleProperties, old []Schedule) []Schedule {
	if len(old) == 0 {
		old = append(old, Schedule{})
	}
	old[0].LoadSDKModel(&info)

	return old
}

func (s *Schedule) ToSDKModel() softwareupdateconfiguration.SUCScheduleProperties {
	if s == nil {
		return softwareupdateconfiguration.SUCScheduleProperties{}
	}
	timeZone := s.TimeZone

	frequency := softwareupdateconfiguration.ScheduleFrequency(s.Frequency)
	res := softwareupdateconfiguration.SUCScheduleProperties{
		StartTimeOffsetMinutes:  utils.Float(s.StartTimeOffsetMinutes),
		ExpiryTimeOffsetMinutes: utils.Float(s.ExpiryTimeOffsetMinutes),
		IsEnabled:               utils.Bool(s.IsEnabled),
		NextRun:                 &s.NextRun,
		NextRunOffsetMinutes:    utils.Float(s.NextRunOffsetMinutes),
		Interval:                utils.Int64(int64(s.Interval)),
		TimeZone:                &timeZone,
		AdvancedSchedule:        &softwareupdateconfiguration.AdvancedSchedule{},
		Description:             utils.String(s.Description),
		Frequency:               &frequency,
	}

	loc, _ := time.LoadLocation(timeZone)
	if s.StartTime != "" {
		startTime, _ := time.Parse(time.RFC3339, s.StartTime)
		res.SetStartTimeAsTime(startTime.In(loc))
	} else {
		res.StartTime = &s.StartTime
	}
	if s.ExpiryTime != "" {
		expiryTime, _ := time.Parse(time.RFC3339, s.ExpiryTime)
		res.SetExpiryTimeAsTime(expiryTime.In(loc))
	} else {
		res.ExpiryTime = &s.ExpiryTime
	}

	if len(s.AdvancedWeekDays) > 0 {
		res.AdvancedSchedule.WeekDays = &s.AdvancedWeekDays
	}

	if len(s.AdvancedMonthDays) > 0 {
		var is []int64
		for _, v := range s.AdvancedMonthDays {
			is = append(is, int64(v))
		}
		res.AdvancedSchedule.MonthDays = &is
	}

	var occ []softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence
	for _, m := range s.MonthlyOccurrence {

		day := softwareupdateconfiguration.ScheduleDay(m.Day)
		occ = append(occ, softwareupdateconfiguration.AdvancedScheduleMonthlyOccurrence{
			Occurrence: utils.Int64(m.Occurrence),
			Day:        &day,
		})
	}

	if len(occ) > 0 {
		res.AdvancedSchedule.MonthlyOccurrences = &occ
	}
	return res
}

type Target struct {
	AzureQueries    []AzureQuery    `tfschema:"azure_query"`
	NonAzureQueries []NonAzureQuery `tfschema:"non_azure_query"`
}

func targetsFromSDK(prop *softwareupdateconfiguration.TargetProperties) []Target {
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
				if setting.Tags != nil {
					q.LoadSDKTags(*setting.Tags)
				}

				if setting.FilterOperator != nil {
					q.TagFilter = string(*setting.FilterOperator)
				}
			}
			t.AzureQueries = append(t.AzureQueries, q)
		}
	}

	if prop.NonAzureQueries != nil {
		for _, az := range *prop.NonAzureQueries {
			q := NonAzureQuery{
				FunctionAlias: utils.NormalizeNilableString(az.FunctionAlias),
				WorkspaceId:   utils.NormalizeNilableString(az.WorkspaceId),
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

func (s *SoftwareUpdateConfigurationModel) ToSDKModel() softwareupdateconfiguration.SoftwareUpdateConfiguration {
	var param softwareupdateconfiguration.SoftwareUpdateConfiguration
	param.Name = utils.String(s.Name)
	upd := softwareupdateconfiguration.UpdateConfiguration{}
	upd.OperatingSystem = softwareupdateconfiguration.OperatingSystemType(s.OperatingSystem)

	if len(s.Linux) > 0 {
		l := s.Linux[0]
		linuxUpdateClasses := softwareupdateconfiguration.LinuxUpdateClasses(l.Classification)
		upd.Linux = &softwareupdateconfiguration.LinuxProperties{
			IncludedPackageClassifications: &linuxUpdateClasses,
		}
		if l.Reboot != "" {
			upd.Linux.RebootSetting = utils.String(l.Reboot)
		}

		upd.Linux.IncludedPackageNameMasks = utils.StringSlice(l.IncludedPackages)
		upd.Linux.ExcludedPackageNameMasks = utils.StringSlice(l.ExcludedPackages)
	}

	if len(s.Windows) > 0 {
		w := s.Windows[0]
		includedUpdateClassifications := softwareupdateconfiguration.WindowsUpdateClasses(strings.Join(w.Classifications, ","))
		upd.Windows = &softwareupdateconfiguration.WindowsProperties{
			IncludedUpdateClassifications: &includedUpdateClassifications,
		}
		if len(w.Classifications) == 0 && w.Classification != "" {
			includedUpdateClassifications = softwareupdateconfiguration.WindowsUpdateClasses(w.Classification)
			upd.Windows.IncludedUpdateClassifications = &includedUpdateClassifications
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
		upd.Targets = &softwareupdateconfiguration.TargetProperties{}
		var azureQueries []softwareupdateconfiguration.AzureQueryProperties
		t := s.Targets[0]
		for _, az := range t.AzureQueries {
			q := softwareupdateconfiguration.AzureQueryProperties{
				Scope:       utils.StringSlice(az.Scope),
				Locations:   utils.StringSlice(az.Locations),
				TagSettings: nil,
			}
			tag := softwareupdateconfiguration.TagSettingsProperties{}
			tag.Tags = az.ToSDKTags()
			// always set filterOperator until issue fixed: https://github.com/Azure/azure-rest-api-specs/issues/21719
			filterOperator := softwareupdateconfiguration.TagOperatorsAll
			tag.FilterOperator = &filterOperator
			if az.TagFilter != "" {
				tagOperators := softwareupdateconfiguration.TagOperators(az.TagFilter)
				tag.FilterOperator = &tagOperators
			}
			q.TagSettings = &tag
			azureQueries = append(azureQueries, q)
		}

		if azureQueries != nil {
			upd.Targets.AzureQueries = &azureQueries
		}

		var nonAzureQueries []softwareupdateconfiguration.NonAzureQueryProperties
		for _, az := range t.NonAzureQueries {
			q := softwareupdateconfiguration.NonAzureQueryProperties{
				FunctionAlias: utils.String(az.FunctionAlias),
				WorkspaceId:   utils.String(az.WorkspaceId),
			}
			nonAzureQueries = append(nonAzureQueries, q)
		}

		if nonAzureQueries != nil {
			upd.Targets.NonAzureQueries = &nonAzureQueries
		}
	}

	prop := softwareupdateconfiguration.SoftwareUpdateConfigurationProperties{}

	if len(s.Schedule) > 0 {
		prop.ScheduleInfo = s.Schedule[0].ToSDKModel()
	}

	tasks := softwareupdateconfiguration.SoftwareUpdateConfigurationTasks{}
	if len(s.PreTask) > 0 {
		tasks.PreTask = s.PreTask[0].ToSDKModel()
	}
	if len(s.PostTask) > 0 {
		tasks.PostTask = s.PostTask[0].ToSDKModel()
	}
	prop.Tasks = &tasks
	prop.UpdateConfiguration = upd
	param.Properties = prop
	return param
}

func (s *SoftwareUpdateConfigurationModel) LoadSDKModel(model *softwareupdateconfiguration.SoftwareUpdateConfiguration) {
	if model == nil {
		return
	}
	props := model.Properties

	if props.Error != nil {
		s.ErrorCode = utils.NormalizeNilableString(props.Error.Code)
		s.ErrorMeesage = utils.NormalizeNilableString(props.Error.Message)
		s.ErrorMessage = utils.NormalizeNilableString(props.Error.Message)
	}

	conf := props.UpdateConfiguration
	s.OperatingSystem = string(conf.OperatingSystem)

	if l := conf.Linux; l != nil {

		includedPackageClassifications := ""
		if l.IncludedPackageClassifications != nil {
			includedPackageClassifications = string(*l.IncludedPackageClassifications)
		}

		s.Linux = []Linux{{
			Reboot:           utils.NormalizeNilableString(l.RebootSetting),
			Classification:   includedPackageClassifications,
			ExcludedPackages: pointer.ToSliceOfStrings(l.ExcludedPackageNameMasks),
			IncludedPackages: pointer.ToSliceOfStrings(l.IncludedPackageNameMasks),
		}}
	}

	if w := conf.Windows; w != nil {

		classification := ""
		if w.IncludedUpdateClassifications != nil {
			classification = string(*w.IncludedUpdateClassifications)
		}

		s.Windows = []Windows{
			{
				Classification: classification,
				ExcludedKbs:    pointer.ToSliceOfStrings(w.ExcludedKbNumbers),
				IncludedKbs:    pointer.ToSliceOfStrings(w.IncludedKbNumbers),
				RebootSetting:  utils.NormalizeNilableString(w.RebootSetting),
			},
		}

		for _, v := range strings.Split(classification, ",") {
			s.Windows[0].Classifications = append(s.Windows[0].Classifications, strings.TrimSpace(v))
		}
	}

	s.Duration = utils.NormalizeNilableString(conf.Duration)
	s.VirtualMachines = pointer.ToSliceOfStrings(conf.AzureVirtualMachines)
	s.NonAzureComputerNames = pointer.ToSliceOfStrings(conf.NonAzureComputerNames)
	s.Targets = targetsFromSDK(conf.Targets)

	// service api response scheduleInfo.advancedSchedule as null, which cause import lost it
	s.Schedule = scheduleFromSDK(props.ScheduleInfo, s.Schedule)
	if tasks := props.Tasks; tasks != nil {
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
			ValidateFunc: softwareupdateconfiguration.ValidateAutomationAccountID,
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
				string(softwareupdateconfiguration.OperatingSystemTypeLinux),
				string(softwareupdateconfiguration.OperatingSystemTypeWindows),
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
							vs = append(vs, softwareupdateconfiguration.PossibleValuesForLinuxUpdateClasses()...)
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
							vs = append(vs, softwareupdateconfiguration.PossibleValuesForWindowsUpdateClasses()...)
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
						Default:      "Etc/UTC",
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
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				if meta.ResourceData.IsNewResource() {
					return meta.ResourceRequiresImport(m.ResourceType(), id)
				}
			}

			param := model.ToSDKModel()
			_, err = client.SoftwareUpdateConfigurationsCreate(ctx, id, param, softwareupdateconfiguration.SoftwareUpdateConfigurationsCreateOperationOptions{})

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
			var output SoftwareUpdateConfigurationModel
			if err := meta.Decode(&output); err != nil {
				return err
			}

			output.Name = id.SoftwareUpdateConfigurationName
			output.AutomationAccountID = softwareupdateconfiguration.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName).ID()
			output.LoadSDKModel(resp.Model)

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
