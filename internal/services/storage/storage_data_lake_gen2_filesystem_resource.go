// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/filesystems"
	"github.com/jackofallops/giovanni/storage/2023-11-03/datalakestore/paths"
	"github.com/jackofallops/giovanni/storage/accesscontrol"
)

func resourceStorageDataLakeGen2FileSystem() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen2FileSystemCreate,
		Read:   resourceStorageDataLakeGen2FileSystemRead,
		Update: resourceStorageDataLakeGen2FileSystemUpdate,
		Delete: resourceStorageDataLakeGen2FileSystemDelete,

		Importer: helpers.ImporterValidatingStorageResourceIdThen(func(id, storageDomainSuffix string) error {
			_, err := filesystems.ParseFileSystemID(id, storageDomainSuffix)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			storageClient := meta.(*clients.Client).Storage
			subscriptionId := meta.(*clients.Client).Account.SubscriptionId
			ctx, cancel := context.WithTimeout(ctx, d.Timeout(pluginsdk.TimeoutRead))
			defer cancel()

			id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing ID %q for import of Data Lake Gen2 File System: %v", d.Id(), err)
			}

			// we then need to look up the Storage Account ID
			account, err := storageClient.FindAccount(ctx, subscriptionId, id.AccountId.AccountName)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving Account %q for Data Lake Gen2 File System %q: %s", id.AccountId.AccountName, id.FileSystemName, err)
			}
			if account == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("unable to locate Storage Account: %q", id.AccountId.AccountName)
			}

			d.Set("storage_account_id", account.StorageAccountId.ID())

			return []*pluginsdk.ResourceData{d}, nil
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
				ValidateFunc: validateStorageDataLakeGen2FileSystemName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"properties": MetaDataSchema(),

			"default_encryption_scope": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true, // needed because a dummy value is returned when unspecified
				ForceNew:     true,
				ValidateFunc: validate.StorageEncryptionScopeName,
			},

			"owner": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"$superuser"}, false)),
			},

			"group": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"$superuser"}, false)),
			},

			"ace": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"scope": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"default", "access"}, false),
							Default:      "access",
						},
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"user", "group", "mask", "other"}, false),
						},
						"id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"permissions": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ADLSAccessControlPermissions,
						},
					},
				},
			},
		},
	}
}

func resourceStorageDataLakeGen2FileSystemCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %v", err)
	}

	if acl != nil && !account.IsHnsEnabled {
		return fmt.Errorf("an ACL can only be configured when Hierarchical Namespace (HNS) is enabled on the Storage Account")
	}

	propertiesRaw := d.Get("properties").(map[string]interface{})
	properties := ExpandMetaData(propertiesRaw)

	resp, err := dataPlaneFilesystemsClient.GetProperties(ctx, id.FileSystemName)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for existence of existing File System %q in %s: %v", id.FileSystemName, accountId, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_filesystem", id.ID())
	}

	log.Printf("[INFO] Creating %s...", id)
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
		sv := v.(string)
		owner = &sv
	}
	var group *string
	if v, ok := d.GetOk("group"); ok {
		sv := v.(string)
		group = &sv
	}

	if acl != nil || owner != nil || group != nil {
		var aclString *string
		if acl != nil {
			log.Printf("[INFO] Creating ACL %q for %s", acl, id)
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

func resourceStorageDataLakeGen2FileSystemUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
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

	// Build the data plane clients
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %v", err)
	}

	if acl != nil && !account.IsHnsEnabled {
		return fmt.Errorf("an ACL can only be configured when Hierarchical Namespace (HNS) is enabled on the Storage Account")
	}

	propertiesRaw := d.Get("properties").(map[string]interface{})
	properties := ExpandMetaData(propertiesRaw)

	log.Printf("[INFO] Updating Properties for %s...", id)
	input := filesystems.SetPropertiesInput{
		Properties: properties,
	}
	if _, err = dataPlaneFilesystemsClient.SetProperties(ctx, id.FileSystemName, input); err != nil {
		return fmt.Errorf("updating Properties for %s: %v", id, err)
	}

	var owner *string
	if v, ok := d.GetOk("owner"); ok {
		sv := v.(string)
		owner = &sv
	}
	var group *string
	if v, ok := d.GetOk("group"); ok {
		sv := v.(string)
		group = &sv
	}

	if acl != nil || owner != nil || group != nil {
		var aclString *string
		if acl != nil {
			log.Printf("[INFO] Creating ACL %q for %s...", acl, id)
			v := acl.String()
			aclString = &v
		}
		accessControlInput := paths.SetAccessControlInput{
			ACL:   aclString,
			Owner: owner,
			Group: group,
		}
		if _, err = dataPlanePathsClient.SetAccessControl(ctx, id.FileSystemName, "/", accessControlInput); err != nil {
			return fmt.Errorf("setting access control for root path in File System %q in Storage Account %q: %v", id.FileSystemName, id.AccountId.AccountName, err)
		}
	}

	return resourceStorageDataLakeGen2FileSystemRead(d, meta)
}

func resourceStorageDataLakeGen2FileSystemRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
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

	// Build the data plane clients
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}
	dataPlanePathsClient, err := storageClient.DataLakePathsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Paths Client: %v", err)
	}

	resp, err := dataPlaneFilesystemsClient.GetProperties(ctx, id.FileSystemName)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] File System %q does not exist in Storage Account %q - removing from state...", id.FileSystemName, id.AccountId.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %v", id, err)
	}

	d.Set("name", id.FileSystemName)
	d.Set("default_encryption_scope", resp.DefaultEncryptionScope)

	if err = d.Set("properties", resp.Properties); err != nil {
		return fmt.Errorf("setting `properties`: %v", err)
	}

	var ace []interface{}
	var owner, group string
	// acl is only enabled when `IsHnsEnabled` is true otherwise the rest api will report error
	if account.IsHnsEnabled {
		// The above `getStatus` API request doesn't return the ACLs
		// Have to make a `getAccessControl` request, but that doesn't return all fields either!
		payload := paths.GetPropertiesInput{
			Action: paths.GetPropertiesActionGetAccessControl,
		}
		pathResponse, err := dataPlanePathsClient.GetProperties(ctx, id.FileSystemName, "/", payload)
		if err == nil {
			acl, err := accesscontrol.ParseACL(pathResponse.ACL)
			if err != nil {
				return fmt.Errorf("parsing response ACL %q: %s", pathResponse.ACL, err)
			}
			ace = FlattenDataLakeGen2AceList(d, acl)
			owner = pathResponse.Owner
			group = pathResponse.Group
		}
	}
	d.Set("ace", ace)
	d.Set("owner", owner)
	d.Set("group", group)

	return nil
}

func resourceStorageDataLakeGen2FileSystemDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseFileSystemID(d.Id(), storageClient.StorageDomainSuffix)
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
	dataPlaneFilesystemsClient, err := storageClient.DataLakeFilesystemsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Data Lake Gen2 Filesystems Client: %v", err)
	}

	resp, err := dataPlaneFilesystemsClient.Delete(ctx, id.FileSystemName)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}

func validateStorageDataLakeGen2FileSystemName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^\$root$|^[0-9a-z-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"only lowercase alphanumeric characters and hyphens allowed in %q: %q",
			k, value))
	}
	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 3 and 63 characters: %q", k, value))
	}
	if regexp.MustCompile(`^-`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q cannot begin with a hyphen: %q", k, value))
	}
	return warnings, errors
}
