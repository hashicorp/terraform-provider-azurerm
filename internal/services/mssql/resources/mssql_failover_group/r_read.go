// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_failover_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/failovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.MSSQL.FailoverGroupsClient

			id, err := failovergroups.ParseFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			serverId := commonids.NewSqlServerID(subscriptionId, id.ResourceGroupName, id.ServerName)

			model := MsSqlFailoverGroupModel{
				Name:     id.FailoverGroupName,
				ServerId: serverId.ID(),
			}

			if existing.Model != nil {
				model.Tags = pointer.From(existing.Model.Tags)

				if props := existing.Model.Properties; props != nil {
					if props.Databases != nil {
						model.Databases = *props.Databases
					}

					model.PartnerServers = r.flattenPartnerServers(props.PartnerServers)

					if props.ReadOnlyEndpoint != nil && pointer.From(props.ReadOnlyEndpoint.FailoverPolicy) == failovergroups.ReadOnlyEndpointFailoverPolicyEnabled {
						model.ReadonlyEndpointFailurePolicyEnabled = true
					}

					model.ReadWriteEndpointFailurePolicy = []ReadWriteEndpointFailurePolicyModel{{
						Mode: string(props.ReadWriteEndpoint.FailoverPolicy),
					}}

					model.ReadWriteEndpointFailurePolicy[0].GraceMinutes = pointer.From(props.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes)
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r Resource) flattenPartnerServers(input []failovergroups.PartnerInfo) []PartnerServerModel {
	output := make([]PartnerServerModel, 0)
	if input == nil {
		return output
	}

	for _, partner := range input {
		model := PartnerServerModel{
			Location: location.NormalizeNilable(partner.Location),
			Role:     pointer.FromEnum(partner.ReplicationRole),
			ID:       partner.Id,
		}

		output = append(output, model)
	}

	return output
}
