package helpers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-03-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AutoHealSettingWindows struct {
	Triggers []AutoHealTriggerWindows `tfschema:"trigger"`
	Actions  []AutoHealActionWindows  `tfschema:"action"`
}

type AutoHealTriggerWindows struct {
	Requests         []AutoHealRequestTrigger         `tfschema:"requests"`
	PrivateMemoryKB  int                              `tfschema:"private_memory_kb"` // Private should be > 102400 KB (100 MB) to 13631488 KB (13 GB), defaults to 0 however and is always present.
	StatusCodes      []AutoHealStatusCodeTrigger      `tfschema:"status_code"`       // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	StatusCodesRange []AutoHealStatusCodeRangeTrigger `tfschema:"status_code_range"` // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	SlowRequests     []AutoHealSlowRequest            `tfschema:"slow_request"`
}

type AutoHealRequestTrigger struct {
	Count    int    `tfschema:"count"`
	Interval string `tfschema:"interval"`
}

type AutoHealStatusCodeTrigger struct {
	StatusCode               string `tfschema:"status_code"`
	StatusCodeRange3Provider string `tfschema:"status_code_range"`
	SubStatus                int    `tfschema:"sub_status"`
	Win32Status              string `tfschema:"win32_status"`
	Path                     string `tfschema:"path"`
	Count                    int    `tfschema:"count"`
	Interval                 string `tfschema:"interval"` // Format - hh:mm:ss
}

type AutoHealStatusCodeRangeTrigger struct {
	StatusCodeRange string `tfschema:"status_code_range"` // Conflicts with `StatusCode`, `Win32Code`, and `SubStatus` when not a single value...
	Path            string `tfschema:"path"`
	Count           int    `tfschema:"count"`
	Interval        string `tfschema:"interval"` // Format - hh:mm:ss
}

type AutoHealSlowRequest struct {
	TimeTaken string `tfschema:"time_taken"`
	Interval  string `tfschema:"interval"`
	Count     int    `tfschema:"count"`
	Path      string `tfschema:"path"`
}

type AutoHealActionWindows struct {
	ActionType         string                 `tfschema:"action_type"`                    // Enum
	CustomAction       []AutoHealCustomAction `tfschema:"custom_action"`                  // Max: 1, needs `action_type` to be "Custom"
	MinimumProcessTime string                 `tfschema:"minimum_process_execution_time"` // Minimum uptime for process before action will trigger
}

type AutoHealCustomAction struct {
	Executable string `tfschema:"executable"`
	Parameters string `tfschema:"parameters"`
}

func autoHealSettingSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaWindows(),

				"action": autoHealActionSchemaWindows(),
			},
		},
		RequiredWith: []string{
			"site_config.0.auto_heal_enabled",
		},
	}
}

func autoHealSettingSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"trigger": autoHealTriggerSchemaWindowsComputed(),

				"action": autoHealActionSchemaWindowsComputed(),
			},
		},
	}
}

func autoHealActionSchemaWindows() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(web.AutoHealActionTypeCustomAction),
						string(web.AutoHealActionTypeLogEvent),
						string(web.AutoHealActionTypeRecycle),
					}, false),
				},

				"custom_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"executable": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"parameters": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					// ValidateFunc: // TODO - Time in hh:mm:ss, because why not...
				},
			},
		},
	}
}

func autoHealActionSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"action_type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"custom_action": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"executable": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"parameters": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"minimum_process_execution_time": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// (@jackofallops) - trigger schemas intentionally left long-hand for now
func autoHealTriggerSchemaWindows() *pluginsdk.Schema {
	s := &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.TimeInterval, // TODO should be hh:mm:ss - This is too loose, need to improve
							},
						},
					},
				},

				"private_memory_kb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(102400, 13631488),
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.TimeInterval,
							},

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"win32_status": {
								Type:     pluginsdk.TypeString,
								Optional: true,
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"status_code_range": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.StatusCodeRange,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"interval": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.TimeInterval,
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.TimeInterval,
							},

							"interval": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validate.TimeInterval,
							},

							"count": {
								Type:         pluginsdk.TypeInt,
								Required:     true,
								ValidateFunc: validation.IntAtLeast(1),
							},

							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}

	if !features.FourPointOhBeta() {
		s.Elem.(*pluginsdk.Resource).Schema["status_code"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// to mak
					"status_code": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"status_code_range": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.StatusCodeRange,
					},

					"count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"interval": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.TimeInterval,
					},

					"sub_status": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"win32_status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		}
	}
	return s
}

func autoHealTriggerSchemaWindowsComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"requests": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"private_memory_kb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"status_code": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"status_code": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"win32_status": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},

				"slow_request": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"time_taken": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"interval": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},

							"count": {
								Type:     pluginsdk.TypeInt,
								Computed: true,
							},

							"path": {
								Type:     pluginsdk.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func expandAutoHealSettingsWindows(autoHealSettings []AutoHealSettingWindows) (*web.AutoHealRules, error) {
	if len(autoHealSettings) == 0 {
		return &web.AutoHealRules{}, nil
	}

	result := &web.AutoHealRules{
		Triggers: &web.AutoHealTriggers{},
		Actions:  &web.AutoHealActions{},
	}

	autoHeal := autoHealSettings[0]

	triggers := autoHeal.Triggers[0]
	if len(triggers.Requests) == 1 {
		result.Triggers.Requests = &web.RequestsBasedTrigger{
			Count:        pointer.To(int32(triggers.Requests[0].Count)),
			TimeInterval: pointer.To(triggers.Requests[0].Interval),
		}
	}

	if len(triggers.SlowRequests) == 1 {
		result.Triggers.SlowRequests = &web.SlowRequestsBasedTrigger{
			TimeTaken:    pointer.To(triggers.SlowRequests[0].TimeTaken),
			TimeInterval: pointer.To(triggers.SlowRequests[0].Interval),
			Count:        pointer.To(int32(triggers.SlowRequests[0].Count)),
		}
		if triggers.SlowRequests[0].Path != "" {
			result.Triggers.SlowRequests.Path = pointer.To(triggers.SlowRequests[0].Path)
		}
	}

	if triggers.PrivateMemoryKB != 0 {
		result.Triggers.PrivateBytesInKB = pointer.To(int32(triggers.PrivateMemoryKB))
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]web.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers3Provider := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := web.StatusCodesBasedTrigger{}
			statusCodeRangeTriggerFor3Provider := web.StatusCodesRangeBasedTrigger{}
			if !features.FourPointOhBeta() {
				if len(strings.Split(s.StatusCodeRange3Provider, "-")) == 2 {
					statusCodeRangeTriggerFor3Provider.StatusCodes = pointer.To(s.StatusCodeRange3Provider)
					statusCodeRangeTriggerFor3Provider.Count = pointer.To(int32(s.Count))
					statusCodeRangeTriggerFor3Provider.TimeInterval = pointer.To(s.Interval)
					if s.Win32Status != "" || s.SubStatus != 0 {
						return nil, fmt.Errorf("`sub_status` and `win32_status` is not supported when setting the status code range")
					}
					if s.Path != "" {
						statusCodeRangeTriggerFor3Provider.Path = pointer.To(s.Path)
					}
					statusCodeRangeTriggers3Provider = append(statusCodeRangeTriggers3Provider, statusCodeRangeTriggerFor3Provider)
				} else {
					statusCode, err := strconv.Atoi(s.StatusCodeRange3Provider)
					if err == nil {
						statusCodeTrigger.Status = pointer.To(int32(statusCode))
					}
					statusCodeTrigger.Count = pointer.To(int32(s.Count))
					statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
					statusCodeTrigger.SubStatus = pointer.To(int32(s.SubStatus))
					win32status, err := strconv.Atoi(s.Win32Status)
					if err == nil {
						statusCodeTrigger.Win32Status = pointer.To(int32(win32status))
					}
					if s.Path != "" {
						statusCodeTrigger.Path = pointer.To(s.Path)
					}
					statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
				}
			} else {
				statusCode, err := strconv.Atoi(s.StatusCode)
				if err == nil {
					statusCodeTrigger.Status = pointer.To(int32(statusCode))
				}
				statusCodeTrigger.Count = pointer.To(int32(s.Count))
				statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
				statusCodeTrigger.SubStatus = pointer.To(int32(s.SubStatus))
				win32status, err := strconv.Atoi(s.Win32Status)
				if err == nil {
					statusCodeTrigger.Win32Status = pointer.To(int32(win32status))
				}
				if s.Path != "" {
					statusCodeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}

		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers3Provider
	}

	if len(triggers.StatusCodesRange) > 0 && features.FourPointOhBeta() {
		statusCodeRangeTriggers := make([]web.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodesRange {
			statusCodeRangeTrigger := web.StatusCodesRangeBasedTrigger{}
			statusCodeRangeTrigger.StatusCodes = pointer.To(s.StatusCodeRange)
			statusCodeRangeTrigger.Count = pointer.To(int32(s.Count))
			statusCodeRangeTrigger.TimeInterval = pointer.To(s.Interval)
			if s.Path != "" {
				statusCodeRangeTrigger.Path = pointer.To(s.Path)
			}
			statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
		}
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	action := autoHeal.Actions[0]
	result.Actions.ActionType = web.AutoHealActionType(action.ActionType)
	result.Actions.MinProcessExecutionTime = pointer.To(action.MinimumProcessTime)
	if len(action.CustomAction) != 0 {
		customAction := action.CustomAction[0]
		result.Actions.CustomAction = &web.AutoHealCustomAction{
			Exe:        pointer.To(customAction.Executable),
			Parameters: pointer.To(customAction.Parameters),
		}
	}

	return result, nil
}

func flattenAutoHealSettingsWindows(autoHealRules *web.AutoHealRules) []AutoHealSettingWindows {
	if autoHealRules == nil {
		return []AutoHealSettingWindows{}
	}

	result := AutoHealSettingWindows{}
	// Triggers
	if autoHealRules.Triggers != nil {
		resultTrigger := AutoHealTriggerWindows{}
		triggers := *autoHealRules.Triggers
		if triggers.Requests != nil {
			count := 0
			if triggers.Requests.Count != nil {
				count = int(*triggers.Requests.Count)
			}
			resultTrigger.Requests = []AutoHealRequestTrigger{{
				Count:    count,
				Interval: pointer.From(triggers.Requests.TimeInterval),
			}}
		}

		if privateBytes := triggers.PrivateBytesInKB; privateBytes != nil && *privateBytes != 0 {
			resultTrigger.PrivateMemoryKB = int(*triggers.PrivateBytesInKB)
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		statusCodeRangeTriggers := make([]AutoHealStatusCodeRangeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}

				if s.Status != nil {
					if !features.FourPointOhBeta() {
						t.StatusCodeRange3Provider = strconv.Itoa(int(*s.Status))
					} else {
						t.StatusCode = strconv.Itoa(int(*s.Status))
					}
				}

				if s.Count != nil {
					t.Count = int(*s.Count)
				}

				if s.SubStatus != nil {
					t.SubStatus = int(*s.SubStatus)
				}

				if s.Win32Status != nil {
					t.Win32Status = strconv.Itoa(int(*s.Win32Status))
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		if triggers.StatusCodesRange != nil {
			for _, s := range *triggers.StatusCodesRange {
				if !features.FourPointOhBeta() {
					tLegacy := AutoHealStatusCodeTrigger{
						Interval: pointer.From(s.TimeInterval),
						Path:     pointer.From(s.Path),
					}
					if s.Count != nil {
						tLegacy.Count = int(*s.Count)
					}

					if s.StatusCodes != nil {
						tLegacy.StatusCodeRange3Provider = *s.StatusCodes
					}
					statusCodeTriggers = append(statusCodeTriggers, tLegacy)
				} else {
					t := AutoHealStatusCodeRangeTrigger{
						Interval: pointer.From(s.TimeInterval),
						Path:     pointer.From(s.Path),
					}
					if s.Count != nil {
						t.Count = int(*s.Count)
					}

					if s.StatusCodes != nil {
						t.StatusCodeRange = *s.StatusCodes
					}
					statusCodeRangeTriggers = append(statusCodeRangeTriggers, t)
				}
			}
		}
		resultTrigger.StatusCodes = statusCodeTriggers

		slowRequestTriggers := make([]AutoHealSlowRequest, 0)
		if triggers.SlowRequests != nil {
			slowRequestTriggers = append(slowRequestTriggers, AutoHealSlowRequest{
				TimeTaken: pointer.From(triggers.SlowRequests.TimeTaken),
				Interval:  pointer.From(triggers.SlowRequests.TimeInterval),
				Count:     int(pointer.From(triggers.SlowRequests.Count)),
				Path:      pointer.From(triggers.SlowRequests.Path),
			})
		}
		resultTrigger.SlowRequests = slowRequestTriggers
		result.Triggers = []AutoHealTriggerWindows{resultTrigger}
	}

	// Actions
	if autoHealRules.Actions != nil {
		actions := *autoHealRules.Actions
		customActions := make([]AutoHealCustomAction, 0)
		if actions.CustomAction != nil {
			customActions = append(customActions, AutoHealCustomAction{
				Executable: pointer.From(actions.CustomAction.Exe),
				Parameters: pointer.From(actions.CustomAction.Parameters),
			})
		}

		resultActions := AutoHealActionWindows{
			ActionType:         string(actions.ActionType),
			CustomAction:       customActions,
			MinimumProcessTime: pointer.From(actions.MinProcessExecutionTime),
		}
		result.Actions = []AutoHealActionWindows{resultActions}
	}

	if result.Actions != nil || result.Triggers != nil {
		return []AutoHealSettingWindows{result}
	}

	return nil
}
