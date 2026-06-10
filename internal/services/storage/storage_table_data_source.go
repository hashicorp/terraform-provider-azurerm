// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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
	StorageAccountName string     `tfschema:"storage_account_name,removedInNextMajorVersion"`
	StorageAccountId   string     `tfschema:"storage_account_id"`
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
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageTableName,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},
	}

	if !features.FivePointOh() {
		s["storage_account_id"].Required = false
		s["storage_account_id"].Optional = true
		s["storage_account_id"].Computed = true
		s["storage_account_id"].ConflictsWith = []string{"storage_account_name"}

		s["storage_account_name"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"storage_account_id"},
			Deprecated:    "`storage_account_name` has been deprecated in favour of `storage_account_id` and will be removed in v5.0 of the AzureRM Provider",
		}
	}

	return s
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

			var accountName string
			var account *client.AccountDetails
			var err error

			if !features.FivePointOh() {
				if model.StorageAccountId != "" {
					storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
					if err != nil {
						return fmt.Errorf("parsing storage_account_id: %v", err)
					}
					accountName = storageAccountId.StorageAccountName
					account, err = storageClient.GetAccount(ctx, *storageAccountId)
					if err != nil {
						return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, model.Name, err)
					}
				} else {
					accountName = model.StorageAccountName
					account, err = storageClient.FindAccount(ctx, metadata.Client.Account.SubscriptionId, accountName)
					if err != nil {
						return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, model.Name, err)
					}
				}
			} else {
				storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
				if err != nil {
					return fmt.Errorf("parsing storage_account_id: %v", err)
				}
				accountName = storageAccountId.StorageAccountName
				account, err = storageClient.GetAccount(ctx, *storageAccountId)
				if err != nil {
					return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", accountName, model.Name, err)
				}
			}

			if account == nil {
				return fmt.Errorf("locating Storage Account %q for Table %q", accountName, model.Name)
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

			if !features.FivePointOh() {
				metadata.SetID(id)
				model.StorageAccountName = accountName
			} else {
				metadata.SetID(resourceManagerId)
			}

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
