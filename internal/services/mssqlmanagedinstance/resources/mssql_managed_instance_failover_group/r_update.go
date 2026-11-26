// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_failover_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/instancefailovergroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
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

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(state.ManagedInstanceId)
			if err != nil {
				return err
			}

			// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
			// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
			partnerId, err := commonids.ParseSqlManagedInstanceIDInsensitively(state.PartnerManagedInstanceId)
			if err != nil {
				return err
			}

			instancesClient := metadata.Client.MSSQLManagedInstance.ManagedInstancesClientForSubscription(partnerId.SubscriptionId)
			partner, err := instancesClient.Get(ctx, partnerId.ResourceGroupName, partnerId.ManagedInstanceName, "")
			if err != nil || partner.Location == nil || *partner.Location == "" {
				return fmt.Errorf("checking for existence and region of Partner of %q: %+v", id, err)
			}

			readOnlyFailoverPolicy := instancefailovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if state.ReadOnlyEndpointFailoverPolicyEnabled {
				readOnlyFailoverPolicy = instancefailovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			parameters := instancefailovergroups.InstanceFailoverGroup{
				Properties: &instancefailovergroups.InstanceFailoverGroupProperties{
					ReadOnlyEndpoint: &instancefailovergroups.InstanceFailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: instancefailovergroups.InstanceFailoverGroupReadWriteEndpoint{},
					PartnerRegions: []instancefailovergroups.PartnerRegionInfo{
						{
							Location: partner.Location,
						},
					},
					ManagedInstancePairs: []instancefailovergroups.ManagedInstancePairInfo{
						{
							PrimaryManagedInstanceId: pointer.To(managedInstanceId.ID()),
							PartnerManagedInstanceId: pointer.To(partnerId.ID()),
						},
					},
					SecondaryType: pointer.ToEnum[instancefailovergroups.SecondaryInstanceType](state.SecondaryType),
				},
			}

			if rwPolicy := state.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				parameters.Properties.ReadWriteEndpoint.FailoverPolicy = instancefailovergroups.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(instancefailovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
					parameters.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = &rwPolicy[0].GraceMinutes
				}
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
