package storage

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storagesync/mgmt/2019-06-01/storagesync"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStorageSync() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStorageSyncCreate,
		Read:   resourceArmStorageSyncRead,
		Update: resourceArmStorageSyncUpdate,
		Delete: resourceArmStorageSyncDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parsers.ParseStorageSyncID(id)
			return err
		}),

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
				ValidateFunc: ValidateArmStorageSyncName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmStorageSyncCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StoragesyncClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Storage Sync(Storage Sync Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_storage_sync", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	parameters := storagesync.ServiceCreateParameters{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.Create(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Storage Sync(Storage Sync Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Storage Sync(Storage Sync Name %q / Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Storage Sync(Storage Sync Name %q / Resource Group %q) ID", name, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmStorageSyncRead(d, meta)
}

func resourceArmStorageSyncRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StoragesyncClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.ParseStorageSyncID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Storage Sync %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Storage Sync(Storage Sync Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	var location string
	if resp.Location != nil {
		location = *resp.Location
	}
	d.Set("location", azure.NormalizeLocation(location))
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmStorageSyncUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StoragesyncClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.ParseStorageSyncID(d.Id())
	if err != nil {
		return err
	}

	update := storagesync.ServiceUpdateParameters{}

	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	_, err = client.Update(ctx, id.ResourceGroup, id.Name, &update)
	if err != nil {
		return fmt.Errorf("Error updating Storage Sync %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmStorageSyncRead(d, meta)
}

func resourceArmStorageSyncDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.StoragesyncClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parsers.ParseStorageSyncID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Error deleting Storage Sync(Storage Sync Name %q / Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func ValidateArmStorageSyncName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[0-9a-zA-Z-_.]*[0-9a-zA-Z-_]$").MatchString(strings.TrimSpace(input)) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of letters, numbers, spaces, and any of the following characters: '.-_' and that does not end with characters: '. '", input))
	}

	return warnings, errors
}
