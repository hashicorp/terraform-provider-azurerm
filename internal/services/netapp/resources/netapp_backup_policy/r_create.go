// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_backup_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/backuppolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupPolicyClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model models.NetAppBackupPolicyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := backuppolicies.NewBackupPolicyID(subscriptionId, model.ResourceGroupName, model.AccountName, model.Name)

			// Validations
			if errorList := validate.ValidateNetAppBackupPolicyCombinedRetention(model.DailyBackupsToKeep, model.WeeklyBackupsToKeep, model.MonthlyBackupsToKeep); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %s", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			parameters := backuppolicies.BackupPolicy{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Properties: backuppolicies.BackupPolicyProperties{
					DailyBackupsToKeep:   pointer.To(model.DailyBackupsToKeep),
					WeeklyBackupsToKeep:  pointer.To(model.WeeklyBackupsToKeep),
					MonthlyBackupsToKeep: pointer.To(model.MonthlyBackupsToKeep),
					Enabled:              pointer.To(model.Enabled),
				},
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
