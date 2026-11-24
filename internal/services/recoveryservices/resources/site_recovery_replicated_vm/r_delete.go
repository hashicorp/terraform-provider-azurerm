// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replicated_vm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryReplicatedItemDelete(d *pluginsdk.ResourceData, meta any) error {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	disableProtectionInput := replicationprotecteditems.DisableProtectionInput{
		Properties: replicationprotecteditems.DisableProtectionInputProperties{
			DisableProtectionReason: pointer.To(replicationprotecteditems.DisableProtectionReasonNotSpecified),
			// It's a workaround for https://github.com/hashicorp/pandora/issues/1864
			ReplicationProviderInput: replicationprotecteditems.BaseDisableProtectionProviderSpecificInputImpl{
				InstanceType: string(siterecovery.InstanceTypeDisableProtectionProviderSpecificInput),
			},
		},
	}

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	if err = client.DeleteThenPoll(ctx, *id, disableProtectionInput); err != nil {
		return fmt.Errorf("deleting site recovery replicated vm %s : %+v", id.String(), err)
	}

	return nil
}
