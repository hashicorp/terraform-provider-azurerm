package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2021-10-01/netapp"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/parse"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppSnapshotPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppSnapshotPolicyCreate,
		Read:   resourceNetAppSnapshotPolicyRead,
		Update: resourceNetAppSnapshotPolicyUpdate,
		Delete: resourceNetAppSnapshotPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SnapshotPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.SnapshotName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.AccountName,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"hourly_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},

						"minute": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"daily_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},

						"hour": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},

						"minute": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"weekly_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},

						"days_of_week": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 7,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsDayOfTheWeek(false),
							},
						},

						"hour": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},

						"minute": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"monthly_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"snapshots_to_keep": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 255),
						},

						"days_of_month": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 30,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeInt,
								ValidateFunc: validation.IntBetween(1, 30),
							},
						},

						"hour": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},

						"minute": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceNetAppSnapshotPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, accountName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing NetApp SnapshotPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_netapp_snapshot_policy", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	enabled := d.Get("enabled").(bool)

	hourlyScheduleRaw := d.Get("hourly_schedule").([]interface{})
	hourlySchedule := expandNetAppSnapshotPolicyHourlySchedule(hourlyScheduleRaw)

	dailyScheduleRaw := d.Get("daily_schedule").([]interface{})
	dailySchedule := expandNetAppSnapshotPolicyDailySchedule(dailyScheduleRaw)

	weeklyScheduleRaw := d.Get("weekly_schedule").([]interface{})
	weeklySchedule := expandNetAppSnapshotPolicyWeeklySchedule(weeklyScheduleRaw)

	monthlyScheduleRaw := d.Get("monthly_schedule").([]interface{})
	monthlySchedule := expandNetAppSnapshotPolicyMonthlySchedule(monthlyScheduleRaw)

	parameters := netapp.SnapshotPolicy{
		Location: utils.String(location),
		Name:     utils.String(name),
		SnapshotPolicyProperties: &netapp.SnapshotPolicyProperties{
			HourlySchedule:  hourlySchedule,
			DailySchedule:   dailySchedule,
			WeeklySchedule:  weeklySchedule,
			MonthlySchedule: monthlySchedule,
			Enabled:         utils.Bool(enabled),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Create(ctx, parameters, resourceGroup, accountName, name); err != nil {
		return fmt.Errorf("creating NetApp SnapshotPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Waiting for snapshot policy be completely provisioned
	id := parse.NewSnapshotPolicyID(client.SubscriptionID, resourceGroup, accountName, name)
	log.Printf("[DEBUG] Waiting for NetApp Snapshot Policy Provisioning Service %q (Resource Group %q) to complete", id.Name, id.ResourceGroup)
	if err := waitForSnapshotPolicyCreation(ctx, client, id, d.Timeout(pluginsdk.TimeoutDelete)); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceNetAppSnapshotPolicyRead(d, meta)
}

func resourceNetAppSnapshotPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)

	location := azure.NormalizeLocation(d.Get("location").(string))

	enabled := d.Get("enabled").(bool)

	hourlyScheduleRaw := d.Get("hourly_schedule").([]interface{})
	hourlySchedule := expandNetAppSnapshotPolicyHourlySchedule(hourlyScheduleRaw)

	dailyScheduleRaw := d.Get("daily_schedule").([]interface{})
	dailySchedule := expandNetAppSnapshotPolicyDailySchedule(dailyScheduleRaw)

	weeklyScheduleRaw := d.Get("weekly_schedule").([]interface{})
	weeklySchedule := expandNetAppSnapshotPolicyWeeklySchedule(weeklyScheduleRaw)

	monthlyScheduleRaw := d.Get("monthly_schedule").([]interface{})
	monthlySchedule := expandNetAppSnapshotPolicyMonthlySchedule(monthlyScheduleRaw)

	parameters := netapp.SnapshotPolicyPatch{
		Location: utils.String(location),
		Name:     utils.String(name),
		SnapshotPolicyProperties: &netapp.SnapshotPolicyProperties{
			HourlySchedule:  hourlySchedule,
			DailySchedule:   dailySchedule,
			WeeklySchedule:  weeklySchedule,
			MonthlySchedule: monthlySchedule,
			Enabled:         utils.Bool(enabled),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, parameters, resourceGroup, accountName, name)
	if err != nil {
		return fmt.Errorf("updating NetApp SnapshotPolicy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceNetAppSnapshotPolicyRead(d, meta)
}

func resourceNetAppSnapshotPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NetApp SnapshotPolicy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading NetApp SnapshotPolicy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.NetAppAccountName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.SnapshotPolicyProperties; props != nil {
		d.Set("enabled", props.Enabled)
		if err := d.Set("hourly_schedule", flattenNetAppVolumeSnapshotPolicyHourlySchedule(props.HourlySchedule)); err != nil {
			return fmt.Errorf("setting `hourly_schedule`: %+v", err)
		}
		if err := d.Set("daily_schedule", flattenNetAppVolumeSnapshotPolicyDailySchedule(props.DailySchedule)); err != nil {
			return fmt.Errorf("setting `daily_schedule`: %+v", err)
		}
		if err := d.Set("weekly_schedule", flattenNetAppVolumeSnapshotPolicyWeeklySchedule(props.WeeklySchedule)); err != nil {
			return fmt.Errorf("setting `weekly_schedule`: %+v", err)
		}
		if err := d.Set("monthly_schedule", flattenNetAppVolumeSnapshotPolicyMonthlySchedule(props.MonthlySchedule)); err != nil {
			return fmt.Errorf("setting `monthly_schedule`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNetAppSnapshotPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	// Deleting snapshot policy and waiting for it fo fully complete the operation
	future, err := client.Delete(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting NetApp Snapshot Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for NetApp SnapshotPolicy Provisioning Service %q (Resource Group %q) to be deleted", id.Name, id.ResourceGroup)
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}
	if err := waitForSnapshotPolicyDeletion(ctx, client, *id, d.Timeout(pluginsdk.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func expandNetAppSnapshotPolicyHourlySchedule(input []interface{}) *netapp.HourlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &netapp.HourlySchedule{}
	}

	hourlyScheduleObject := netapp.HourlySchedule{}

	hourlyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := hourlyScheduleRaw["snapshots_to_keep"]; ok {
		hourlyScheduleObject.SnapshotsToKeep = utils.Int32(int32(v.(int)))
	}
	if v, ok := hourlyScheduleRaw["minute"]; ok {
		hourlyScheduleObject.Minute = utils.Int32(int32(v.(int)))
	}

	return &hourlyScheduleObject
}

func expandNetAppSnapshotPolicyDailySchedule(input []interface{}) *netapp.DailySchedule {
	if len(input) == 0 || input[0] == nil {
		return &netapp.DailySchedule{}
	}

	dailyScheduleObject := netapp.DailySchedule{}

	dailyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := dailyScheduleRaw["snapshots_to_keep"]; ok {
		dailyScheduleObject.SnapshotsToKeep = utils.Int32(int32(v.(int)))
	}
	if v, ok := dailyScheduleRaw["hour"]; ok {
		dailyScheduleObject.Hour = utils.Int32(int32(v.(int)))
	}
	if v, ok := dailyScheduleRaw["minute"]; ok {
		dailyScheduleObject.Minute = utils.Int32(int32(v.(int)))
	}

	return &dailyScheduleObject
}

func expandNetAppSnapshotPolicyWeeklySchedule(input []interface{}) *netapp.WeeklySchedule {
	if len(input) == 0 || input[0] == nil {
		return &netapp.WeeklySchedule{}
	}

	weeklyScheduleObject := netapp.WeeklySchedule{}

	weeklyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := weeklyScheduleRaw["snapshots_to_keep"]; ok {
		weeklyScheduleObject.SnapshotsToKeep = utils.Int32(int32(v.(int)))
	}
	if _, ok := weeklyScheduleRaw["days_of_week"]; ok {
		weeklyScheduleObject.Day = utils.ExpandStringSliceWithDelimiter(weeklyScheduleRaw["days_of_week"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := weeklyScheduleRaw["hour"]; ok {
		weeklyScheduleObject.Hour = utils.Int32(int32(v.(int)))
	}
	if v, ok := weeklyScheduleRaw["minute"]; ok {
		weeklyScheduleObject.Minute = utils.Int32(int32(v.(int)))
	}

	return &weeklyScheduleObject
}

func expandNetAppSnapshotPolicyMonthlySchedule(input []interface{}) *netapp.MonthlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &netapp.MonthlySchedule{}
	}

	monthlyScheduleObject := netapp.MonthlySchedule{}

	monthlyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := monthlyScheduleRaw["snapshots_to_keep"]; ok {
		monthlyScheduleObject.SnapshotsToKeep = utils.Int32(int32(v.(int)))
	}
	if _, ok := monthlyScheduleRaw["days_of_month"]; ok {
		monthlyScheduleObject.DaysOfMonth = utils.ExpandIntSliceWithDelimiter(monthlyScheduleRaw["days_of_month"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := monthlyScheduleRaw["hour"]; ok {
		monthlyScheduleObject.Hour = utils.Int32(int32(v.(int)))
	}
	if v, ok := monthlyScheduleRaw["minute"]; ok {
		monthlyScheduleObject.Minute = utils.Int32(int32(v.(int)))
	}

	return &monthlyScheduleObject
}

func flattenNetAppVolumeSnapshotPolicyHourlySchedule(input *netapp.HourlySchedule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyDailySchedule(input *netapp.DailySchedule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyWeeklySchedule(input *netapp.WeeklySchedule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	weekDays := make([]interface{}, 0)
	if input.Day != nil {
		for _, day := range strings.Split(*input.Day, ",") {
			weekDays = append(weekDays, day)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"days_of_week":      weekDays,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyMonthlySchedule(input *netapp.MonthlySchedule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	daysOfMonth := make([]interface{}, 0)
	if input.DaysOfMonth != nil {
		for _, day := range strings.Split(*input.DaysOfMonth, ",") {
			intDay, _ := strconv.Atoi(day)
			daysOfMonth = append(daysOfMonth, intDay)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"days_of_month":     daysOfMonth,
			"hour":              input.Hour,
			"minute":            input.Minute,
		},
	}
}

func waitForSnapshotPolicyCreation(ctx context.Context, client *netapp.SnapshotPoliciesClient, id parse.SnapshotPolicyId, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappSnapshotPolicyStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting NetApp Volume Provisioning Service %q (Resource Group %q) to complete: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForSnapshotPolicyDeletion(ctx context.Context, client *netapp.SnapshotPoliciesClient, id parse.SnapshotPolicyId, timeout time.Duration) error {
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappSnapshotPolicyStateRefreshFunc(ctx, client, id),
		Timeout:                   timeout,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for NetApp SnapshotPolicy Provisioning Service %q (Resource Group %q) to be deleted: %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func netappSnapshotPolicyStateRefreshFunc(ctx context.Context, client *netapp.SnapshotPoliciesClient, id parse.SnapshotPolicyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.NetAppAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("retrieving NetApp SnapshotPolicy %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
			}
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
