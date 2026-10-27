// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replicated_vm

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			disableProtectionInput := replicationprotecteditems.DisableProtectionInput{
				Properties: replicationprotecteditems.DisableProtectionInputProperties{
					DisableProtectionReason: pointer.To(replicationprotecteditems.DisableProtectionReasonNotSpecified),
					// It's a workaround for https://github.com/hashicorp/pandora/issues/1864
					ReplicationProviderInput: replicationprotecteditems.BaseDisableProtectionProviderSpecificInputImpl{
						InstanceType: string(siterecovery.InstanceTypeDisableProtectionProviderSpecificInput),
					},
				},
			}

			if err = client.DeleteThenPoll(ctx, *id, disableProtectionInput); err != nil {
				return fmt.Errorf("deleting %s : %+v", id.String(), err)
			}

			return nil
		},
	}
}
