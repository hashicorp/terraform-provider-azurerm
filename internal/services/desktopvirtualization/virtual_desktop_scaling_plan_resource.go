// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = DesktopVirtualizationScalingPlanResource{}
	_ sdk.ResourceWithUpdate = DesktopVirtualizationScalingPlanResource{}
)

type DesktopVirtualizationScalingPlanResource struct{}

func (DesktopVirtualizationScalingPlanResource) ModelObject() interface{} {
	return &DesktopVirtualizationScalingPlanResourceModel{}
}

func (DesktopVirtualizationScalingPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scalingplan.ValidateHostPoolID
}

func (DesktopVirtualizationScalingPlanResource) ResourceType() string {
	return "azurerm_virtual_desktop_scaling_plan"
}

type DesktopVirtualizationScalingPlanResourceModel struct {
	Name              string                                     `tfschema:"name"`
	Location          string                                     `tfschema:"location"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	FriendlyName      string                                     `tfschema:"friendly_name"`
	Description       string                                     `tfschema:"description"`
	TimeZone          string                                     `tfschema:"time_zone"`
	ExclusionTag      string                                     `tfschema:"exclusion_tag"`
	Schedule          []DesktopVirtualizationScalingPlanSchedule `tfschema:"schedule"`
	HostPool          []DesktopVirtualizationScalingPlanHostPool `tfschema:"host_pool"`
	Tags              map[string]string                          `tfschema:"tags"`
}

type DesktopVirtualizationScalingPlanSchedule struct {
	Name                             string   `tfschema:"name"`
	DaysOfWeek                       []string `tfschema:"days_of_week"`
	RampUpStartTime                  string   `tfschema:"ramp_up_start_time"`
	RampUpLoadBalancingAlgorithm     string   `tfschema:"ramp_up_load_balancing_algorithm"`
	RampUpMinimumHostsPercent        int64    `tfschema:"ramp_up_minimum_hosts_percent"`
	RampUpCapacityThresholdPercent   int64    `tfschema:"ramp_up_capacity_threshold_percent"`
	PeakStartTime                    string   `tfschema:"peak_start_time"`
	PeakLoadBalancingAlgorithm       string   `tfschema:"peak_load_balancing_algorithm"`
	RampDownStartTime                string   `tfschema:"ramp_down_start_time"`
	RampDownLoadBalancingAlgorithm   string   `tfschema:"ramp_down_load_balancing_algorithm"`
	RampDownMinimumHostsPercent      int64    `tfschema:"ramp_down_minimum_hosts_percent"`
	RampDownCapacityThresholdPercent int64    `tfschema:"ramp_down_capacity_threshold_percent"`
	RampDownForceLogoffUsers         bool     `tfschema:"ramp_down_force_logoff_users"`
	RampDownStopHostsWhen            string   `tfschema:"ramp_down_stop_hosts_when"`
	RampDownWaitTimeMinutes          int64    `tfschema:"ramp_down_wait_time_minutes"`
	RampDownNotificationMessage      string   `tfschema:"ramp_down_notification_message"`
	OffPeakStartTime                 string   `tfschema:"off_peak_start_time"`
	OffPeakLoadBalancingAlgorithm    string   `tfschema:"off_peak_load_balancing_algorithm"`
}

type DesktopVirtualizationScalingPlanHostPool struct {
	HostpoolId         string `tfschema:"hostpool_id"`
	ScalingPlanEnabled bool   `tfschema:"scaling_plan_enabled"`
}

func (r DesktopVirtualizationScalingPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"exclusion_tag": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"days_of_week": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForDaysOfWeek(), false),
						},
					},

					"ramp_up_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_up_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_up_minimum_hosts_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_up_capacity_threshold_percent": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"peak_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"peak_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"ramp_down_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},

					"ramp_down_minimum_hosts_percent": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_capacity_threshold_percent": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},

					"ramp_down_force_logoff_users": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"ramp_down_stop_hosts_when": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForStopHostsWhen(), false),
					},

					"ramp_down_wait_time_minutes": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"ramp_down_notification_message": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"off_peak_start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validateTime(),
					},

					"off_peak_load_balancing_algorithm": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForSessionHostLoadBalancingAlgorithm(), false),
					},
				},
			},
		},

		"host_pool": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"hostpool_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: hostpool.ValidateHostPoolID,
					},
					"scaling_plan_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r DesktopVirtualizationScalingPlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func validateTime() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`), `The time must be in the format HH:MM.`)
}

func (r DesktopVirtualizationScalingPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DesktopVirtualizationScalingPlanResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan create")

			id := scalingplan.NewScalingPlanID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s): %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			hostPoolType := scalingplan.ScalingHostPoolTypePooled // Only one possible value for this
			payload := scalingplan.ScalingPlan{
				Name:     pointer.To(model.Name),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: scalingplan.ScalingPlanProperties{
					Description:        pointer.To(model.Description),
					FriendlyName:       pointer.To(model.FriendlyName),
					TimeZone:           model.TimeZone,
					HostPoolType:       &hostPoolType,
					ExclusionTag:       pointer.To(model.ExclusionTag),
					Schedules:          expandScalingPlanSchedule(model.Schedule),
					HostPoolReferences: expandScalingPlanHostpoolReference(model.HostPool),
				},
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r DesktopVirtualizationScalingPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			state := DesktopVirtualizationScalingPlanResourceModel{}

			id, err := scalingplan.ParseScalingPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.Name = id.ScalingPlanName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Description = pointer.From(model.Properties.Description)
				state.FriendlyName = pointer.From(model.Properties.FriendlyName)
				state.TimeZone = model.Properties.TimeZone
				state.ExclusionTag = pointer.From(model.Properties.ExclusionTag)
				state.Schedule = flattenScalingPlanSchedule(model.Properties.Schedules)
				state.HostPool = flattenScalingHostpoolReference(model.Properties.HostPoolReferences)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationScalingPlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			var model DesktopVirtualizationScalingPlanResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan update")

			id, err := scalingplan.ParseScalingPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			payload := scalingplan.ScalingPlanPatch{
				Tags: pointer.To(model.Tags),
				Properties: &scalingplan.ScalingPlanPatchProperties{
					Description:        pointer.To(model.Description),
					FriendlyName:       pointer.To(model.FriendlyName),
					TimeZone:           pointer.To(model.TimeZone),
					ExclusionTag:       pointer.To(model.ExclusionTag),
					Schedules:          expandScalingPlanSchedule(model.Schedule),
					HostPoolReferences: expandScalingPlanHostpoolReference(model.HostPool),
				},
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationScalingPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			id, err := scalingplan.ParseScalingPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandScalingPlanSchedule(input []DesktopVirtualizationScalingPlanSchedule) *[]scalingplan.ScalingSchedule {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingSchedule, 0)
	for _, item := range input {
		name := item.Name
		daysOfWeekRaw := item.DaysOfWeek
		daysOfWeek := make([]scalingplan.DaysOfWeek, 0)
		for _, weekday := range daysOfWeekRaw {
			daysOfWeek = append(daysOfWeek, scalingplan.DaysOfWeek(weekday))
		}

		results = append(results, scalingplan.ScalingSchedule{
			Name:                           pointer.To(name),
			DaysOfWeek:                     pointer.To(daysOfWeek),
			RampUpStartTime:                expandScalingPlanScheduleTime(item.RampUpStartTime),
			RampUpLoadBalancingAlgorithm:   pointer.To(scalingplan.SessionHostLoadBalancingAlgorithm(item.RampUpLoadBalancingAlgorithm)),
			RampUpMinimumHostsPct:          pointer.To(item.RampUpMinimumHostsPercent),
			RampUpCapacityThresholdPct:     pointer.To(item.RampUpCapacityThresholdPercent),
			PeakStartTime:                  expandScalingPlanScheduleTime(item.PeakStartTime),
			PeakLoadBalancingAlgorithm:     pointer.To(scalingplan.SessionHostLoadBalancingAlgorithm(item.PeakLoadBalancingAlgorithm)),
			RampDownStartTime:              expandScalingPlanScheduleTime(item.RampDownStartTime),
			RampDownLoadBalancingAlgorithm: pointer.To(scalingplan.SessionHostLoadBalancingAlgorithm(item.RampDownLoadBalancingAlgorithm)),
			RampDownMinimumHostsPct:        pointer.To(item.RampDownMinimumHostsPercent),
			RampDownCapacityThresholdPct:   pointer.To(item.RampDownCapacityThresholdPercent),
			RampDownForceLogoffUsers:       pointer.To(item.RampDownForceLogoffUsers),
			RampDownStopHostsWhen:          pointer.To(scalingplan.StopHostsWhen(item.RampDownStopHostsWhen)),
			RampDownWaitTimeMinutes:        pointer.To(item.RampDownWaitTimeMinutes),
			RampDownNotificationMessage:    pointer.To(item.RampDownNotificationMessage),
			OffPeakStartTime:               expandScalingPlanScheduleTime(item.OffPeakStartTime),
			OffPeakLoadBalancingAlgorithm:  pointer.To(scalingplan.SessionHostLoadBalancingAlgorithm(item.OffPeakLoadBalancingAlgorithm)),
		})
	}

	return &results
}

func expandScalingPlanScheduleTime(input string) *scalingplan.Time {
	if len(input) == 0 {
		return nil
	}

	time := strings.Split(input, ":")
	hour, _ := strconv.Atoi(time[0])
	minute, _ := strconv.Atoi(time[1])

	return &scalingplan.Time{
		Hour:   int64(hour),
		Minute: int64(minute),
	}
}

func expandScalingPlanHostpoolReference(input []DesktopVirtualizationScalingPlanHostPool) *[]scalingplan.ScalingHostPoolReference {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingHostPoolReference, 0)
	for _, item := range input {
		hostPoolArmPath := item.HostpoolId
		scalingPlanEnabled := item.ScalingPlanEnabled

		results = append(results, scalingplan.ScalingHostPoolReference{
			HostPoolArmPath:    pointer.To(hostPoolArmPath),
			ScalingPlanEnabled: pointer.To(scalingPlanEnabled),
		})
	}
	return &results
}

func flattenScalingPlanSchedule(input *[]scalingplan.ScalingSchedule) []DesktopVirtualizationScalingPlanSchedule {
	results := make([]DesktopVirtualizationScalingPlanSchedule, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		name := ""
		if item.Name != nil {
			name = *item.Name
		}
		rampUpStartTime := ""
		if item.RampUpStartTime != nil {
			rampUpStartTime = fmt.Sprintf("%02d:%02d", item.RampUpStartTime.Hour, item.RampUpStartTime.Minute)
		}
		rampUpMinimumHostsPct := int64(0)
		if item.RampUpMinimumHostsPct != nil {
			rampUpMinimumHostsPct = *item.RampUpMinimumHostsPct
		}
		rampUpCapacityThresholdPct := int64(0)
		if item.RampUpCapacityThresholdPct != nil {
			rampUpCapacityThresholdPct = *item.RampUpCapacityThresholdPct
		}
		peakStartTime := ""
		if item.PeakStartTime != nil {
			peakStartTime = fmt.Sprintf("%02d:%02d", item.PeakStartTime.Hour, item.PeakStartTime.Minute)
		}
		rampDownStartTime := ""
		if item.RampDownStartTime != nil {
			rampDownStartTime = fmt.Sprintf("%02d:%02d", item.RampDownStartTime.Hour, item.RampDownStartTime.Minute)
		}
		rampDownMinimumHostsPct := int64(0)
		if item.RampDownMinimumHostsPct != nil {
			rampDownMinimumHostsPct = *item.RampDownMinimumHostsPct
		}
		rampDownCapacityThresholdPct := int64(0)
		if item.RampDownCapacityThresholdPct != nil {
			rampDownCapacityThresholdPct = *item.RampDownCapacityThresholdPct
		}
		rampDownForceLogoffUsers := false
		if item.RampDownForceLogoffUsers != nil {
			rampDownForceLogoffUsers = *item.RampDownForceLogoffUsers
		}
		rampDownWaitTimeMinutes := int64(0)
		if item.RampDownWaitTimeMinutes != nil {
			rampDownWaitTimeMinutes = *item.RampDownWaitTimeMinutes
		}
		rampDownNotificationMessage := ""
		if item.RampDownNotificationMessage != nil {
			rampDownNotificationMessage = *item.RampDownNotificationMessage
		}
		offPeakStartTime := ""
		if item.OffPeakStartTime != nil {
			offPeakStartTime = fmt.Sprintf("%02d:%02d", item.OffPeakStartTime.Hour, item.OffPeakStartTime.Minute)
		}
		daysOfWeek := make([]string, 0)
		if item.DaysOfWeek != nil {
			for _, weekday := range *item.DaysOfWeek {
				daysOfWeek = append(daysOfWeek, string(weekday))
			}
		}

		results = append(results, DesktopVirtualizationScalingPlanSchedule{
			Name:                             name,
			DaysOfWeek:                       daysOfWeek,
			RampUpStartTime:                  rampUpStartTime,
			RampUpLoadBalancingAlgorithm:     string(pointer.From(item.RampUpLoadBalancingAlgorithm)),
			RampUpMinimumHostsPercent:        rampUpMinimumHostsPct,
			RampUpCapacityThresholdPercent:   rampUpCapacityThresholdPct,
			PeakStartTime:                    peakStartTime,
			PeakLoadBalancingAlgorithm:       string(pointer.From(item.PeakLoadBalancingAlgorithm)),
			RampDownStartTime:                rampDownStartTime,
			RampDownLoadBalancingAlgorithm:   string(pointer.From(item.RampDownLoadBalancingAlgorithm)),
			RampDownMinimumHostsPercent:      rampDownMinimumHostsPct,
			RampDownCapacityThresholdPercent: rampDownCapacityThresholdPct,
			RampDownForceLogoffUsers:         rampDownForceLogoffUsers,
			RampDownStopHostsWhen:            string(pointer.From(item.RampDownStopHostsWhen)),
			RampDownWaitTimeMinutes:          rampDownWaitTimeMinutes,
			RampDownNotificationMessage:      rampDownNotificationMessage,
			OffPeakStartTime:                 offPeakStartTime,
			OffPeakLoadBalancingAlgorithm:    string(pointer.From(item.OffPeakLoadBalancingAlgorithm)),
		})
	}
	return results
}

func flattenScalingHostpoolReference(input *[]scalingplan.ScalingHostPoolReference) []DesktopVirtualizationScalingPlanHostPool {
	results := make([]DesktopVirtualizationScalingPlanHostPool, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		hostPoolArmPath := ""
		if item.HostPoolArmPath != nil {
			hostPoolArmPath = *item.HostPoolArmPath
		}
		scalingPlanEnabled := false
		if item.ScalingPlanEnabled != nil {
			scalingPlanEnabled = *item.ScalingPlanEnabled
		}
		results = append(results, DesktopVirtualizationScalingPlanHostPool{
			HostpoolId:         hostPoolArmPath,
			ScalingPlanEnabled: scalingPlanEnabled,
		})
	}
	return results
}
