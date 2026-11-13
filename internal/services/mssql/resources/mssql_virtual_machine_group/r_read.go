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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.VirtualMachineGroupsClient

			id, err := sqlvirtualmachinegroups.ParseSqlVirtualMachineGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := MsSqlVirtualMachineGroupModel{
				Name:          id.SqlVirtualMachineGroupName,
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.SqlImageOffer = pointer.From(props.SqlImageOffer)
					state.SqlImageSku = pointer.FromEnum(props.SqlImageSku)

					var oldModel MsSqlVirtualMachineGroupModel
					if err = metadata.Decode(&oldModel); err != nil {
						return err
					}
					storageAccountPrimaryKey := ""
					if len(oldModel.WsfcDomainProfile) != 0 {
						storageAccountPrimaryKey = oldModel.WsfcDomainProfile[0].StorageAccountPrimaryKey
					}
					state.WsfcDomainProfile = flattenMsSqlVirtualMachineGroupWsfcDomainProfile(props.WsfcDomainProfile, storageAccountPrimaryKey)
				}
				state.Location = location.Normalize(model.Location)

				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenMsSqlVirtualMachineGroupWsfcDomainProfile(domainProfile *sqlvirtualmachinegroups.WsfcDomainProfile, storageAccountPrimaryKey string) []WsfcDomainProfile {
	if domainProfile == nil {
		return []WsfcDomainProfile{}
	}

	return []WsfcDomainProfile{
		{
			Fqdn:                        pointer.From(domainProfile.DomainFqdn),
			OrganizationalUnitPath:      pointer.From(domainProfile.OuPath),
			ClusterBootstrapAccountName: pointer.From(domainProfile.ClusterBootstrapAccount),
			ClusterOperatorAccountName:  pointer.From(domainProfile.ClusterOperatorAccount),
			SqlServiceAccountName:       pointer.From(domainProfile.SqlServiceAccount),
			StorageAccountUrl:           pointer.From(domainProfile.StorageAccountURL),
			ClusterSubnetType:           pointer.FromEnum(domainProfile.ClusterSubnetType),
			StorageAccountPrimaryKey:    storageAccountPrimaryKey,
		},
	}
}

type WsfcDomainProfile struct {
	Fqdn                        string `tfschema:"fqdn"`
	OrganizationalUnitPath      string `tfschema:"organizational_unit_path"`
	ClusterBootstrapAccountName string `tfschema:"cluster_bootstrap_account_name"`
	ClusterOperatorAccountName  string `tfschema:"cluster_operator_account_name"`
	SqlServiceAccountName       string `tfschema:"sql_service_account_name"`
	StorageAccountUrl           string `tfschema:"storage_account_url"`
	StorageAccountPrimaryKey    string `tfschema:"storage_account_primary_key"`
	ClusterSubnetType           string `tfschema:"cluster_subnet_type"`
}
