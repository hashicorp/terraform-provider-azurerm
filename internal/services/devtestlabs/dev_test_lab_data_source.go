package devtestlabs

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceDevTestLab() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDevTestLabRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"storage_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

			"artifacts_storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_premium_storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"key_vault_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"premium_data_disk_storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"unique_identifier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDevTestLabRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDevTestLabID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	read, err := client.Get(ctx, id.ResourceGroup, id.LabName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", read.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := read.LabProperties; props != nil {
		d.Set("storage_type", string(props.LabStorageType))

		// Computed fields
		d.Set("artifacts_storage_account_id", props.ArtifactsStorageAccount)
		d.Set("default_storage_account_id", props.DefaultStorageAccount)
		d.Set("default_premium_storage_account_id", props.DefaultPremiumStorageAccount)
		d.Set("key_vault_id", props.VaultName)
		d.Set("premium_data_disk_storage_account_id", props.PremiumDataDiskStorageAccount)
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	return tags.FlattenAndSet(d, read.Tags)
}
