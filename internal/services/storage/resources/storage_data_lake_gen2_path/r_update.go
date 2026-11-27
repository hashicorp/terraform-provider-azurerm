// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_data_lake_gen2_path

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
)

func resourceStorageDataLakeGen2PathUpdate(d *pluginsdk.ResourceData, meta any) error {
	storageClient := client.From(meta)
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForUpdate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := paths.ParsePathID(d.Id(), storageClient.StorageDomainSuffix)
	if err != nil {
		return err
	}

	// Retrieve the storage account properties
	account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
	if err != nil {
		return fmt.Errorf("retrieving Account %q for Data Lake Gen2 Filesystem %q: %v", id.AccountId.AccountName, id.FileSystemName, err)
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", id.AccountId.AccountName)
	}

	// Build the data plane client
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	path := d.Get("path").(string)

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
		if _, err = dataPlanePathsClient.SetAccessControl(ctx, id.FileSystemName, path, accessControlInput); err != nil {
			return fmt.Errorf("setting access control for %s: %s", id, err)
		}
	}

	return resourceStorageDataLakeGen2PathRead(d, meta)
}
