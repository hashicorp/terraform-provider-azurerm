// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_filesystem

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
)

func resourceStorageDataLakeGen2FileSystemCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	filesystemName := d.Get("name").(string)

	// Parse the storage_account_id which is a resource manager ID
	accountResourceManagerId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// Confirm the storage account exists and retrieve its properties
	account, err := storageClient.FindAccount(ctx, subscriptionId, accountResourceManagerId.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Data Lake Gen2 Filesystem %q: %v", accountResourceManagerId.StorageAccountName, filesystemName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", accountResourceManagerId.StorageAccountName)
	}

	// Build the data plane clients
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	// Determine the dfs endpoint, so we can build a data plane ID
	endpoint, err := account.DataPlaneEndpoint(client.EndpointTypeDfs)
	if err != nil {
		return fmt.Errorf("determining Data Lake Gen2 Filesystems endpoint: %v", err)
	}

	// Parse the dfs endpoint as a data plane account ID
	accountId, err := accounts.ParseAccountID(*endpoint, storageClient.StorageDomainSuffix)
	if err != nil {
		return fmt.Errorf("parsing Account ID: %v", err)
	}

	// Finally, build the data plane ID for this filesystem
	id := filesystems.NewFileSystemID(*accountId, filesystemName)

	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := helpers.ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %v", err)
	}

	if acl != nil && !account.IsHnsEnabled {
		return fmt.Errorf("an ACL can only be configured when Hierarchical Namespace (HNS) is enabled on the Storage Account")
	}

	propertiesRaw := d.Get("properties").(map[string]any)
	properties := helpers.ExpandMetaData(propertiesRaw)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := dataPlaneFilesystemsClient.GetProperties(ctx, id.FileSystemName)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existence of existing File System %q in %s: %v", id.FileSystemName, accountId, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_filesystem", id.ID())
		}
	}

	input := filesystems.CreateInput{
		Properties: properties,
	}
	if encryptionScope := d.Get("default_encryption_scope"); encryptionScope.(string) != "" {
		input.DefaultEncryptionScope = encryptionScope.(string)
	}

	if _, err = dataPlaneFilesystemsClient.Create(ctx, id.FileSystemName, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}

	var owner *string
	if v, ok := d.GetOk("owner"); ok {
		owner = pointer.To(v.(string))
	}
	var group *string
	if v, ok := d.GetOk("group"); ok {
		group = pointer.To(v.(string))
	}

	if acl != nil || owner != nil || group != nil {
		var aclString *string
		if acl != nil {
			v := acl.String()
			aclString = &v
		}
		accessControlInput := paths.SetAccessControlInput{
			ACL:   aclString,
			Owner: owner,
			Group: group,
		}
		if _, err = dataPlanePathsClient.SetAccessControl(ctx, id.FileSystemName, "/", accessControlInput); err != nil {
			return fmt.Errorf("setting access control for root path in File System %q in %s: %v", id.FileSystemName, accountId, err)
		}
	}

	d.SetId(id.ID())

	return resourceStorageDataLakeGen2FileSystemRead(d, meta)
}
