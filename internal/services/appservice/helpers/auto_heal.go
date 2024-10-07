// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
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
	Requests             []AutoHealRequestTrigger      `tfschema:"requests"`
	PrivateMemoryKB      int64                         `tfschema:"private_memory_kb"` // Private should be > 102400 KB (100 MB) to 13631488 KB (13 GB), defaults to 0 however and is always present.
	StatusCodes          []AutoHealStatusCodeTrigger   `tfschema:"status_code"`       // 0 or more, ranges split by `-`, ranges cannot use sub-status or win32 code
	SlowRequests         []AutoHealSlowRequest         `tfschema:"slow_request"`
	SlowRequestsWithPath []AutoHealSlowRequestWithPath `tfschema:"slow_request_with_path"`
}

type AutoHealRequestTrigger struct {
	Count    int64  `tfschema:"count"`
	Interval string `tfschema:"interval"`
}

type AutoHealStatusCodeTrigger struct {
	StatusCodeRange string `tfschema:"status_code_range"` // Conflicts with `StatusCode`, `Win32Code`, and `SubStatus` when not a single value...
	SubStatus       int64  `tfschema:"sub_status"`
	Win32Status     int64  `tfschema:"win32_status_code"`
	Path            string `tfschema:"path"`
	Count           int64  `tfschema:"count"`
	Interval        string `tfschema:"interval"` // Format - hh:mm:ss
}

type AutoHealSlowRequest struct {
	TimeTaken string `tfschema:"time_taken"`
	Interval  string `tfschema:"interval"`
	Count     int64  `tfschema:"count"`
	Path      string `tfschema:"path,removedInNextMajorVersion"`
}

type AutoHealSlowRequestWithPath struct {
	TimeTaken string `tfschema:"time_taken"`
	Interval  string `tfschema:"interval"`
	Count     int64  `tfschema:"count"`
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
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForAutoHealActionType(), false),
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
					Type:     pluginsdk.TypeSet,
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

							"sub_status": {
								Type:     pluginsdk.TypeInt,
								Optional: true,
							},

							"win32_status_code": {
								Type:     pluginsdk.TypeInt,
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
						},
					},
				},

				"slow_request_with_path": {
					Type:     pluginsdk.TypeList,
					Optional: true,
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
		s.Elem.(*pluginsdk.Resource).Schema["slow_request"] = &pluginsdk.Schema{
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
						Deprecated:   "`path` will be removed in `slow_request` and please use `slow_request_with_path` to set the path in version 4.0 of the AzureRM Provider.",
					},
				},
			},
		}
	}
	return s
}

func autoHealTriggerSchemaWindowsComputed() *pluginsdk.Schema {
	s := &pluginsdk.Schema{
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
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"status_code_range": {
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

							"win32_status_code": {
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
						},
					},
				},

				"slow_request_with_path": {
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
	if !features.FourPointOh() {
		s.Elem.(*pluginsdk.Resource).Schema["slow_request"] = &pluginsdk.Schema{
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
						Type:       pluginsdk.TypeString,
						Computed:   true,
						Deprecated: "`path` will be removed in `slow_request` and please use `slow_request_with_path` to set the path in version 4.0 of the AzureRM Provider.",
					},
				},
			},
		}
	}
	return s
}

func expandAutoHealSettingsWindows(autoHealSettings []AutoHealSettingWindows) *webapps.AutoHealRules {
	if len(autoHealSettings) == 0 {
		return &webapps.AutoHealRules{}
	}

	result := &webapps.AutoHealRules{
		Triggers: &webapps.AutoHealTriggers{},
		Actions:  &webapps.AutoHealActions{},
	}

	autoHeal := autoHealSettings[0]
	if len(autoHeal.Triggers) == 0 {
		return result
	}

	triggers := autoHeal.Triggers[0]
	if len(triggers.Requests) == 1 {
		result.Triggers.Requests = &webapps.RequestsBasedTrigger{
			Count:        pointer.To(triggers.Requests[0].Count),
			TimeInterval: pointer.To(triggers.Requests[0].Interval),
		}
	}

	if len(triggers.SlowRequests) == 1 {
		result.Triggers.SlowRequests = &webapps.SlowRequestsBasedTrigger{
			TimeTaken:    pointer.To(triggers.SlowRequests[0].TimeTaken),
			TimeInterval: pointer.To(triggers.SlowRequests[0].Interval),
			Count:        pointer.To(triggers.SlowRequests[0].Count),
		}
		if !features.FourPointOh() {
			if triggers.SlowRequests[0].Path != "" {
				result.Triggers.SlowRequests.Path = pointer.To(triggers.SlowRequests[0].Path)
			}
		}
	}

	if len(triggers.SlowRequestsWithPath) > 0 {
		slowRequestWithPathTriggers := make([]webapps.SlowRequestsBasedTrigger, 0)
		for _, sr := range triggers.SlowRequestsWithPath {
			trigger := webapps.SlowRequestsBasedTrigger{
				TimeTaken:    pointer.To(sr.TimeTaken),
				TimeInterval: pointer.To(sr.Interval),
				Count:        pointer.To(sr.Count),
			}
			if sr.Path != "" {
				trigger.Path = pointer.To(sr.Path)
			}
			slowRequestWithPathTriggers = append(slowRequestWithPathTriggers, trigger)
		}
		result.Triggers.SlowRequestsWithPath = &slowRequestWithPathTriggers
	}

	if triggers.PrivateMemoryKB != 0 {
		result.Triggers.PrivateBytesInKB = pointer.To(triggers.PrivateMemoryKB)
	}

	if len(triggers.StatusCodes) > 0 {
		statusCodeTriggers := make([]webapps.StatusCodesBasedTrigger, 0)
		statusCodeRangeTriggers := make([]webapps.StatusCodesRangeBasedTrigger, 0)
		for _, s := range triggers.StatusCodes {
			statusCodeTrigger := webapps.StatusCodesBasedTrigger{}
			statusCodeRangeTrigger := webapps.StatusCodesRangeBasedTrigger{}
			parts := strings.Split(s.StatusCodeRange, "-")
			if len(parts) == 2 {
				statusCodeRangeTrigger.StatusCodes = pointer.To(s.StatusCodeRange)
				statusCodeRangeTrigger.Count = pointer.To(s.Count)
				statusCodeRangeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeRangeTrigger.Path = pointer.To(s.Path)
				}
				statusCodeRangeTriggers = append(statusCodeRangeTriggers, statusCodeRangeTrigger)
			} else {
				statusCode, err := strconv.Atoi(s.StatusCodeRange)
				if err == nil {
					statusCodeTrigger.Status = pointer.To(int64(statusCode))
				}
				statusCodeTrigger.Count = pointer.To(s.Count)
				statusCodeTrigger.TimeInterval = pointer.To(s.Interval)
				if s.Path != "" {
					statusCodeTrigger.Path = pointer.To(s.Path)
				}
				if s.SubStatus != 0 {
					statusCodeTrigger.SubStatus = pointer.To(s.SubStatus)
				}
				if s.Win32Status != 0 {
					statusCodeTrigger.Win32Status = pointer.To((s.Win32Status))
				}
				statusCodeTriggers = append(statusCodeTriggers, statusCodeTrigger)
			}
		}
		result.Triggers.StatusCodes = &statusCodeTriggers
		result.Triggers.StatusCodesRange = &statusCodeRangeTriggers
	}

	if len(autoHeal.Actions) > 0 {
		action := autoHeal.Actions[0]
		result.Actions.ActionType = pointer.To(webapps.AutoHealActionType(action.ActionType))
		result.Actions.MinProcessExecutionTime = pointer.To(action.MinimumProcessTime)
		if len(action.CustomAction) != 0 {
			customAction := action.CustomAction[0]
			result.Actions.CustomAction = &webapps.AutoHealCustomAction{
				Exe:        pointer.To(customAction.Executable),
				Parameters: pointer.To(customAction.Parameters),
			}
		}
	}
	return result
}

func flattenAutoHealSettingsWindows(autoHealRules *webapps.AutoHealRules) []AutoHealSettingWindows {
	if autoHealRules == nil {
		return []AutoHealSettingWindows{}
	}

	result := AutoHealSettingWindows{}
	// Triggers
	if autoHealRules.Triggers != nil {
		resultTrigger := AutoHealTriggerWindows{}
		triggers := *autoHealRules.Triggers
		if triggers.Requests != nil {
			resultTrigger.Requests = []AutoHealRequestTrigger{{
				Count:    pointer.From(triggers.Requests.Count),
				Interval: pointer.From(triggers.Requests.TimeInterval),
			}}
		}

		if privateBytes := triggers.PrivateBytesInKB; privateBytes != nil && *privateBytes != 0 {
			resultTrigger.PrivateMemoryKB = *triggers.PrivateBytesInKB
		}

		statusCodeTriggers := make([]AutoHealStatusCodeTrigger, 0)
		if triggers.StatusCodes != nil {
			for _, s := range *triggers.StatusCodes {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}

				if s.Status != nil {
					t.StatusCodeRange = strconv.FormatInt(*s.Status, 10)
				}

				if s.Count != nil {
					t.Count = *s.Count
				}

				if s.SubStatus != nil {
					t.SubStatus = *s.SubStatus
				}

				if s.Win32Status != nil {
					t.Win32Status = pointer.From(s.Win32Status)
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		if triggers.StatusCodesRange != nil {
			for _, s := range *triggers.StatusCodesRange {
				t := AutoHealStatusCodeTrigger{
					Interval: pointer.From(s.TimeInterval),
					Path:     pointer.From(s.Path),
				}
				if s.Count != nil {
					t.Count = *s.Count
				}

				if s.StatusCodes != nil {
					t.StatusCodeRange = *s.StatusCodes
				}
				statusCodeTriggers = append(statusCodeTriggers, t)
			}
		}
		resultTrigger.StatusCodes = statusCodeTriggers

		slowRequestTriggers := make([]AutoHealSlowRequest, 0)
		if triggers.SlowRequests != nil {
			slowRequestTriggers = append(slowRequestTriggers, AutoHealSlowRequest{
				TimeTaken: pointer.From(triggers.SlowRequests.TimeTaken),
				Interval:  pointer.From(triggers.SlowRequests.TimeInterval),
				Count:     pointer.From(triggers.SlowRequests.Count),
				Path:      pointer.From(triggers.SlowRequests.Path),
			})
		}

		slowRequestTriggersWithPaths := make([]AutoHealSlowRequestWithPath, 0)
		if triggers.SlowRequestsWithPath != nil {
			for _, v := range *triggers.SlowRequestsWithPath {
				sr := AutoHealSlowRequestWithPath{
					TimeTaken: pointer.From(v.TimeTaken),
					Interval:  pointer.From(v.TimeInterval),
					Count:     pointer.From(v.Count),
					Path:      pointer.From(v.Path),
				}
				slowRequestTriggersWithPaths = append(slowRequestTriggersWithPaths, sr)
			}
		}

		resultTrigger.SlowRequests = slowRequestTriggers
		resultTrigger.SlowRequestsWithPath = slowRequestTriggersWithPaths
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
			ActionType:         string(pointer.From(actions.ActionType)),
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
