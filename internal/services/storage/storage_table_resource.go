// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-05-01/tableservice"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func resourceStorageTable() *pluginsdk.Resource {
	r := &pluginsdk.Resource{
		Create: resourceStorageTableCreate,
		Read:   resourceStorageTableRead,
		Delete: resourceStorageTableDelete,
		Update: resourceStorageTableUpdate,

		Importer: helpers.ImporterValidatingStorageResourceId(func(id, storageDomainSuffix string) error {
			if !features.FivePointOh() {
				if strings.HasPrefix(id, "/subscriptions/") {
					_, err := tableservice.ParseTableID(id)
					return err
				}
				_, err := tables.ParseTableID(id, storageDomainSuffix)
				return err
			}

			_, err := tableservice.ParseTableID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.TableV0ToV1{},
			1: migration.TableV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageTableName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"acl": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"access_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"permissions": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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
		},
	}

	if !features.FivePointOh() {
		r.Schema["storage_account_name"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.StorageAccountName,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
			Deprecated:   "the `storage_account_name` property has been deprecated in favour of `storage_account_id` and will be removed in version 5.0 of the Provider.",
		}

		r.Schema["storage_account_id"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			ExactlyOneOf: []string{"storage_account_id", "storage_account_name"},
		}

		r.Schema["resource_manager_id"] = &pluginsdk.Schema{
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Resource Manager ID of this Storage Table.",
		}

		r.CustomizeDiff = func(ctx context.Context, diff *pluginsdk.ResourceDiff, i interface{}) error {
			// Resource Manager ID in use, but change to `storage_account_id` should recreate
			if strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_id") {
				return diff.ForceNew("storage_account_id")
			}

			// using legacy Data Plane ID but attempting to change the storage_account_name should recreate
			if diff.Id() != "" && !strings.HasPrefix(diff.Id(), "/subscriptions/") && diff.HasChange("storage_account_name") {
				// converting from storage_account_id to the deprecated storage_account_name is not supported
				oldAccountId, _ := diff.GetChange("storage_account_id")
				oldName, newName := diff.GetChange("storage_account_name")

				if oldAccountId.(string) != "" && newName.(string) != "" {
					return diff.ForceNew("storage_account_name")
				}

				if oldName.(string) != "" && newName.(string) != "" {
					return diff.ForceNew("storage_account_name")
				}
			}

			return nil
		}
	}

	return r
}

func resourceStorageTableCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	tableClient := meta.(*clients.Client).Storage.ResourceManager.TableService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tableName := d.Get("name").(string)

	if !features.FivePointOh() {
		if accountName := d.Get("storage_account_name").(string); accountName != "" {
			storageClient := meta.(*clients.Client).Storage

			aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
			acls := expandStorageTableACLsDeprecated(aclsRaw)

			account, err := storageClient.FindAccount(ctx, subscriptionId, accountName)
			if err != nil {
				return fmt.Errorf("retrieving Account %q for Table %q: %s", accountName, tableName, err)
			}
			if account == nil {
				return fmt.Errorf("locating Storage Account %q", accountName)
			}

			tablesDataPlaneClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Tables Client: %s", err)
			}

			// Determine the table endpoint, so we can build a data plane ID
			endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
			if err != nil {
				return fmt.Errorf("determining Tables endpoint: %v", err)
			}

			// Parse the table endpoint as a data plane account ID
			accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
			if err != nil {
				return fmt.Errorf("parsing Account ID: %v", err)
			}

			id := tables.NewTableID(*accountId, tableName)

			exists, err := tablesDataPlaneClient.Exists(ctx, tableName)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %v", id, err)
			}
			if exists != nil && *exists {
				return tf.ImportAsExistsError("azurerm_storage_table", id.ID())
			}

			if err = tablesDataPlaneClient.Create(ctx, tableName); err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			d.SetId(id.ID())

			// Setting ACLs only supports shared key authentication (@manicminer, 2024-02-29)
			aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
			if err != nil {
				return fmt.Errorf("building Tables Client: %v", err)
			}

			if err = aclClient.UpdateACLs(ctx, tableName, acls); err != nil {
				return fmt.Errorf("setting ACLs for %s: %v", id, err)
			}

			return resourceStorageTableRead(d, meta)
		}
	}

	accountId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	id := tableservice.NewTableID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, tableName)

	existing, err := tableClient.TableGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %q: %v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_table", id.ID())
	}

	payload := tableservice.Table{
		Properties: &tableservice.TableProperties{
			SignedIdentifiers: pointer.To(expandStorageTableACLs(d.Get("acl").(*pluginsdk.Set).List())),
		},
	}

	if _, err := tableClient.TableCreate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	d.SetId(id.ID())

	return resourceStorageTableRead(d, meta)
}

func resourceStorageTableRead(d *pluginsdk.ResourceData, meta interface{}) error {
	tableClient := meta.(*clients.Client).Storage.ResourceManager.TableService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() {
		if !strings.HasPrefix(d.Id(), "/subscriptions/") {
			if said := d.Get("storage_account_id").(string); said == "" {
				storageClient := meta.(*clients.Client).Storage

				id, err := tables.ParseTableID(d.Id(), storageClient.StorageDomainSuffix)
				if err != nil {
					return err
				}

				account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
				if err != nil {
					return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", id.AccountId.AccountName, id.TableName, err)
				}
				if account == nil {
					log.Printf("Unable to determine Resource Group for Storage Table %q (Account %s) - assuming removed & removing from state", id.TableName, id.AccountId.AccountName)
					d.SetId("")
					return nil
				}

				client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
				if err != nil {
					return fmt.Errorf("building Tables Client: %v", err)
				}

				exists, err := client.Exists(ctx, id.TableName)
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				if exists == nil || !*exists {
					log.Printf("[DEBUG] %s not found, removing from state", id)
					d.SetId("")
					return nil
				}

				// Retrieving ACLs only supports shared key authentication (@manicminer, 2024-02-29)
				aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
				if err != nil {
					return fmt.Errorf("building Tables Client: %v", err)
				}

				acls, err := aclClient.GetACLs(ctx, id.TableName)
				if err != nil {
					return fmt.Errorf("retrieving ACLs for %s: %v", id, err)
				}

				d.Set("name", id.TableName)
				d.Set("storage_account_name", id.AccountId.AccountName)
				d.Set("resource_manager_id", parse.NewStorageTableResourceManagerID(subscriptionId, account.StorageAccountId.ResourceGroupName, id.AccountId.AccountName, "default", id.TableName).ID())
				d.Set("url", id.ID())

				if err = d.Set("acl", flattenStorageTableACLsDeprecated(acls)); err != nil {
					return fmt.Errorf("setting `acl`: %v", err)
				}

				return nil
			} else {
				// Deal with the ID changing if the user changes from `storage_account_name` to `storage_account_id`
				accountId, err := commonids.ParseStorageAccountID(said)
				if err != nil {
					return err
				}

				id := tableservice.NewTableID(subscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, d.Get("name").(string))
				d.SetId(id.ID())
				// Continue the code flow outside this block
			}
		}
	}

	id, err := tableservice.ParseTableID(d.Id())
	if err != nil {
		return err
	}

	existing, err := tableClient.TableGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found, removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %v", *id, err)
	}

	d.Set("name", id.TableName)
	d.Set("storage_account_id", commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID())

	if !features.FivePointOh() {
		d.Set("resource_manager_id", id.String())
		d.Set("storage_account_name", "")
	}

	if model := existing.Model; model != nil {
		if prop := model.Properties; prop != nil {
			if acls := prop.SignedIdentifiers; acls != nil {
				acl, err := flattenStorageTableACLs(*acls)
				if err != nil {
					return fmt.Errorf("flattening `acl`: %v", err)
				}
				if err := d.Set("acl", acl); err != nil {
					return fmt.Errorf("setting `acl`: %s", err)
				}
			}
		}
	}

	account, err := meta.(*clients.Client).Storage.GetAccount(ctx, commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName))
	if err != nil {
		return fmt.Errorf("retrieving Account for Table %q: %v", id, err)
	}

	// Determine the table endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeTable)
	if err != nil {
		return fmt.Errorf("determining Table endpoint: %v", err)
	}

	// Parse the table endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, meta.(*clients.Client).Storage.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	d.Set("url", tables.NewTableID(*accountId, id.TableName).ID())

	return nil
}

func resourceStorageTableDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	tableClient := meta.(*clients.Client).Storage.ResourceManager.TableService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage
		id, err := tables.ParseTableID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

		account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
		if err != nil {
			return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", id.AccountId.AccountName, id.TableName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
		}

		client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Tables Client: %v", err)
		}

		if err = client.Delete(ctx, id.TableName); err != nil {
			if strings.Contains(err.Error(), "unexpected status 404") {
				return nil
			}
			return fmt.Errorf("deleting %s: %v", id, err)
		}

		return nil
	}

	id, err := tableservice.ParseTableID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := tableClient.TableDelete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}

func resourceStorageTableUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tableClient := meta.(*clients.Client).Storage.ResourceManager.TableService
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if !features.FivePointOh() && !strings.HasPrefix(d.Id(), "/subscriptions/") {
		storageClient := meta.(*clients.Client).Storage
		id, err := tables.ParseTableID(d.Id(), storageClient.StorageDomainSuffix)
		if err != nil {
			return err
		}

		account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
		if err != nil {
			return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", id.AccountId.AccountName, id.TableName, err)
		}
		if account == nil {
			return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
		}

		if d.HasChange("acl") {
			log.Printf("[DEBUG] Updating ACLs for %s", id)

			aclsRaw := d.Get("acl").(*pluginsdk.Set).List()
			acls := expandStorageTableACLsDeprecated(aclsRaw)

			// Setting ACLs only supports shared key authentication (@manicminer, 2024-02-29)
			aclClient, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingOnlySharedKeyAuth())
			if err != nil {
				return fmt.Errorf("building Tables Client: %v", err)
			}

			if err = aclClient.UpdateACLs(ctx, id.TableName, acls); err != nil {
				return fmt.Errorf("updating ACLs for %s: %v", id, err)
			}

			log.Printf("[DEBUG] Updated ACLs for %s", id)
		}

		return resourceStorageTableRead(d, meta)
	}

	id, err := tableservice.ParseTableID(d.Id())
	if err != nil {
		return err
	}

	existing, err := tableClient.TableGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %q: %v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("unexpected null model after retrieving %v", id)
	}

	payload := tableservice.Table{
		Properties: existing.Model.Properties,
	}

	if d.HasChange("acl") {
		payload.Properties.SignedIdentifiers = pointer.To(expandStorageTableACLs(d.Get("acl").(*pluginsdk.Set).List()))
	}

	if _, err := tableClient.TableCreate(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %v", id, err)
	}

	return resourceStorageTableRead(d, meta)
}

func expandStorageTableACLs(input []interface{}) []tableservice.TableSignedIdentifier {
	results := make([]tableservice.TableSignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := tableservice.TableSignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: &tableservice.TableAccessPolicy{
				StartTime:  pointer.To(policy["start"].(string)),
				ExpiryTime: pointer.To(policy["expiry"].(string)),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func expandStorageTableACLsDeprecated(input []interface{}) []tables.SignedIdentifier {
	results := make([]tables.SignedIdentifier, 0)

	for _, v := range input {
		vals := v.(map[string]interface{})

		policies := vals["access_policy"].([]interface{})
		policy := policies[0].(map[string]interface{})

		identifier := tables.SignedIdentifier{
			Id: vals["id"].(string),
			AccessPolicy: tables.AccessPolicy{
				Start:      policy["start"].(string),
				Expiry:     policy["expiry"].(string),
				Permission: policy["permissions"].(string),
			},
		}
		results = append(results, identifier)
	}

	return results
}

func flattenStorageTableACLs(input []tableservice.TableSignedIdentifier) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, v := range input {
		var startTime, expiryTime string
		var err error
		if policy := v.AccessPolicy; policy != nil {
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
		}
		output := map[string]interface{}{
			"id": v.Id,
			"access_policy": []interface{}{
				map[string]interface{}{
					"start":       startTime,
					"expiry":      expiryTime,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result, nil
}

func flattenStorageTableACLsDeprecated(input *[]tables.SignedIdentifier) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		output := map[string]interface{}{
			"id": v.Id,
			"access_policy": []interface{}{
				map[string]interface{}{
					"start":       v.AccessPolicy.Start,
					"expiry":      v.AccessPolicy.Expiry,
					"permissions": v.AccessPolicy.Permission,
				},
			},
		}

		result = append(result, output)
	}

	return result
}

// convertTimeFormat converts the ISO8601 time format from "2006-01-02T15:04:05Z" to "2006-01-02T15:04:05.0000000Z".
// The storage table data plane API accepts multiple formats, but always return "2006-01-02T15:04:05.0000000Z".
// The storage table mgmt plane API accepts multiple formats, but always return "2006-01-02T15:04:05Z".
func convertTimeFormat(input string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05Z", input)
	if err != nil {
		return "", fmt.Errorf("parsing time %q: %v", input, err)
	}
	return t.Format("2006-01-02T15:04:05.0000000Z"), nil
}
