// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2024-03-01/backuppolicy"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NetAppBackupPolicyResource struct{}

var _ sdk.Resource = NetAppBackupPolicyResource{}

func (r NetAppBackupPolicyResource) ModelObject() interface{} {
	return &netAppModels.NetAppBackupPolicyModel{}
}

func (r NetAppBackupPolicyResource) ResourceType() string {
	return "azurerm_netapp_backup_policy"
}

func (r NetAppBackupPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backuppolicy.ValidateBackupPolicyID
}

func (r NetAppBackupPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeQuotaRuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"tags": commonschema.Tags(),

		"daily_backups_to_keep": {
			Type:         pluginsdk.TypeInt,
			Default:      2,
			Optional:     true,
			ValidateFunc: validation.IntBetween(2, 1019),
		},

		"weekly_backups_to_keep": {
			Type:         pluginsdk.TypeInt,
			Default:      1,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 1019),
		},

		"monthly_backups_to_keep": {
			Type:         pluginsdk.TypeInt,
			Default:      1,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 1019),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},
	}
}

func (r NetAppBackupPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppBackupPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppBackupPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := backuppolicy.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			// Validations
			if errorList := netAppValidate.ValidateNetAppBackupPolicyCombinedRetention(model.DailyBackupsToKeep, model.WeeklyBackupsToKeep, model.MonthlyBackupsToKeep); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.BackupPoliciesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			parameters := backuppolicy.BackupPolicy{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: backuppolicy.BackupPolicyProperties{
					DailyBackupsToKeep:   pointer.To(model.DailyBackupsToKeep),
					WeeklyBackupsToKeep:  pointer.To(model.WeeklyBackupsToKeep),
					MonthlyBackupsToKeep: pointer.To(model.MonthlyBackupsToKeep),
					Enabled:              pointer.To(model.Enabled),
				},
			}

			err = client.BackupPoliciesCreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppBackupPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupPolicyClient

			id, err := backuppolicy.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupPolicyModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			var shouldUpdate = false

			update := backuppolicy.BackupPolicyPatch{
				Properties: &backuppolicy.BackupPolicyProperties{},
			}

			// Checking properties with changes
			if metadata.ResourceData.HasChange("tags") {
				update.Tags = pointer.To(state.Tags)
				shouldUpdate = true
			}

			if metadata.ResourceData.HasChange("daily_backups_to_keep") {
				update.Properties.DailyBackupsToKeep = pointer.To(state.DailyBackupsToKeep)
				shouldUpdate = true
			}

			if metadata.ResourceData.HasChange("weekly_backups_to_keep") {
				update.Properties.WeeklyBackupsToKeep = pointer.To(state.WeeklyBackupsToKeep)
				shouldUpdate = true
			}

			if metadata.ResourceData.HasChange("monthly_backups_to_keep") {
				update.Properties.MonthlyBackupsToKeep = pointer.To(state.MonthlyBackupsToKeep)
				shouldUpdate = true
			}

			if metadata.ResourceData.HasChange("enabled") {
				update.Properties.Enabled = pointer.To(state.Enabled)
				shouldUpdate = true
			}

			if shouldUpdate {
				metadata.Logger.Infof("Updating %s", id)

				if err := client.BackupPoliciesUpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				metadata.SetID(id)
			}

			return nil
		},
	}
}

func (r NetAppBackupPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.BackupPolicyClient

			id, err := backuppolicy.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupPolicyModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.BackupPoliciesGet(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if model := existing.Model; model != nil {
				state.Location = location.NormalizeNilable(pointer.To(model.Location))
				state.Tags = pointer.From(model.Tags)
				state.AccountName = id.NetAppAccountName
				state.Name = id.BackupPolicyName
				state.ResourceGroupName = id.ResourceGroupName

				state.DailyBackupsToKeep = pointer.From(model.Properties.DailyBackupsToKeep)
				state.WeeklyBackupsToKeep = pointer.From(model.Properties.WeeklyBackupsToKeep)
				state.MonthlyBackupsToKeep = pointer.From(model.Properties.MonthlyBackupsToKeep)
				state.Enabled = pointer.From(model.Properties.Enabled)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func (r NetAppBackupPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.BackupPolicyClient

			id, err := backuppolicy.ParseBackupPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.BackupPoliciesGet(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if err = client.BackupPoliciesDeleteThenPoll(ctx, pointer.From(id)); err != nil {
				return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
			}

			if err = waitForBackupPolicyDeletion(ctx, client, *id); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func waitForBackupPolicyDeletion(ctx context.Context, client *backuppolicy.BackupPolicyClient, id backuppolicy.BackupPolicyId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netappBackupPolicyStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func netappBackupPolicyStateRefreshFunc(ctx context.Context, client *backuppolicy.BackupPolicyClient, id backuppolicy.BackupPolicyId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.BackupPoliciesGet(ctx, id)
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
