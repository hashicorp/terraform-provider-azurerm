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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/failovergroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.FailoverGroupsClient
			serversClient := metadata.Client.MSSQL.ServersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MsSqlFailoverGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serverId, err := commonids.ParseSqlServerID(model.ServerId)
			if err != nil {
				return err
			}

			if _, err = serversClient.Get(ctx, *serverId, servers.DefaultGetOperationOptions()); err != nil {
				return fmt.Errorf("retrieving %s: %+v", serverId, err)
			}

			id := failovergroups.NewFailoverGroupID(subscriptionId, serverId.ResourceGroupName, serverId.ServerName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			readOnlyFailoverPolicy := failovergroups.ReadOnlyEndpointFailoverPolicyDisabled
			if model.ReadonlyEndpointFailurePolicyEnabled {
				readOnlyFailoverPolicy = failovergroups.ReadOnlyEndpointFailoverPolicyEnabled
			}

			properties := failovergroups.FailoverGroup{
				Properties: &failovergroups.FailoverGroupProperties{
					Databases: &model.Databases,
					ReadOnlyEndpoint: &failovergroups.FailoverGroupReadOnlyEndpoint{
						FailoverPolicy: &readOnlyFailoverPolicy,
					},
					ReadWriteEndpoint: failovergroups.FailoverGroupReadWriteEndpoint{},
					PartnerServers:    r.expandPartnerServers(model.PartnerServers),
				},
				Tags: pointer.To(model.Tags),
			}

			if rwPolicy := model.ReadWriteEndpointFailurePolicy; len(rwPolicy) > 0 {
				properties.Properties.ReadWriteEndpoint.FailoverPolicy = failovergroups.ReadWriteEndpointFailoverPolicy(rwPolicy[0].Mode)
				if rwPolicy[0].Mode == string(failovergroups.ReadWriteEndpointFailoverPolicyAutomatic) {
					properties.Properties.ReadWriteEndpoint.FailoverWithDataLossGracePeriodMinutes = pointer.To(rwPolicy[0].GraceMinutes)
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandPartnerServers(input []PartnerServerModel) []failovergroups.PartnerInfo {
	partnerServers := make([]failovergroups.PartnerInfo, 0, len(input))
	if input == nil {
		return partnerServers
	}

	for _, v := range input {
		partnerServers = append(partnerServers, failovergroups.PartnerInfo{
			Id: v.ID,
		})
	}

	return partnerServers
}
