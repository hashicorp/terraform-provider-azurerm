// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model MsSqlVirtualMachineGroupModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.MSSQL.VirtualMachineGroupsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := sqlvirtualmachinegroups.NewSqlVirtualMachineGroupID(subscriptionId, model.ResourceGroup, model.Name)

			if _, err := client.Get(ctx, id); err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     pointer.To(model.SqlImageOffer),
					SqlImageSku:       pointer.ToEnum[sqlvirtualmachinegroups.SqlVMGroupImageSku](model.SqlImageSku),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			return nil
		},
	}
}
