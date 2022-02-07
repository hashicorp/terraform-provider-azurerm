package devtestlabs

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2018-09-15/dtl"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	keyvaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDevTestLab() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDevTestLabCreateUpdate,
		Read:   resourceDevTestLabRead,
		Update: resourceDevTestLabCreateUpdate,
		Delete: resourceDevTestLabDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DevTestLabID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestLabUpgradeV0ToV1{},
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
				ValidateFunc: validate.DevTestLabName(),
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"storage_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(dtl.Premium),
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.Standard),
					string(dtl.Premium),
				}, false),
			},

			"tags": tags.Schema(),

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

func resourceDevTestLabCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Lab creation")

	id := parse.NewDevTestLabID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LabName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_dev_test_lab", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storageType := d.Get("storage_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := dtl.Lab{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		LabProperties: &dtl.LabProperties{
			LabStorageType: dtl.StorageType(storageType),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LabName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDevTestLabRead(d, meta)
}

func resourceDevTestLabRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DevTestLabID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LabName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.LabName)
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

		kvId := ""
		if props.VaultName != nil {
			id, err := keyvaultParse.VaultID(*props.VaultName)
			if err != nil {
				return fmt.Errorf("parsing %q: %+v", *props.VaultName, err)
			}
			kvId = id.ID()
		}
		d.Set("key_vault_id", kvId)
		d.Set("premium_data_disk_storage_account_id", props.PremiumDataDiskStorageAccount)
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceDevTestLabDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DevTestLabID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LabName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.LabName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
