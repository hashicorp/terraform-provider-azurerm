// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_quota_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model models.NetAppVolumeQuotaRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(model.VolumeID)
			if err != nil {
				return fmt.Errorf("error parsing volume id %s: %+v", model.VolumeID, err)
			}

			id := volumequotarules.NewVolumeQuotaRuleID(subscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, model.Name)

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

			// Performing some validations that are not possible in the schema
			if errorList := validate.ValidateNetAppVolumeQuotaRule(ctx, pointer.From(volumeID), metadata.Client, pointer.To(model)); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			parameters := volumequotarules.VolumeQuotaRule{
				Location: location.Normalize(model.Location),
				Properties: &volumequotarules.VolumeQuotaRulesProperties{
					QuotaSizeInKiBs: pointer.To(model.QuotaSizeInKiB),
					QuotaType:       pointer.ToEnum[volumequotarules.QuotaType](model.QuotaType),
					QuotaTarget:     pointer.To(model.QuotaTarget),
				},
			}

			if err = client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Waiting for quota rule to be created
			if err := waitForQuotaRuleCreateOrUpdate(ctx, client, id); err != nil {
				return fmt.Errorf("waiting create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
