// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			parameters := sqlvirtualmachinegroups.SqlVirtualMachineGroup{
				Properties: &sqlvirtualmachinegroups.SqlVirtualMachineGroupProperties{
					SqlImageOffer:     pointer.To(model.SqlImageOffer),
					SqlImageSku:       pointer.ToEnum[sqlvirtualmachinegroups.SqlVMGroupImageSku](model.SqlImageSku),
					WsfcDomainProfile: expandMsSqlVirtualMachineGroupWsfcDomainProfile(model.WsfcDomainProfile),
				},

				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandMsSqlVirtualMachineGroupWsfcDomainProfile(wsfcDomainProfile []WsfcDomainProfile) *sqlvirtualmachinegroups.WsfcDomainProfile {
	if wsfcDomainProfile == nil {
		return nil
	}

	result := sqlvirtualmachinegroups.WsfcDomainProfile{
		ClusterSubnetType:        pointer.ToEnum[sqlvirtualmachinegroups.ClusterSubnetType](wsfcDomainProfile[0].ClusterSubnetType),
		DomainFqdn:               pointer.To(wsfcDomainProfile[0].Fqdn),
		OuPath:                   pointer.To(wsfcDomainProfile[0].OrganizationalUnitPath),
		ClusterBootstrapAccount:  pointer.To(wsfcDomainProfile[0].ClusterBootstrapAccountName),
		ClusterOperatorAccount:   pointer.To(wsfcDomainProfile[0].ClusterOperatorAccountName),
		SqlServiceAccount:        pointer.To(wsfcDomainProfile[0].SqlServiceAccountName),
		StorageAccountURL:        pointer.To(wsfcDomainProfile[0].StorageAccountUrl),
		StorageAccountPrimaryKey: pointer.To(wsfcDomainProfile[0].StorageAccountPrimaryKey),
	}

	return &result
}
