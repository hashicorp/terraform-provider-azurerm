// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_hyperv_replication_policy_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	recoveryservicesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
)

func (h Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := recoveryservicesClient.From(metadata.Client).ContainerMappingClient

			id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			if err := client.DeleteThenPoll(ctx, *id, replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
