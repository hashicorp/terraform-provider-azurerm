// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/snapshotpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
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
			_, err := snapshotpolicy.ParseSnapshotPolicyID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.SnapshotName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

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

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("hourly_schedule", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),

			pluginsdk.ForceNewIfChange("daily_schedule", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),

			pluginsdk.ForceNewIfChange("weekly_schedule", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),

			pluginsdk.ForceNewIfChange("monthly_schedule", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) > 0 && len(new.([]interface{})) == 0
			}),
		),
	}
}

func resourceNetAppSnapshotPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshotpolicy.NewSnapshotPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.SnapshotPoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_snapshot_policy", id.ID())
		}
	}

	parameters := snapshotpolicy.SnapshotPolicy{
		Location: azure.NormalizeLocation(d.Get("location").(string)),
		Name:     utils.String(id.SnapshotPolicyName),
		Properties: snapshotpolicy.SnapshotPolicyProperties{
			HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(d.Get("hourly_schedule").([]interface{})),
			DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(d.Get("daily_schedule").([]interface{})),
			WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(d.Get("weekly_schedule").([]interface{})),
			MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(d.Get("monthly_schedule").([]interface{})),
			Enabled:         utils.Bool(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.SnapshotPoliciesCreate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Waiting for snapshot policy be completely provisioned
	log.Printf("[DEBUG] Waiting for %s to complete", id)
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

	id, err := snapshotpolicy.ParseSnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	parameters := snapshotpolicy.SnapshotPolicyPatch{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Name:     utils.String(id.SnapshotPolicyName),
		Properties: &snapshotpolicy.SnapshotPolicyProperties{
			HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(d.Get("hourly_schedule").([]interface{})),
			DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(d.Get("daily_schedule").([]interface{})),
			WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(d.Get("weekly_schedule").([]interface{})),
			MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(d.Get("monthly_schedule").([]interface{})),
			Enabled:         utils.Bool(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err = client.SnapshotPoliciesUpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNetAppSnapshotPolicyRead(d, meta)
}

func resourceNetAppSnapshotPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshotpolicy.ParseSnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SnapshotPoliciesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] NetApp SnapshotPolicy %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.SnapshotPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.NetAppAccountName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		props := model.Properties
		d.Set("enabled", props.Enabled)
		if err = d.Set("hourly_schedule", flattenNetAppVolumeSnapshotPolicyHourlySchedule(props.HourlySchedule)); err != nil {
			return fmt.Errorf("setting `hourly_schedule`: %+v", err)
		}
		if err = d.Set("daily_schedule", flattenNetAppVolumeSnapshotPolicyDailySchedule(props.DailySchedule)); err != nil {
			return fmt.Errorf("setting `daily_schedule`: %+v", err)
		}
		if err = d.Set("weekly_schedule", flattenNetAppVolumeSnapshotPolicyWeeklySchedule(props.WeeklySchedule)); err != nil {
			return fmt.Errorf("setting `weekly_schedule`: %+v", err)
		}
		if err = d.Set("monthly_schedule", flattenNetAppVolumeSnapshotPolicyMonthlySchedule(props.MonthlySchedule)); err != nil {
			return fmt.Errorf("setting `monthly_schedule`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetAppSnapshotPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	volumeClient := meta.(*clients.Client).NetApp.VolumeClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshotpolicy.ParseSnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	// Try to delete the snapshot policy using DeleteThenPoll
	err = client.SnapshotPoliciesDeleteThenPoll(ctx, *id)
	if err != nil {
		// Check if error is about snapshot policy being in use
		if strings.Contains(err.Error(), "SnapshotPolicy is used") {
			// Get all volumes in the account that might be using this snapshot policy
			volumeIds, err := findVolumesUsingSnapshotPolicy(ctx, meta.(*clients.Client), *id)
			if err != nil {
				return fmt.Errorf("finding volumes using snapshot policy %s: %+v", *id, err)
			}

			// Disassociate snapshot policy from each volume
			for _, volumeId := range volumeIds {
				volId, err := volumes.ParseVolumeID(volumeId)
				if err != nil {
					return fmt.Errorf("parsing volume ID %q: %+v", volumeId, err)
				}

				// Update volume to remove snapshot policy
				update := volumes.VolumePatch{
					Properties: &volumes.VolumePatchProperties{
						DataProtection: &volumes.VolumePatchPropertiesDataProtection{
							Snapshot: &volumes.VolumeSnapshotProperties{
								SnapshotPolicyId: pointer.To(""),
							},
						},
					},
				}

				locks.ByID(volumeId)

				if err = volumeClient.UpdateThenPoll(ctx, *volId, update); err != nil {
					locks.UnlockByID(volumeId)
					return fmt.Errorf("removing snapshot policy from volume %s: %+v", *volId, err)
				}

				// Wait for the update to complete
				if err := waitForVolumeCreateOrUpdate(ctx, volumeClient, *volId); err != nil {
					locks.UnlockByID(volumeId)
					return fmt.Errorf("waiting for snapshot policy removal from volume %s: %+v", *volId, err)
				}

				locks.UnlockByID(volumeId)
			}

			// Try deleting the snapshot policy again
			if err = client.SnapshotPoliciesDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s after volume disassociation: %+v", *id, err)
			}
		} else {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	if err := waitForSnapshotPolicyDeletion(ctx, client, *id, d.Timeout(pluginsdk.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func findVolumesUsingSnapshotPolicy(ctx context.Context, client *clients.Client, snapshotPolicyId snapshotpolicy.SnapshotPolicyId) ([]string, error) {
	volumeIds := make([]string, 0)
	poolClient := client.NetApp.PoolClient
	accountId := capacitypools.NewNetAppAccountID(snapshotPolicyId.SubscriptionId, snapshotPolicyId.ResourceGroupName, snapshotPolicyId.NetAppAccountName)

	poolsResult, err := poolClient.PoolsList(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("listing capacity pools in account %s: %+v", snapshotPolicyId.NetAppAccountName, err)
	}

	if model := poolsResult.Model; model != nil {
		volumeClient := client.NetApp.VolumeClient

		for _, pool := range *model {
			if pool.Name == nil {
				continue
			}

			poolNameParts := strings.Split(pointer.From(pool.Name), "/")
			poolName := poolNameParts[len(poolNameParts)-1]
			volumeId := volumes.NewCapacityPoolID(snapshotPolicyId.SubscriptionId, snapshotPolicyId.ResourceGroupName, snapshotPolicyId.NetAppAccountName, poolName)

			volumesResult, err := volumeClient.List(ctx, volumeId)
			if err != nil {
				return nil, fmt.Errorf("listing volumes in pool %s: %+v", poolName, err)
			}

			if volumesModel := volumesResult.Model; volumesModel != nil {
				for _, volume := range *volumesModel {
					if volume.Id == nil || volume.Properties.DataProtection == nil ||
						volume.Properties.DataProtection.Snapshot == nil || volume.Properties.DataProtection.Snapshot.SnapshotPolicyId == nil {
						continue
					}

					if strings.EqualFold(*volume.Properties.DataProtection.Snapshot.SnapshotPolicyId, snapshotPolicyId.ID()) {
						volumeIds = append(volumeIds, *volume.Id)
					}
				}
			}
		}
	}

	return volumeIds, nil
}

func expandNetAppSnapshotPolicyHourlySchedule(input []interface{}) *snapshotpolicy.HourlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicy.HourlySchedule{}
	}

	hourlyScheduleObject := snapshotpolicy.HourlySchedule{}

	hourlyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := hourlyScheduleRaw["snapshots_to_keep"]; ok {
		hourlyScheduleObject.SnapshotsToKeep = utils.Int64(int64(v.(int)))
	}
	if v, ok := hourlyScheduleRaw["minute"]; ok {
		hourlyScheduleObject.Minute = utils.Int64(int64(v.(int)))
	}

	return &hourlyScheduleObject
}

func expandNetAppSnapshotPolicyDailySchedule(input []interface{}) *snapshotpolicy.DailySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicy.DailySchedule{}
	}

	dailyScheduleObject := snapshotpolicy.DailySchedule{}

	dailyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := dailyScheduleRaw["snapshots_to_keep"]; ok {
		dailyScheduleObject.SnapshotsToKeep = utils.Int64(int64(v.(int)))
	}
	if v, ok := dailyScheduleRaw["hour"]; ok {
		dailyScheduleObject.Hour = utils.Int64(int64(v.(int)))
	}
	if v, ok := dailyScheduleRaw["minute"]; ok {
		dailyScheduleObject.Minute = utils.Int64(int64(v.(int)))
	}

	return &dailyScheduleObject
}

func expandNetAppSnapshotPolicyWeeklySchedule(input []interface{}) *snapshotpolicy.WeeklySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicy.WeeklySchedule{}
	}

	weeklyScheduleObject := snapshotpolicy.WeeklySchedule{}

	weeklyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := weeklyScheduleRaw["snapshots_to_keep"]; ok {
		weeklyScheduleObject.SnapshotsToKeep = utils.Int64(int64(v.(int)))
	}
	if _, ok := weeklyScheduleRaw["days_of_week"]; ok {
		weeklyScheduleObject.Day = utils.ExpandStringSliceWithDelimiter(weeklyScheduleRaw["days_of_week"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := weeklyScheduleRaw["hour"]; ok {
		weeklyScheduleObject.Hour = utils.Int64(int64(v.(int)))
	}
	if v, ok := weeklyScheduleRaw["minute"]; ok {
		weeklyScheduleObject.Minute = utils.Int64(int64(v.(int)))
	}

	return &weeklyScheduleObject
}

func expandNetAppSnapshotPolicyMonthlySchedule(input []interface{}) *snapshotpolicy.MonthlySchedule {
	if len(input) == 0 || input[0] == nil {
		return &snapshotpolicy.MonthlySchedule{}
	}

	monthlyScheduleObject := snapshotpolicy.MonthlySchedule{}

	monthlyScheduleRaw := input[0].(map[string]interface{})

	if v, ok := monthlyScheduleRaw["snapshots_to_keep"]; ok {
		monthlyScheduleObject.SnapshotsToKeep = utils.Int64(int64(v.(int)))
	}
	if _, ok := monthlyScheduleRaw["days_of_month"]; ok {
		monthlyScheduleObject.DaysOfMonth = utils.ExpandIntSliceWithDelimiter(monthlyScheduleRaw["days_of_month"].(*pluginsdk.Set).List(), ",")
	}
	if v, ok := monthlyScheduleRaw["hour"]; ok {
		monthlyScheduleObject.Hour = utils.Int64(int64(v.(int)))
	}
	if v, ok := monthlyScheduleRaw["minute"]; ok {
		monthlyScheduleObject.Minute = utils.Int64(int64(v.(int)))
	}

	return &monthlyScheduleObject
}

func flattenNetAppVolumeSnapshotPolicyHourlySchedule(input *snapshotpolicy.HourlySchedule) []interface{} {
	if input == nil || (input.Minute == nil && input.SnapshotsToKeep == nil) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"snapshots_to_keep": input.SnapshotsToKeep,
			"minute":            input.Minute,
		},
	}
}

func flattenNetAppVolumeSnapshotPolicyDailySchedule(input *snapshotpolicy.DailySchedule) []interface{} {
	if input == nil || (input.SnapshotsToKeep == nil && input.Hour == nil && input.Minute == nil) {
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

func flattenNetAppVolumeSnapshotPolicyWeeklySchedule(input *snapshotpolicy.WeeklySchedule) []interface{} {
	if input == nil || (input.SnapshotsToKeep == nil && input.Day == nil && input.Hour == nil && input.Minute == nil) {
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

func flattenNetAppVolumeSnapshotPolicyMonthlySchedule(input *snapshotpolicy.MonthlySchedule) []interface{} {
	if input == nil || (input.SnapshotsToKeep == nil && input.DaysOfMonth == nil && input.Hour == nil && input.Minute == nil) {
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

func waitForSnapshotPolicyCreation(ctx context.Context, client *snapshotpolicy.SnapshotPolicyClient, id snapshotpolicy.SnapshotPolicyId, timeout time.Duration) error {
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
		return fmt.Errorf("waiting for %s to complete: %+v", id, err)
	}

	return nil
}

func waitForSnapshotPolicyDeletion(ctx context.Context, client *snapshotpolicy.SnapshotPolicyClient, id snapshotpolicy.SnapshotPolicyId, timeout time.Duration) error {
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
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappSnapshotPolicyStateRefreshFunc(ctx context.Context, client *snapshotpolicy.SnapshotPolicyClient, id snapshotpolicy.SnapshotPolicyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.SnapshotPoliciesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}
