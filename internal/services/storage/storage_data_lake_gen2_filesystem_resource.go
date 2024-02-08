// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/paths"
	"github.com/tombuildsstuff/giovanni/storage/accesscontrol"
)

func resourceStorageDataLakeGen2FileSystem() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen2FileSystemCreate,
		Read:   resourceStorageDataLakeGen2FileSystemRead,
		Update: resourceStorageDataLakeGen2FileSystemUpdate,
		Delete: resourceStorageDataLakeGen2FileSystemDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := filesystems.ParseResourceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			storageClients := meta.(*clients.Client).Storage
			id, err := filesystems.ParseResourceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("parsing ID %q for import of Data Lake Gen2 File System: %v", d.Id(), err)
			}

			// we then need to look up the Storage Account ID
			account, err := storageClients.FindAccount(ctx, id.AccountName)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("retrieving Account %q for Data Lake Gen2 File System %q: %s", id.AccountName, id.DirectoryName, err)
			}
			if account == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
			}

			d.Set("storage_account_id", account.ID)

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
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.FileSystemsClient
	pathClient := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageID, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %s", err)
	}

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroupName, storageID.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("%s was not found", storageID)
		}

		return fmt.Errorf("checking for existence of %s: %+v", storageID, err)
	}

	if acl != nil && (storageAccount.AccountProperties == nil ||
		storageAccount.AccountProperties.IsHnsEnabled == nil ||
		!*storageAccount.AccountProperties.IsHnsEnabled) {
		return fmt.Errorf("ACL is enabled only when the Hierarchical Namespace (HNS) feature is turned ON")
	}

	fileSystemName := d.Get("name").(string)
	propertiesRaw := d.Get("properties").(map[string]interface{})
	properties := ExpandMetaData(propertiesRaw)

	id := client.GetResourceID(storageID.StorageAccountName, fileSystemName)

	resp, err := client.GetProperties(ctx, storageID.StorageAccountName, fileSystemName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existence of existing File System %q in %s: %+v", fileSystemName, storageID, err)
		}
	}

	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_filesystem", id)
	}

	log.Printf("[INFO] Creating File System %q in %s.", fileSystemName, storageID)
	input := filesystems.CreateInput{
		Properties: properties,
	}
	if _, err := client.Create(ctx, storageID.StorageAccountName, fileSystemName, input); err != nil {
		return fmt.Errorf("creating File System %q in %s: %+v", fileSystemName, storageID, err)
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
			log.Printf("[INFO] Creating acl %q in File System %q in %s", acl, fileSystemName, storageID)
			v := acl.String()
			aclString = &v
		}
		accessControlInput := paths.SetAccessControlInput{
			ACL:   aclString,
			Owner: owner,
			Group: group,
		}
		if _, err := pathClient.SetAccessControl(ctx, storageID.StorageAccountName, fileSystemName, "/", accessControlInput); err != nil {
			return fmt.Errorf("setting access control for root path in File System %q in %s: %+v", fileSystemName, storageID, err)
		}
	}

	d.SetId(id)
	return resourceStorageDataLakeGen2FileSystemRead(d, meta)
}

func resourceStorageDataLakeGen2FileSystemUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.FileSystemsClient
	pathClient := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	storageId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("parsing ace list: %s", err)
	}

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageId.ResourceGroupName, storageId.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("%s was not found", storageId)
		}

		return fmt.Errorf("checking for existence of %s: %+v", storageId, err)
	}

	if acl != nil && (storageAccount.AccountProperties == nil ||
		storageAccount.AccountProperties.IsHnsEnabled == nil ||
		!*storageAccount.AccountProperties.IsHnsEnabled) {
		return fmt.Errorf("ACL is enabled only when the Hierarchical Namespace (HNS) feature is turned ON")
	}

	propertiesRaw := d.Get("properties").(map[string]interface{})
	properties := ExpandMetaData(propertiesRaw)

	log.Printf("[INFO] Updating Properties for File System %q in Storage Account %q.", id.DirectoryName, id.AccountName)
	input := filesystems.SetPropertiesInput{
		Properties: properties,
	}
	if _, err = client.SetProperties(ctx, id.AccountName, id.DirectoryName, input); err != nil {
		return fmt.Errorf("updating Properties for File System %q in Storage Account %q: %s", id.DirectoryName, id.AccountName, err)
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
			log.Printf("[INFO] Creating acl %q in File System %q in Storage Account %q.", acl, id.DirectoryName, id.AccountName)
			v := acl.String()
			aclString = &v
		}
		accessControlInput := paths.SetAccessControlInput{
			ACL:   aclString,
			Owner: owner,
			Group: group,
		}
		if _, err := pathClient.SetAccessControl(ctx, id.AccountName, id.DirectoryName, "/", accessControlInput); err != nil {
			return fmt.Errorf("setting access control for root path in File System %q in Storage Account %q: %s", id.DirectoryName, id.AccountName, err)
		}
	}

	return resourceStorageDataLakeGen2FileSystemRead(d, meta)
}

func resourceStorageDataLakeGen2FileSystemRead(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.FileSystemsClient
	pathClient := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	storageID, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroupName, storageID.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			log.Printf("[INFO] Storage Account %q does not exist removing from state...", id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("checking for existence of %s for File System %q: %+v", storageID, id.DirectoryName, err)
	}

	resp, err := client.GetProperties(ctx, id.AccountName, id.DirectoryName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] File System %q does not exist in Storage Account %q - removing from state...", id.DirectoryName, id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving File System %q in Storage Account %q: %+v", id.DirectoryName, id.AccountName, err)
	}

	d.Set("name", id.DirectoryName)

	if err := d.Set("properties", resp.Properties); err != nil {
		return fmt.Errorf("setting `properties`: %+v", err)
	}

	var ace []interface{}
	var owner, group string
	// acl is only enabled when `IsHnsEnabled` is true otherwise the rest api will report error
	if storageAccount.AccountProperties != nil && storageAccount.AccountProperties.IsHnsEnabled != nil &&
		*storageAccount.AccountProperties.IsHnsEnabled {
		// The above `getStatus` API request doesn't return the ACLs
		// Have to make a `getAccessControl` request, but that doesn't return all fields either!
		pathResponse, err := pathClient.GetProperties(ctx, id.AccountName, id.DirectoryName, "/", paths.GetPropertiesActionGetAccessControl)
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
	client := meta.(*clients.Client).Storage.FileSystemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := filesystems.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.AccountName, id.DirectoryName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting File System %q in Storage Account %q: %+v", id.DirectoryName, id.AccountName, err)
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
