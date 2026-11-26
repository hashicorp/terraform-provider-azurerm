// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_failover_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/instancefailovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceFailoverGroupsClient

			id, err := instancefailovergroups.ParseInstanceFailoverGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MsSqlManagedInstanceFailoverGroupModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := MsSqlManagedInstanceFailoverGroupModel{
				Name:     id.InstanceFailoverGroupName,
				Location: id.LocationName,
			}

			if result.Model != nil {
				if props := result.Model.Properties; props != nil {
					model.Role = pointer.FromEnum(props.ReplicationRole)

					if instancePairs := props.ManagedInstancePairs; len(instancePairs) == 1 {
						if primaryId := instancePairs[0].PrimaryManagedInstanceId; primaryId != nil {
							id, err := commonids.ParseSqlManagedInstanceIDInsensitively(*primaryId)
							if err != nil {
								return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
							}

							model.ManagedInstanceId = id.ID()
						}

						if partnerId := instancePairs[0].PartnerManagedInstanceId; partnerId != nil {
							id, err := commonids.ParseSqlManagedInstanceIDInsensitively(*partnerId)
							if err != nil {
								return fmt.Errorf("parsing `PrimaryManagedInstanceID` from response: %v", err)
							}

							model.PartnerManagedInstanceId = id.ID()
						}
					}

					for _, partnerRegion := range props.PartnerRegions {
						model.PartnerRegion = append(model.PartnerRegion, MsSqlManagedInstancePartnerRegionModel{
							Location: pointer.From(partnerRegion.Location),
							Role:     pointer.FromEnum(partnerRegion.ReplicationRole),
						})
					}

					if readOnlyEndpoint := props.ReadOnlyEndpoint; readOnlyEndpoint != nil {
						if *readOnlyEndpoint.FailoverPolicy == instancefailovergroups.ReadOnlyEndpointFailoverPolicyEnabled {
							model.ReadOnlyEndpointFailoverPolicyEnabled = true
						}
					}

					model.SecondaryType = pointer.FromEnum(props.SecondaryType)

					model.ReadWriteEndpointFailurePolicy = []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel{
						{
							Mode:         string(props.ReadWriteEndpoint.FailoverPolicy),
							GraceMinutes: pointer.From(props.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes),
						},
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}
