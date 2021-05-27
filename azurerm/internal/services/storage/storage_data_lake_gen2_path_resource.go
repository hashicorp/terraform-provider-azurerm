package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2019-12-12/datalakestore/paths"
	"github.com/tombuildsstuff/giovanni/storage/accesscontrol"
)

func resourceStorageDataLakeGen2Path() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageDataLakeGen2PathCreate,
		Read:   resourceStorageDataLakeGen2PathRead,
		Update: resourceStorageDataLakeGen2PathUpdate,
		Delete: resourceStorageDataLakeGen2PathDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := paths.ParseResourceID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			storageClients := meta.(*clients.Client).Storage

			id, err := paths.ParseResourceID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("Error parsing ID %q for import of Data Lake Gen2 Path: %v", d.Id(), err)
			}

			// we then need to look up the Storage Account ID
			account, err := storageClients.FindAccount(ctx, id.AccountName)
			if err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("Error retrieving Account %q for Data Lake Gen2 Path %q in File System %q: %s", id.AccountName, id.Path, id.FileSystemName, err)
			}
			if account == nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
			}

			if _, err = storageClients.FileSystemsClient.GetProperties(ctx, id.AccountName, id.FileSystemName); err != nil {
				return []*pluginsdk.ResourceData{d}, fmt.Errorf("Error retrieving File System %q for Data Lake Gen 2 Path %q in Account %q: %s", id.FileSystemName, id.Path, id.AccountName, err)
			}

			d.Set("storage_account_id", account.ID)
			d.Set("filesystem_name", id.FileSystemName)

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"filesystem_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStorageDataLakeGen2FileSystemName,
			},

			"path": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true, // TODO handle rename
			},

			"resource": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"directory"}, false),
			},

			"owner": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},

			"group": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
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

func resourceStorageDataLakeGen2PathCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageID, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroup, storageID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q was not found in Resource Group %q!", storageID.Name, storageID.ResourceGroup)
		}

		return fmt.Errorf("Error checking for existence of Storage Account %q (Resource Group %q): %+v", storageID.Name, storageID.ResourceGroup, err)
	}

	fileSystemName := d.Get("filesystem_name").(string)
	path := d.Get("path").(string)

	id := client.GetResourceID(storageID.Name, fileSystemName, path)
	resp, err := client.GetProperties(ctx, storageID.Name, fileSystemName, path, paths.GetPropertiesActionGetStatus)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for existence of existing Path %q in  File System %q (Account %q): %+v", path, fileSystemName, storageID.Name, err)
		}
	}
	if !utils.ResponseWasNotFound(resp.Response) {
		return tf.ImportAsExistsError("azurerm_storage_data_lake_gen2_path", id)
	}

	resourceString := d.Get("resource").(string)
	var resource paths.PathResource
	switch resourceString {
	case "directory":
		resource = paths.PathResourceDirectory
	default:
		return fmt.Errorf("Unhandled resource type %q", resourceString)
	}
	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("Error parsing ace list: %s", err)
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

	log.Printf("[INFO] Creating Path %q in File System %q in Storage Account %q.", path, fileSystemName, storageID.Name)
	input := paths.CreateInput{
		Resource: resource,
	}

	if _, err := client.Create(ctx, storageID.Name, fileSystemName, path, input); err != nil {
		return fmt.Errorf("Error creating Path %q in File System %q in Storage Account %q: %s", path, fileSystemName, storageID.Name, err)
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
		if _, err := client.SetAccessControl(ctx, storageID.Name, fileSystemName, path, accessControlInput); err != nil {
			return fmt.Errorf("Error setting access control for Path %q in File System %q in Storage Account %q: %s", path, fileSystemName, storageID.Name, err)
		}
	}

	d.SetId(id)
	return resourceStorageDataLakeGen2PathRead(d, meta)
}

func resourceStorageDataLakeGen2PathUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	accountsClient := meta.(*clients.Client).Storage.AccountsClient
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := paths.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	storageID, err := parse.StorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	path := d.Get("path").(string)

	aceRaw := d.Get("ace").(*pluginsdk.Set).List()
	acl, err := ExpandDataLakeGen2AceList(aceRaw)
	if err != nil {
		return fmt.Errorf("Error parsing ace list: %s", err)
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

	// confirm the storage account exists, otherwise Data Plane API requests will fail
	storageAccount, err := accountsClient.GetProperties(ctx, storageID.ResourceGroup, storageID.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q was not found in Resource Group %q!", storageID.Name, storageID.ResourceGroup)
		}

		return fmt.Errorf("Error checking for existence of Storage Account %q (Resource Group %q): %+v", storageID.Name, storageID.ResourceGroup, err)
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
		if _, err := client.SetAccessControl(ctx, id.AccountName, id.FileSystemName, path, accessControlInput); err != nil {
			return fmt.Errorf("Error setting access control for Path %q in File System %q in Storage Account %q: %s", path, id.FileSystemName, id.AccountName, err)
		}
	}

	return resourceStorageDataLakeGen2PathRead(d, meta)
}

func resourceStorageDataLakeGen2PathRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := paths.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, id.AccountName, id.FileSystemName, id.Path, paths.GetPropertiesActionGetStatus)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Path %q does not exist in File System %q in Storage Account %q - removing from state...", id.Path, id.FileSystemName, id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Path %q in File System %q in Storage Account %q: %+v", id.Path, id.FileSystemName, id.AccountName, err)
	}

	d.Set("path", id.Path)
	d.Set("resource", resp.ResourceType)
	d.Set("owner", resp.Owner)
	d.Set("group", resp.Group)

	// The above `getStatus` API request doesn't return the ACLs
	// Have to make a `getAccessControl` request, but that doesn't return all fields either!
	resp, err = client.GetProperties(ctx, id.AccountName, id.FileSystemName, id.Path, paths.GetPropertiesActionGetAccessControl)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Path %q does not exist in File System %q in Storage Account %q - removing from state...", id.Path, id.FileSystemName, id.AccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ACLs for Path %q in File System %q in Storage Account %q: %+v", id.Path, id.FileSystemName, id.AccountName, err)
	}

	acl, err := accesscontrol.ParseACL(resp.ACL)
	if err != nil {
		return fmt.Errorf("Error parsing response ACL %q: %s", resp.ACL, err)
	}
	d.Set("ace", FlattenDataLakeGen2AceList(acl))

	return nil
}

func resourceStorageDataLakeGen2PathDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ADLSGen2PathsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := paths.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.AccountName, id.FileSystemName, id.Path)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Path %q in File System %q in Storage Account %q: %+v", id.Path, id.FileSystemName, id.AccountName, err)
		}
	}

	return nil
}
