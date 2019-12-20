package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/blob/containers"
)

func resourceArmStorageContainer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceArmStorageContainerCreate,
		Read:          resourceArmStorageContainerRead,
		Delete:        resourceArmStorageContainerDelete,
		Update:        resourceArmStorageContainerUpdate,
		MigrateState:  resourceStorageContainerMigrateState,
		SchemaVersion: 1,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageContainerName,
			},

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArmStorageAccountName,
			},

			"container_access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					string(containers.Blob),
					string(containers.Container),
					"private",
				}, false),
			},

			"metadata": storage.MetaDataComputedSchema(),

			// TODO: support for ACL's, Legal Holds and Immutability Policies
			"has_immutability_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"has_legal_hold": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameDeprecated(),

			"properties": {
				Type:       schema.TypeMap,
				Computed:   true,
				Deprecated: "This field will be removed in version 2.0 of the Azure Provider",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmStorageContainerCreate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	containerName := d.Get("name").(string)
	accountName := d.Get("storage_account_name").(string)
	accessLevelRaw := d.Get("container_access_type").(string)
	accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	metaData := storage.ExpandMetaData(metaDataRaw)

	account, err := storageClient.FindAccount(ctx, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", accountName, containerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", accountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client: %s", err)
	}

	id := client.GetResourceID(accountName, containerName)
	if features.ShouldResourcesBeImported() {
		existing, err := client.GetProperties(ctx, accountName, containerName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for existence of existing Container %q (Account %q / Resource Group %q): %+v", containerName, accountName, account.ResourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_storage_container", id)
		}
	}

	log.Printf("[INFO] Creating Container %q in Storage Account %q", containerName, accountName)
	input := containers.CreateInput{
		AccessLevel: accessLevel,
		MetaData:    metaData,
	}
	if _, err := client.Create(ctx, accountName, containerName, input); err != nil {
		return fmt.Errorf("Error creating Container %q (Account %q / Resource Group %q): %s", containerName, accountName, account.ResourceGroup, err)
	}

	d.SetId(id)
	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if d.HasChange("container_access_type") {
		log.Printf("[DEBUG] Updating the Access Control for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		accessLevelRaw := d.Get("container_access_type").(string)
		accessLevel := expandStorageContainerAccessLevel(accessLevelRaw)

		if _, err := client.SetAccessControl(ctx, id.AccountName, id.ContainerName, accessLevel); err != nil {
			return fmt.Errorf("Error updating the Access Control for Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the Access Control for Container %q (Storage Account %q / Resource Group %q)", id.ContainerName, id.AccountName, account.ResourceGroup)
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for Container %q (Storage Account %q / Resource Group %q)..", id.ContainerName, id.AccountName, account.ResourceGroup)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		metaData := storage.ExpandMetaData(metaDataRaw)

		if _, err := client.SetMetaData(ctx, id.AccountName, id.ContainerName, metaData); err != nil {
			return fmt.Errorf("Error updating the MetaData for Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the MetaData for Container %q (Storage Account %q / Resource Group %q)", id.ContainerName, id.AccountName, account.ResourceGroup)
	}

	return resourceArmStorageContainerRead(d, meta)
}

func resourceArmStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		log.Printf("[DEBUG] Unable to locate Account %q for Storage Container %q - assuming removed & removing from state", id.AccountName, id.ContainerName)
		d.SetId("")
		return nil
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	props, err := client.GetProperties(ctx, id.AccountName, id.ContainerName)
	if err != nil {
		if utils.ResponseWasNotFound(props.Response) {
			log.Printf("[DEBUG] Container %q was not found in Account %q / Resource Group %q - assuming removed & removing from state", id.ContainerName, id.AccountName, account.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Container %q (Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	d.Set("name", id.ContainerName)
	d.Set("storage_account_name", id.AccountName)
	d.Set("resource_group_name", account.ResourceGroup)

	d.Set("container_access_type", flattenStorageContainerAccessLevel(props.AccessLevel))

	if err := d.Set("metadata", storage.FlattenMetaData(props.MetaData)); err != nil {
		return fmt.Errorf("Error setting `metadata`: %+v", err)
	}

	if err := d.Set("properties", flattenStorageContainerProperties(props)); err != nil {
		return fmt.Errorf("Error setting `properties`: %+v", err)
	}

	d.Set("has_immutability_policy", props.HasImmutabilityPolicy)
	d.Set("has_legal_hold", props.HasLegalHold)

	return nil
}

func resourceArmStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := containers.ParseResourceID(d.Id())
	if err != nil {
		return err
	}

	account, err := storageClient.FindAccount(ctx, id.AccountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Account %q for Container %q: %s", id.AccountName, id.ContainerName, err)
	}
	if account == nil {
		return fmt.Errorf("Unable to locate Storage Account %q!", id.AccountName)
	}

	client, err := storageClient.ContainersClient(ctx, *account)
	if err != nil {
		return fmt.Errorf("Error building Containers Client for Storage Account %q (Resource Group %q): %s", id.AccountName, account.ResourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.AccountName, id.ContainerName); err != nil {
		return fmt.Errorf("Error deleting Container %q (Storage Account %q / Resource Group %q): %s", id.ContainerName, id.AccountName, account.ResourceGroup, err)
	}

	return nil
}

func flattenStorageContainerProperties(input containers.ContainerProperties) map[string]interface{} {
	output := map[string]interface{}{
		"last_modified":  input.Header.Get("Last-Modified"),
		"lease_duration": "",
		"lease_state":    string(input.LeaseState),
		"lease_status":   string(input.LeaseStatus),
	}

	if input.LeaseDuration != nil {
		output["lease_duration"] = string(*input.LeaseDuration)
	}

	return output
}

func expandStorageContainerAccessLevel(input string) containers.AccessLevel {
	// for historical reasons, "private" above is an empty string in the API
	// so the enum doesn't 1:1 match. You could argue the SDK should handle this
	// but this is suitable for now
	if input == "private" {
		return containers.Private
	}

	return containers.AccessLevel(input)
}

func flattenStorageContainerAccessLevel(input containers.AccessLevel) string {
	// for historical reasons, "private" above is an empty string in the API
	if input == containers.Private {
		return "private"
	}

	return string(input)
}
