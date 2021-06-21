package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/identity"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dataprotection/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataProtectionBackupVaultCreate,
		Read:   resourceDataProtectionBackupVaultRead,
		Update: resourceDataProtectionBackupVaultUpdate,
		Delete: resourceDataProtectionBackupVaultDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BackupVaultID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dataprotection.StorageSettingStoreTypesArchiveStore),
					string(dataprotection.StorageSettingStoreTypesSnapshotStore),
					string(dataprotection.StorageSettingStoreTypesVaultStore),
				}, false),
			},

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(dataprotection.GeoRedundant),
					string(dataprotection.LocallyRedundant),
				}, false),
			},

			"identity": identity.SystemAssigned{}.Schema(),

			"tags": tags.Schema(),
		},
	}
}
func resourceDataProtectionBackupVaultCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewBackupVaultID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing DataProtection BackupVault (%q): %+v", id, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_data_protection_backup_vault", id.ID())
	}

	parameters := dataprotection.BackupVaultResource{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &dataprotection.BackupVault{
			StorageSettings: &[]dataprotection.StorageSetting{
				{
					DatastoreType: dataprotection.StorageSettingStoreTypes(d.Get("datastore_type").(string)),
					Type:          dataprotection.StorageSettingTypes(d.Get("redundancy").(string)),
				}},
		},
		Identity: expandBackupVaultDppIdentityDetails(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	future, err := client.CreateOrUpdate(ctx, id.Name, id.ResourceGroup, parameters)
	if err != nil {
		return fmt.Errorf("creating DataProtection BackupVault (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataProtectionBackupVaultRead(d, meta)
}

func resourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupVaultID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if props.StorageSettings != nil && len(*props.StorageSettings) > 0 {
			d.Set("datastore_type", (*props.StorageSettings)[0].DatastoreType)
			d.Set("redundancy", (*props.StorageSettings)[0].Type)
		}
	}
	if err := d.Set("identity", flattenBackupVaultDppIdentityDetails(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDataProtectionBackupVaultUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupVaultID(d.Id())
	if err != nil {
		return err
	}

	parameters := dataprotection.PatchResourceRequestInput{}
	if d.HasChange("identity") {
		parameters.Identity = expandBackupVaultDppIdentityDetails(d.Get("identity").([]interface{}))
	}
	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Patch(ctx, id.Name, id.ResourceGroup, parameters)
	if err != nil {
		return fmt.Errorf("updating DataProtection BackupVault (%q): %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of the DataProtection BackupVault %q: %+v", id, err)
	}
	return resourceDataProtectionBackupVaultRead(d, meta)
}

func resourceDataProtectionBackupVaultDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BackupVaultID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, id.Name, id.ResourceGroup); err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}
		return fmt.Errorf("deleting DataProtection BackupVault (%q): %+v", id, err)
	}
	return nil
}

func expandBackupVaultDppIdentityDetails(input []interface{}) *dataprotection.DppIdentityDetails {
	config, _ := identity.SystemAssigned{}.Expand(input)
	return &dataprotection.DppIdentityDetails{
		Type: utils.String(config.Type),
	}
}

func flattenBackupVaultDppIdentityDetails(input *dataprotection.DppIdentityDetails) []interface{} {
	var config *identity.ExpandedConfig
	if input != nil {
		config = &identity.ExpandedConfig{
			Type:        *input.Type,
			PrincipalId: input.PrincipalID,
			TenantId:    input.TenantID,
		}
	}
	return identity.SystemAssigned{}.Flatten(config)
}
