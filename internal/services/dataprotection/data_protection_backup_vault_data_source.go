package dataprotection

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/legacysdk/dataprotection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDataProtectionBackupVault() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataProtectionBackupVaultRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.BackupVaultID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{2,50}$"),
					"DataProtection BackupVault name must be 2 - 50 characters long, contain only letters, numbers and hyphens.).",
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"datastore_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"redundancy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceDataProtectionBackupVaultRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataProtection.BackupVaultClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewBackupVaultID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id.Name, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataProtection BackupVault %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataProtection BackupVault (%q): %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if props.StorageSettings != nil && len(*props.StorageSettings) > 0 {
			d.Set("datastore_type", (*props.StorageSettings)[0].DatastoreType)
			d.Set("redundancy", (*props.StorageSettings)[0].Type)
		}
	}
	if err := d.Set("identity", dataSourceFlattenBackupVaultDppIdentityDetails(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func dataSourceFlattenBackupVaultDppIdentityDetails(input *dataprotection.DppIdentityDetails) []interface{} {
	var config *identity.SystemAssigned
	if input != nil {
		principalId := ""
		if input.PrincipalID != nil {
			principalId = *input.PrincipalID
		}

		tenantId := ""
		if input.TenantID != nil {
			tenantId = *input.TenantID
		}
		config = &identity.SystemAssigned{
			Type:        identity.Type(*input.Type),
			PrincipalId: principalId,
			TenantId:    tenantId,
		}
	}
	return identity.FlattenSystemAssigned(config)
}
