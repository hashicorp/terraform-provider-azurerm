// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/tableservice"
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
	Name             string     `tfschema:"name"`
	StorageAccountId string     `tfschema:"storage_account_id"`
	ACL              []ACLModel `tfschema:"acl"`
	Id               string     `tfschema:"id"`
	URL              string     `tfschema:"url"`

	// TODO 5.0: Remove this
	StorageAccountName string `tfschema:"storage_account_name"`

	// TODO 5.0: Remove this
	ResourceManagerId string `tfschema:"resource_manager_id"`
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
	r := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageTableName,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageAccountName,
		},
	}

	if !features.FivePointOh() {
		r["storage_account_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.StorageAccountName,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
			Deprecated:   "the `storage_account_name` property has been deprecated in favour of `storage_account_id` and will be removed in version 5.0 of the Provider.",
		}

		r["storage_account_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
		}
	}

	return r
}

func (k storageTableDataSource) Attributes() map[string]*pluginsdk.Schema {
	r := map[string]*pluginsdk.Schema{
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

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}

	if !features.FivePointOh() {
		r["resource_manager_id"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Computed: true,
		}
	}

	return r
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

			if !features.FivePointOh() {
				if model.StorageAccountName != "" {
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

					model.ACL = flattenStorageTableACLsWithMetadataDeprecated(acls)

					resourceManagerId := parse.NewStorageTableResourceManagerID(account.StorageAccountId.SubscriptionId, account.StorageAccountId.ResourceGroupName, account.StorageAccountId.StorageAccountName, "default", model.Name)
					model.ResourceManagerId = resourceManagerId.ID()
					model.URL = id.ID()
					metadata.SetID(id)

					return metadata.Encode(&model)
				}
			}

			tableClient := metadata.Client.Storage.ResourceManager.TableService
			accountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			id := tableservice.NewTableID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, model.Name)
			resp, err := tableClient.TableGet(ctx, id)
			if err != nil {
				return fmt.Errorf("retrieving %q: %v", id, err)
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					acl, err := flattenStorageTableACLsWithMetadata(props.SignedIdentifiers)
					if err != nil {
						return fmt.Errorf("flattening `acl`: %v", err)
					}
					model.ACL = acl
				}
			}

			account, err := metadata.Client.Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
			if err != nil {
				return fmt.Errorf("retrieving Account for Table %q: %v", id, err)
			}
			// Determine the table endpoint, so we can build a data plane ID
			endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
			if err != nil {
				return fmt.Errorf("determining Table endpoint: %v", err)
			}
			// Parse the table endpoint as a data plane account ID
			accountDpId, err := accounts.ParseAccountID(*endpoint, metadata.Client.Storage.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}
			model.URL = tables.NewTableID(*accountDpId, id.TableName).ID()

			if !features.FivePointOh() {
				model.ResourceManagerId = id.ID()
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func flattenStorageTableACLsWithMetadata(acls *[]tableservice.TableSignedIdentifier) ([]ACLModel, error) {
	if acls == nil {
		return []ACLModel{}, nil
	}

	output := make([]ACLModel, 0, len(*acls))
	for _, acl := range *acls {
		var startTime, expiryTime, permission string
		var err error
		if policy := acl.AccessPolicy; policy != nil {
			if policy.StartTime != nil {
				startTime = *policy.StartTime
				startTime, err = convertTimeFormat(startTime)
				if err != nil {
					return nil, err
				}
			}
			if policy.ExpiryTime != nil {
				expiryTime = *policy.ExpiryTime
				expiryTime, err = convertTimeFormat(expiryTime)
				if err != nil {
					return nil, err
				}
			}
			permission = policy.Permission
		}

		var accessPolicies []AccessPolicyModel
		accessPolicies = append(accessPolicies, AccessPolicyModel{
			Start:       startTime,
			Expiry:      expiryTime,
			Permissions: permission,
		})

		output = append(output, ACLModel{
			Id:           acl.Id,
			AccessPolicy: accessPolicies,
		})
	}

	return output, nil
}

func flattenStorageTableACLsWithMetadataDeprecated(acls *[]tables.SignedIdentifier) []ACLModel {
	if acls == nil {
		return []ACLModel{}
	}

	output := make([]ACLModel, 0, len(*acls))
	for _, acl := range *acls {
		var accessPolicies []AccessPolicyModel
		policy := acl.AccessPolicy
		accessPolicies = append(accessPolicies, AccessPolicyModel{
			Start:       policy.Start,
			Expiry:      policy.Expiry,
			Permissions: policy.Permission,
		})

		output = append(output, ACLModel{
			Id:           acl.Id,
			AccessPolicy: accessPolicies,
		})
	}

	return output
}
