package recoveryservices

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2019-05-13/backup"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmBackupProtectionPolicyFileShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmBackupProtectionPolicyFileShareCreateUpdate,
		Read:   resourceArmBackupProtectionPolicyFileShareRead,
		Update: resourceArmBackupProtectionPolicyFileShareCreateUpdate,
		Delete: resourceArmBackupProtectionPolicyFileShareDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-_!a-zA-Z0-9]{2,149}$"),
					"Backup Policy name must be 3 - 150 characters long, start with a letter, contain only letters and numbers.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},

			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},

			"backup": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"frequency": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(backup.ScheduleRunTypeDaily),
							}, false),
						},

						"time": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^([01][0-9]|[2][0-3]):([03][0])$"), // time must be on the hour or half past
								"Time of day must match the format HH:mm where HH is 00-23 and mm is 00 or 30",
							),
						},
					},
				},
			},

			"retention_daily": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 180),
						},
					},
				},
			},
		},
	}
}

func resourceArmBackupProtectionPolicyFileShareCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	policyName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	log.Printf("[DEBUG] Creating/updating Recovery Service Protection Policy %s (resource group %q)", policyName, resourceGroup)

	// getting this ready now because its shared between *everything*, time is... complicated for this resource
	timeOfDay := d.Get("backup.0.time").(string)
	dateOfDay, err := time.Parse(time.RFC3339, fmt.Sprintf("2018-07-30T%s:00Z", timeOfDay))
	if err != nil {
		return fmt.Errorf("Error generating time from %q for policy %q (Resource Group %q): %+v", timeOfDay, policyName, resourceGroup, err)
	}
	times := append(make([]date.Time, 0), date.Time{Time: dateOfDay})

	if d.IsNewResource() {
		existing, err2 := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err2 != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err2)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_backup_policy_file_share", *existing.ID)
		}
	}

	policy := backup.ProtectionPolicyResource{
		Properties: &backup.AzureFileShareProtectionPolicy{
			TimeZone:             utils.String(d.Get("timezone").(string)),
			BackupManagementType: backup.BackupManagementTypeAzureStorage,
			WorkLoadType:         backup.WorkloadTypeAzureFileShare,
			SchedulePolicy:       expandArmBackupProtectionPolicyFileShareSchedule(d, times),
			RetentionPolicy: &backup.LongTermRetentionPolicy{ // SimpleRetentionPolicy only has duration property ¯\_(ツ)_/¯
				RetentionPolicyType: backup.RetentionPolicyTypeLongTermRetentionPolicy,
				DailySchedule:       expandArmBackupProtectionPolicyFileShareRetentionDaily(d, times),
			},
		},
	}
	if _, err = client.CreateOrUpdate(ctx, vaultName, resourceGroup, policyName, policy); err != nil {
		return fmt.Errorf("Error creating/updating Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	resp, err := resourceArmBackupProtectionPolicyFileShareWaitForUpdate(ctx, client, vaultName, resourceGroup, policyName, d)
	if err != nil {
		return err
	}

	id := strings.Replace(*resp.ID, "Subscriptions", "subscriptions", 1)
	d.SetId(id)

	return resourceArmBackupProtectionPolicyFileShareRead(d, meta)
}

func resourceArmBackupProtectionPolicyFileShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	policyName := id.Path["backupPolicies"]
	vaultName := id.Path["vaults"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Reading Recovery Service Protection Policy %q (resource group %q)", policyName, resourceGroup)

	resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	d.Set("name", policyName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("recovery_vault_name", vaultName)

	if properties, ok := resp.Properties.AsAzureFileShareProtectionPolicy(); ok && properties != nil {
		d.Set("timezone", properties.TimeZone)

		if schedule, ok := properties.SchedulePolicy.AsSimpleSchedulePolicy(); ok && schedule != nil {
			if err := d.Set("backup", flattenArmBackupProtectionPolicyFileShareSchedule(schedule)); err != nil {
				return fmt.Errorf("Error setting `backup`: %+v", err)
			}
		}

		if retention, ok := properties.RetentionPolicy.AsLongTermRetentionPolicy(); ok && retention != nil {
			if s := retention.DailySchedule; s != nil {
				if err := d.Set("retention_daily", flattenArmBackupProtectionPolicyFileShareRetentionDaily(s)); err != nil {
					return fmt.Errorf("Error setting `retention_daily`: %+v", err)
				}
			} else {
				d.Set("retention_daily", nil)
			}
		}
	}

	return nil
}

func resourceArmBackupProtectionPolicyFileShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ProtectionPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	policyName := id.Path["backupPolicies"]
	resourceGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]

	log.Printf("[DEBUG] Deleting Recovery Service Protection Policy %q (resource group %q)", policyName, resourceGroup)

	resp, err := client.Delete(ctx, vaultName, resourceGroup, policyName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
		}
	}

	if _, err := resourceArmBackupProtectionPolicyFileShareWaitForDeletion(ctx, client, vaultName, resourceGroup, policyName, d); err != nil {
		return err
	}

	return nil
}

func expandArmBackupProtectionPolicyFileShareSchedule(d *schema.ResourceData, times []date.Time) *backup.SimpleSchedulePolicy {
	if bb, ok := d.Get("backup").([]interface{}); ok && len(bb) > 0 {
		block := bb[0].(map[string]interface{})

		schedule := backup.SimpleSchedulePolicy{ // LongTermSchedulePolicy has no properties
			SchedulePolicyType: backup.SchedulePolicyTypeSimpleSchedulePolicy,
			ScheduleRunTimes:   &times,
		}

		if v, ok := block["frequency"].(string); ok {
			schedule.ScheduleRunFrequency = backup.ScheduleRunType(v)
		}

		return &schedule
	}

	return nil
}

func expandArmBackupProtectionPolicyFileShareRetentionDaily(d *schema.ResourceData, times []date.Time) *backup.DailyRetentionSchedule {
	if rb, ok := d.Get("retention_daily").([]interface{}); ok && len(rb) > 0 {
		block := rb[0].(map[string]interface{})

		return &backup.DailyRetentionSchedule{
			RetentionTimes: &times,
			RetentionDuration: &backup.RetentionDuration{
				Count:        utils.Int32(int32(block["count"].(int))),
				DurationType: backup.RetentionDurationTypeDays,
			},
		}
	}

	return nil
}

func flattenArmBackupProtectionPolicyFileShareSchedule(schedule *backup.SimpleSchedulePolicy) []interface{} {
	block := map[string]interface{}{}

	block["frequency"] = string(schedule.ScheduleRunFrequency)

	if times := schedule.ScheduleRunTimes; times != nil && len(*times) > 0 {
		block["time"] = (*times)[0].Format("15:04")
	}

	return []interface{}{block}
}

func flattenArmBackupProtectionPolicyFileShareRetentionDaily(daily *backup.DailyRetentionSchedule) []interface{} {
	block := map[string]interface{}{}

	if duration := daily.RetentionDuration; duration != nil {
		if v := duration.Count; v != nil {
			block["count"] = *v
		}
	}

	return []interface{}{block}
}

func resourceArmBackupProtectionPolicyFileShareWaitForUpdate(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string, d *schema.ResourceData) (backup.ProtectionPolicyResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"NotFound"},
		Target:     []string{"Found"},
		Refresh:    resourceArmBackupProtectionPolicyFileShareRefreshFunc(ctx, client, vaultName, resourceGroup, policyName),
	}

	if d.IsNewResource() {
		state.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		state.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	resp, err := state.WaitForState()
	if err != nil {
		return resp.(backup.ProtectionPolicyResource), fmt.Errorf("Error waiting for the Recovery Service Protection Policy %q to update (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	return resp.(backup.ProtectionPolicyResource), nil
}

func resourceArmBackupProtectionPolicyFileShareWaitForDeletion(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string, d *schema.ResourceData) (backup.ProtectionPolicyResource, error) {
	state := &resource.StateChangeConf{
		MinTimeout: 30 * time.Second,
		Delay:      10 * time.Second,
		Pending:    []string{"Found"},
		Target:     []string{"NotFound"},
		Refresh:    resourceArmBackupProtectionPolicyFileShareRefreshFunc(ctx, client, vaultName, resourceGroup, policyName),
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	resp, err := state.WaitForState()
	if err != nil {
		return resp.(backup.ProtectionPolicyResource), fmt.Errorf("Error waiting for the Recovery Service Protection Policy %q to be missing (Resource Group %q): %+v", policyName, resourceGroup, err)
	}

	return resp.(backup.ProtectionPolicyResource), nil
}

func resourceArmBackupProtectionPolicyFileShareRefreshFunc(ctx context.Context, client *backup.ProtectionPoliciesClient, vaultName, resourceGroup, policyName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, vaultName, resourceGroup, policyName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}

			return resp, "Error", fmt.Errorf("Error making Read request on Recovery Service Protection Policy %q (Resource Group %q): %+v", policyName, resourceGroup, err)
		}

		return resp, "Found", nil
	}
}
