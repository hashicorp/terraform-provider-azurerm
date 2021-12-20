package desktopvirtualization

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2021-09-03-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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
			_, err := parse.ScalingPlanID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Monday",
									"Tuesday",
									"Wednesday",
									"Thursday",
									"Friday",
									"Saturday",
									"Sunday",
								}, false),
							},
						},

						"ramp_up_start_time": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validateTime(),
						},

						"ramp_up_load_balancing_algorithm": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmBreadthFirst),
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmDepthFirst),
							}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmBreadthFirst),
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmDepthFirst),
							}, false),
						},

						"ramp_down_start_time": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validateTime(),
						},

						"ramp_down_load_balancing_algorithm": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmBreadthFirst),
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmDepthFirst),
							}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.StopHostsWhenZeroActiveSessions),
								string(desktopvirtualization.StopHostsWhenZeroSessions),
							}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmBreadthFirst),
								string(desktopvirtualization.SessionHostLoadBalancingAlgorithmDepthFirst),
							}, false),
						},
					},
				},
			},

			"host_pool": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"hostpool_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.HostPoolID,
						},
						"scaling_plan_enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
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

	log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan create")

	id := parse.NewScalingPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s): %+v", id, err)
			}
		}

		if existing.ScalingPlanProperties != nil {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.ScalingPlan{
		Name:     utils.String(d.Get("name").(string)),
		Location: &location,
		Tags:     tags.Expand(t),
		ScalingPlanProperties: &desktopvirtualization.ScalingPlanProperties{
			Description:        utils.String(d.Get("description").(string)),
			FriendlyName:       utils.String(d.Get("friendly_name").(string)),
			TimeZone:           utils.String(d.Get("time_zone").(string)),
			HostPoolType:       desktopvirtualization.ScalingHostPoolTypePooled, // Only one possible value for this
			ExclusionTag:       utils.String(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]interface{})),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]interface{})),
		},
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.Name, context); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan update")

	id, err := parse.ScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	context := desktopvirtualization.ScalingPlanPatch{
		Tags: tags.Expand(t),
		ScalingPlanPatchProperties: &desktopvirtualization.ScalingPlanPatchProperties{
			Description:        utils.String(d.Get("description").(string)),
			FriendlyName:       utils.String(d.Get("friendly_name").(string)),
			TimeZone:           utils.String(d.Get("time_zone").(string)),
			ExclusionTag:       utils.String(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]interface{})),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]interface{})),
		},
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, &context); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}

func resourceVirtualDesktopScalingPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ScalingPlanProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("time_zone", props.TimeZone)
		d.Set("exclusion_tag", props.ExclusionTag)
		d.Set("schedule", flattenScalingPlanSchedule(props.Schedules))
		d.Set("host_pool", flattenScalingHostpoolReference(props.HostPoolReferences))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualDesktopScalingPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandScalingPlanSchedule(input []interface{}) *[]desktopvirtualization.ScalingSchedule {
	if len(input) == 0 {
		return nil
	}

	results := make([]desktopvirtualization.ScalingSchedule, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		name := v["name"].(string)
		daysOfWeek := v["days_of_week"].(*pluginsdk.Set).List()
		rampUpStartTime := v["ramp_up_start_time"].(string)
		rampUpLoadBalancingAlgorithm := v["ramp_up_load_balancing_algorithm"].(string)
		rampUpMinimumHostsPct := v["ramp_up_minimum_hosts_percent"].(int)
		rampUpCapacityThresholdPct := v["ramp_up_capacity_threshold_percent"].(int)
		peakStartTime := v["peak_start_time"].(string)
		peakLoadBalancingAlgorithm := v["peak_load_balancing_algorithm"].(string)
		rampDownStartTime := v["ramp_down_start_time"].(string)
		rampDownLoadBalancingAlgorithm := v["ramp_down_load_balancing_algorithm"].(string)
		rampDownMinimumHostsPct := v["ramp_down_minimum_hosts_percent"].(int)
		rampDownCapacityThresholdPct := v["ramp_down_capacity_threshold_percent"].(int)
		rampDownForceLogoffUsers := v["ramp_down_force_logoff_users"].(bool)
		rampDownStopHostsWhen := v["ramp_down_stop_hosts_when"].(string)
		rampDownWaitTimeMinutes := v["ramp_down_wait_time_minutes"].(int)
		rampDownNotificationMessage := v["ramp_down_notification_message"].(string)
		offPeakStartTime := v["off_peak_start_time"].(string)
		offPeakLoadBalancingAlgorithm := v["off_peak_load_balancing_algorithm"].(string)

		results = append(results, desktopvirtualization.ScalingSchedule{
			Name:                           utils.String(name),
			DaysOfWeek:                     utils.ExpandStringSlice(daysOfWeek),
			RampUpStartTime:                expandScalingPlanScheduleTime(rampUpStartTime),
			RampUpLoadBalancingAlgorithm:   desktopvirtualization.SessionHostLoadBalancingAlgorithm(rampUpLoadBalancingAlgorithm),
			RampUpMinimumHostsPct:          utils.Int32(int32(rampUpMinimumHostsPct)),
			RampUpCapacityThresholdPct:     utils.Int32(int32(rampUpCapacityThresholdPct)),
			PeakStartTime:                  expandScalingPlanScheduleTime(peakStartTime),
			PeakLoadBalancingAlgorithm:     desktopvirtualization.SessionHostLoadBalancingAlgorithm(peakLoadBalancingAlgorithm),
			RampDownStartTime:              expandScalingPlanScheduleTime(rampDownStartTime),
			RampDownLoadBalancingAlgorithm: desktopvirtualization.SessionHostLoadBalancingAlgorithm(rampDownLoadBalancingAlgorithm),
			RampDownMinimumHostsPct:        utils.Int32(int32(rampDownMinimumHostsPct)),
			RampDownCapacityThresholdPct:   utils.Int32(int32(rampDownCapacityThresholdPct)),
			RampDownForceLogoffUsers:       utils.Bool(rampDownForceLogoffUsers),
			RampDownStopHostsWhen:          desktopvirtualization.StopHostsWhen(rampDownStopHostsWhen),
			RampDownWaitTimeMinutes:        utils.Int32(int32(rampDownWaitTimeMinutes)),
			RampDownNotificationMessage:    utils.String(rampDownNotificationMessage),
			OffPeakStartTime:               expandScalingPlanScheduleTime(offPeakStartTime),
			OffPeakLoadBalancingAlgorithm:  desktopvirtualization.SessionHostLoadBalancingAlgorithm(offPeakLoadBalancingAlgorithm),
		})
	}

	return &results
}

func expandScalingPlanScheduleTime(input string) *desktopvirtualization.Time {
	if len(input) == 0 {
		return nil
	}

	time := strings.Split(input, ":")
	hour, _ := strconv.Atoi(time[0])
	minute, _ := strconv.Atoi(time[1])

	return &desktopvirtualization.Time{
		Hour:   utils.Int32(int32(hour)),
		Minute: utils.Int32(int32(minute)),
	}
}

func expandScalingPlanHostpoolReference(input []interface{}) *[]desktopvirtualization.ScalingHostPoolReference {
	if len(input) == 0 {
		return nil
	}

	results := make([]desktopvirtualization.ScalingHostPoolReference, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		hostPoolArmPath := v["hostpool_id"].(string)
		scalingPlanEnabled := v["scaling_plan_enabled"].(bool)

		results = append(results, desktopvirtualization.ScalingHostPoolReference{
			HostPoolArmPath:    utils.String(hostPoolArmPath),
			ScalingPlanEnabled: utils.Bool(scalingPlanEnabled),
		})
	}
	return &results
}

func flattenScalingPlanSchedule(input *[]desktopvirtualization.ScalingSchedule) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		name := ""
		if item.Name != nil {
			name = *item.Name
		}
		rampUpStartTime := ""
		if item.RampUpStartTime != nil && item.RampUpStartTime.Hour != nil && item.RampUpStartTime.Minute != nil {
			rampUpStartTime = fmt.Sprintf("%02d:%02d", *item.RampUpStartTime.Hour, *item.RampUpStartTime.Minute)
		}
		rampUpMinimumHostsPct := int32(0)
		if item.RampUpMinimumHostsPct != nil {
			rampUpMinimumHostsPct = *item.RampUpMinimumHostsPct
		}
		rampUpCapacityThresholdPct := int32(0)
		if item.RampUpCapacityThresholdPct != nil {
			rampUpCapacityThresholdPct = *item.RampUpCapacityThresholdPct
		}
		peakStartTime := ""
		if item.PeakStartTime != nil && item.PeakStartTime.Hour != nil && item.PeakStartTime.Minute != nil {
			peakStartTime = fmt.Sprintf("%02d:%02d", *item.PeakStartTime.Hour, *item.PeakStartTime.Minute)
		}
		rampDownStartTime := ""
		if item.RampDownStartTime != nil && item.RampDownStartTime.Hour != nil && item.RampDownStartTime.Minute != nil {
			rampDownStartTime = fmt.Sprintf("%02d:%02d", *item.RampDownStartTime.Hour, *item.RampDownStartTime.Minute)
		}
		rampDownMinimumHostsPct := int32(0)
		if item.RampDownMinimumHostsPct != nil {
			rampDownMinimumHostsPct = *item.RampDownMinimumHostsPct
		}
		rampDownCapacityThresholdPct := int32(0)
		if item.RampDownCapacityThresholdPct != nil {
			rampDownCapacityThresholdPct = *item.RampDownCapacityThresholdPct
		}
		rampDownForceLogoffUsers := false
		if item.RampDownForceLogoffUsers != nil {
			rampDownForceLogoffUsers = *item.RampDownForceLogoffUsers
		}
		rampDownWaitTimeMinutes := int32(0)
		if item.RampDownWaitTimeMinutes != nil {
			rampDownWaitTimeMinutes = *item.RampDownWaitTimeMinutes
		}
		rampDownNotificationMessage := ""
		if item.RampDownNotificationMessage != nil {
			rampDownNotificationMessage = *item.RampDownNotificationMessage
		}
		offPeakStartTime := ""
		if item.OffPeakStartTime != nil && item.OffPeakStartTime.Hour != nil && item.OffPeakStartTime.Minute != nil {
			offPeakStartTime = fmt.Sprintf("%02d:%02d", *item.OffPeakStartTime.Hour, *item.OffPeakStartTime.Minute)
		}
		results = append(results, map[string]interface{}{
			"name":                                 name,
			"days_of_week":                         utils.FlattenStringSlice(item.DaysOfWeek),
			"ramp_up_start_time":                   rampUpStartTime,
			"ramp_up_load_balancing_algorithm":     item.RampUpLoadBalancingAlgorithm,
			"ramp_up_minimum_hosts_percent":        rampUpMinimumHostsPct,
			"ramp_up_capacity_threshold_percent":   rampUpCapacityThresholdPct,
			"peak_start_time":                      peakStartTime,
			"peak_load_balancing_algorithm":        item.PeakLoadBalancingAlgorithm,
			"ramp_down_start_time":                 rampDownStartTime,
			"ramp_down_load_balancing_algorithm":   item.RampDownLoadBalancingAlgorithm,
			"ramp_down_minimum_hosts_percent":      rampDownMinimumHostsPct,
			"ramp_down_capacity_threshold_percent": rampDownCapacityThresholdPct,
			"ramp_down_force_logoff_users":         rampDownForceLogoffUsers,
			"ramp_down_stop_hosts_when":            item.RampDownStopHostsWhen,
			"ramp_down_wait_time_minutes":          rampDownWaitTimeMinutes,
			"ramp_down_notification_message":       rampDownNotificationMessage,
			"off_peak_start_time":                  offPeakStartTime,
			"off_peak_load_balancing_algorithm":    item.OffPeakLoadBalancingAlgorithm,
		})
	}
	return results
}

func flattenScalingHostpoolReference(input *[]desktopvirtualization.ScalingHostPoolReference) []interface{} {
	results := make([]interface{}, 0)
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
		results = append(results, map[string]interface{}{
			"hostpool_id":          hostPoolArmPath,
			"scaling_plan_enabled": scalingPlanEnabled,
		})
	}
	return results
}
