// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_failover_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient

			id, err := failovergroups.ParseFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MsSqlFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			readOnlyFailoverPolicy := failovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = failovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := failovergroups.FailoverGroup{
				Properties: &failovergroups.FailoverGroupProperties{
					Databases: &state.Databases,
					ReadOnlyEndpoint: &failovergroups.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: failovergroups.FailoverGroupReadWriteEndpoint{
						FailoverPolicy: failovergroups.ReadWriteEndpointFailoverPolicy(state.ReadWriteEndpointFailurePolicy[0].Mode),
					},
					PartnerServers: r.expandPartnerServers(state.PartnerServers),
				},
				Tags: pointer.To(state.Tags),
			}

			if state.ReadWriteEndpointFailurePolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
				properties.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = pointer.To(state.ReadWriteEndpointFailurePolicy[0].GraceMinutes)
			}

			// client.Update doesn't support changing the PartnerServers
			if err = client.CreateOrUpdateThenPoll(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
