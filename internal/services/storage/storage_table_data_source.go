// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

type storageTableDataSource struct{}

var _ sdk.DataSource = storageTableDataSource{}

type TableDataSourceModel struct {
	Name               string     `tfschema:"name"`
	StorageAccountName string     `tfschema:"storage_account_name"`
	ACL                []ACLModel `tfschema:"acl"`
	Id                 string     `tfschema:"id"`
	ResourceManagerId  string     `tfschema:"resource_manager_id"`
}

type ACLModel struct {
	Id           string              `tfschema:"id"`
	AccessPolicy []AccessPolicyModel `tfschema:"access_policy"`
}

type AccessPolicyModel struct {
	Start       string `tfschema:"start"`
	Expiry      string `tfschema:"expiry"`
	Permissions string `tfschema:"permissions"`
}

func (k storageTableDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageTableName,
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageAccountName,
		},
	}
}

func (k storageTableDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"acl": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"access_policy": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"expiry": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"permissions": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_manager_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (k storageTableDataSource) ModelObject() interface{} {
	return &TableDataSourceModel{}
}

func (k storageTableDataSource) ResourceType() string {
	return "azurerm_storage_table"
}

func (k storageTableDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model TableDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			storageClient := metadata.Client.Storage

			account, err := storageClient.FindAccount(ctx, metadata.Client.Account.SubscriptionId, model.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", model.StorageAccountName, model.Name, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q for Table %q", model.StorageAccountName, model.Name)
			}

			// Determine the table endpoint, so we can build a data plane ID
			endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
			if err != nil {
				return fmt.Errorf("determining Table endpoint: %v", err)
			}

			// Parse the table endpoint as a data plane account ID
			accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}

			id := tables.NewTableID(*accountId, model.Name)

			aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
			if err != nil {
				return fmt.Errorf("building Tables Client: %v", err)
			}

			acls, err := aclClient.GetACLs(ctx, model.Name)
			if err != nil {
				return fmt.Errorf("retrieving ACLs for %s: %v", id, err)
			}

			model.ACL = flattenStorageTableACLsWithMetadata(acls)

			resourceManagerId := parse.NewStorageTableResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", model.Name)
			model.ResourceManagerId = resourceManagerId.ID()
			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func flattenStorageTableACLsWithMetadata(acls *[]tables.SignedIdentifier) []ACLModel {
	if acls == nil {
		return []ACLModel{}
	}

	output := make([]ACLModel, 0, len(*acls))
	for _, acl := range *acls {
		var accessPolicies []AccessPolicyModel
		for _, policy := range []tables.AccessPolicy{acl.AccessPolicy} {
			accessPolicies = append(accessPolicies, AccessPolicyModel{
				Start:       policy.Start,
				Expiry:      policy.Expiry,
				Permissions: policy.Permission,
			})
		}

		output = append(output, ACLModel{
			Id:           acl.Id,
			AccessPolicy: accessPolicies,
		})
	}

	return output
}
