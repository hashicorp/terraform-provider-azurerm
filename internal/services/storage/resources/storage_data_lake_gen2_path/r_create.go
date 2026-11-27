// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_path

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
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
)

func resourceStorageDataLakeGen2PathCreate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	filesystemName := d.Get("filesystem_name").(string)
	path := d.Get("path").(string)

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

	// Build the data plane client
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

	id := paths.NewPathID(*accountId, filesystemName, path)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := dataPlanePathsClient.GetProperties(ctx, filesystemName, path, paths.GetPropertiesInput{Action: paths.GetPropertiesActionGetStatus})
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existence of existing Path %q in File System %q in %s: %v", path, filesystemName, accountResourceManagerId, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_path", id.ID())
		}
	}

	resourceString := d.Get("resource").(string)
	var resource paths.PathResource
	switch resourceString {
	case "directory":
		resource = paths.PathResourceDirectory
	default:
		return fmt.Errorf("unhandled resource type %q", resourceString)
	}
	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := helpers.ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %v", err)
	}

	var owner *string
	if v, ok := d.GetOk("owner"); ok {
		owner = pointer.To(v.(string))
	}
	var group *string
	if v, ok := d.GetOk("group"); ok {
		group = pointer.To(v.(string))
	}

	input := paths.CreateInput{
		Resource: resource,
	}

	if _, err = dataPlanePathsClient.Create(ctx, filesystemName, path, input); err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}
	d.SetId(id.ID())

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
		if _, err = dataPlanePathsClient.SetAccessControl(ctx, filesystemName, path, accessControlInput); err != nil {
			return fmt.Errorf("setting access control for %s: %+v", id, err)
		}
	}

	return resourceStorageDataLakeGen2PathRead(d, meta)
}
