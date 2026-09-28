// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var scalingPlanResourceType = "azurerm_virtual_desktop_scaling_plan"

func resourceVirtualDesktopScalingPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopScalingPlanCreate,
		Read:   resourceVirtualDesktopScalingPlanRead,
		Update: resourceVirtualDesktopScalingPlanUpdate,
		Delete: resourceVirtualDesktopScalingPlanDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scalingplan.ParseScalingPlanID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
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
								ValidateFunc: validation.StringInSlice(scalingplan.PossibleValuesForDayOfWeek(), false),
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
				Computed: true, // azignore:AZS007 - pre-existing violation
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
		},
	}
}

func validateTime() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`), `The time must be in the format HH:MM.`)
}

func resourceVirtualDesktopScalingPlanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scalingplan.NewScalingPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s): %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	payload := scalingplan.ScalingPlan{
		Name:     pointer.To(d.Get("name").(string)),
		Location: location,
		Tags:     tags.Expand(t),
		Properties: scalingplan.ScalingPlanProperties{
			Description:        pointer.To(d.Get("description").(string)),
			FriendlyName:       pointer.To(d.Get("friendly_name").(string)),
			TimeZone:           d.Get("time_zone").(string),
			HostPoolType:       pointer.To(scalingplan.ScalingHostPoolTypePooled),
			ExclusionTag:       pointer.To(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]interface{})),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]interface{})),
		},
	}

	if _, err := client.Create(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scalingplan.ParseScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	payload := scalingplan.ScalingPlanPatch{
		Tags: tags.Expand(t),
		Properties: &scalingplan.ScalingPlanPatchProperties{
			Description:        pointer.To(d.Get("description").(string)),
			FriendlyName:       pointer.To(d.Get("friendly_name").(string)),
			TimeZone:           pointer.To(d.Get("time_zone").(string)),
			ExclusionTag:       pointer.To(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]interface{})),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]interface{})),
		},
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scalingplan.ParseScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ScalingPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("description", model.Properties.Description)
		d.Set("friendly_name", model.Properties.FriendlyName)
		d.Set("time_zone", model.Properties.TimeZone)
		d.Set("exclusion_tag", model.Properties.ExclusionTag)
		d.Set("schedule", flattenScalingPlanSchedule(model.Properties.Schedules))
		d.Set("host_pool", flattenScalingHostpoolReference(model.Properties.HostPoolReferences))

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceVirtualDesktopScalingPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scalingplan.ParseScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandScalingPlanSchedule(input []interface{}) *[]scalingplan.ScalingSchedule {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingSchedule, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		name := v["name"].(string)
		daysOfWeekRaw := v["days_of_week"].(*pluginsdk.Set).List()
		daysOfWeek := make([]scalingplan.DayOfWeek, 0)
		for _, weekday := range daysOfWeekRaw {
			daysOfWeek = append(daysOfWeek, scalingplan.DayOfWeek(weekday.(string)))
		}

		rampUpStartTime := v["ramp_up_start_time"].(string)
		rampUpMinimumHostsPct := v["ramp_up_minimum_hosts_percent"].(int)
		rampUpCapacityThresholdPct := v["ramp_up_capacity_threshold_percent"].(int)
		peakStartTime := v["peak_start_time"].(string)
		rampDownStartTime := v["ramp_down_start_time"].(string)
		rampDownMinimumHostsPct := v["ramp_down_minimum_hosts_percent"].(int)
		rampDownCapacityThresholdPct := v["ramp_down_capacity_threshold_percent"].(int)
		rampDownForceLogoffUsers := v["ramp_down_force_logoff_users"].(bool)
		rampDownWaitTimeMinutes := v["ramp_down_wait_time_minutes"].(int)
		rampDownNotificationMessage := v["ramp_down_notification_message"].(string)
		offPeakStartTime := v["off_peak_start_time"].(string)

		results = append(results, scalingplan.ScalingSchedule{
			Name:                           pointer.To(name),
			DaysOfWeek:                     &daysOfWeek,
			RampUpStartTime:                expandScalingPlanScheduleTime(rampUpStartTime),
			RampUpLoadBalancingAlgorithm:   pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["ramp_up_load_balancing_algorithm"].(string)),
			RampUpMinimumHostsPct:          pointer.To(int64(rampUpMinimumHostsPct)),
			RampUpCapacityThresholdPct:     pointer.To(int64(rampUpCapacityThresholdPct)),
			PeakStartTime:                  expandScalingPlanScheduleTime(peakStartTime),
			PeakLoadBalancingAlgorithm:     pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["peak_load_balancing_algorithm"].(string)),
			RampDownStartTime:              expandScalingPlanScheduleTime(rampDownStartTime),
			RampDownLoadBalancingAlgorithm: pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["ramp_down_load_balancing_algorithm"].(string)),
			RampDownMinimumHostsPct:        pointer.To(int64(rampDownMinimumHostsPct)),
			RampDownCapacityThresholdPct:   pointer.To(int64(rampDownCapacityThresholdPct)),
			RampDownForceLogoffUsers:       pointer.To(rampDownForceLogoffUsers),
			RampDownStopHostsWhen:          pointer.ToEnum[scalingplan.StopHostsWhen](v["ramp_down_stop_hosts_when"].(string)),
			RampDownWaitTimeMinutes:        pointer.To(int64(rampDownWaitTimeMinutes)),
			RampDownNotificationMessage:    pointer.To(rampDownNotificationMessage),
			OffPeakStartTime:               expandScalingPlanScheduleTime(offPeakStartTime),
			OffPeakLoadBalancingAlgorithm:  pointer.ToEnum[scalingplan.SessionHostLoadBalancingAlgorithm](v["off_peak_load_balancing_algorithm"].(string)),
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

func expandScalingPlanHostpoolReference(input []interface{}) *[]scalingplan.ScalingHostPoolReference {
	if len(input) == 0 {
		return nil
	}

	results := make([]scalingplan.ScalingHostPoolReference, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		hostPoolArmPath := v["hostpool_id"].(string)
		scalingPlanEnabled := v["scaling_plan_enabled"].(bool)

		results = append(results, scalingplan.ScalingHostPoolReference{
			HostPoolArmPath:    pointer.To(hostPoolArmPath),
			ScalingPlanEnabled: pointer.To(scalingPlanEnabled),
		})
	}
	return &results
}

func flattenScalingPlanSchedule(input *[]scalingplan.ScalingSchedule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		rampUpStartTime := ""
		if item.RampUpStartTime != nil {
			rampUpStartTime = fmt.Sprintf("%02d:%02d", item.RampUpStartTime.Hour, item.RampUpStartTime.Minute)
		}
		peakStartTime := ""
		if item.PeakStartTime != nil {
			peakStartTime = fmt.Sprintf("%02d:%02d", item.PeakStartTime.Hour, item.PeakStartTime.Minute)
		}
		rampDownStartTime := ""
		if item.RampDownStartTime != nil {
			rampDownStartTime = fmt.Sprintf("%02d:%02d", item.RampDownStartTime.Hour, item.RampDownStartTime.Minute)
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

		results = append(results, map[string]interface{}{
			"name":                                 pointer.From(item.Name),
			"days_of_week":                         daysOfWeek,
			"ramp_up_start_time":                   rampUpStartTime,
			"ramp_up_load_balancing_algorithm":     item.RampUpLoadBalancingAlgorithm,
			"ramp_up_minimum_hosts_percent":        pointer.From(item.RampUpMinimumHostsPct),
			"ramp_up_capacity_threshold_percent":   pointer.From(item.RampUpCapacityThresholdPct),
			"peak_start_time":                      peakStartTime,
			"peak_load_balancing_algorithm":        item.PeakLoadBalancingAlgorithm,
			"ramp_down_start_time":                 rampDownStartTime,
			"ramp_down_load_balancing_algorithm":   item.RampDownLoadBalancingAlgorithm,
			"ramp_down_minimum_hosts_percent":      pointer.From(item.RampDownMinimumHostsPct),
			"ramp_down_capacity_threshold_percent": pointer.From(item.RampDownCapacityThresholdPct),
			"ramp_down_force_logoff_users":         pointer.From(item.RampDownForceLogoffUsers),
			"ramp_down_stop_hosts_when":            item.RampDownStopHostsWhen,
			"ramp_down_wait_time_minutes":          pointer.From(item.RampDownWaitTimeMinutes),
			"ramp_down_notification_message":       pointer.From(item.RampDownNotificationMessage),
			"off_peak_start_time":                  offPeakStartTime,
			"off_peak_load_balancing_algorithm":    item.OffPeakLoadBalancingAlgorithm,
		})
	}
	return results
}

func flattenScalingHostpoolReference(input *[]scalingplan.ScalingHostPoolReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"hostpool_id":          pointer.From(item.HostPoolArmPath),
			"scaling_plan_enabled": pointer.From(item.ScalingPlanEnabled),
		})
	}
	return results
}
