// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_vmware_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.RecoveryServices.ContainerMappingClient

			input := replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{
				Properties: &replicationprotectioncontainermappings.RemoveProtectionContainerMappingInputProperties{
					ProviderSpecificInput: &replicationprotectioncontainermappings.ReplicationProviderContainerUnmappingInput{},
				},
			}

			if err = client.DeleteThenPoll(ctx, *id, input); err != nil {
				return fmt.Errorf("deleting %s : %+v", id.String(), err)
			}

			return nil
		},
	}
}
