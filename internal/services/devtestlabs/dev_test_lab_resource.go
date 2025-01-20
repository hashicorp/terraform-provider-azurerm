// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devtestlabs

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDevTestLab() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceDevTestLabCreateUpdate,
		Read:   resourceDevTestLabRead,
		Update: resourceDevTestLabCreateUpdate,
		Delete: resourceDevTestLabDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := labs.ParseLabID(id)
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

			"location": commonschema.Location(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

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

	return resource
}

func resourceDevTestLabCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Lab creation")

	id := labs.NewLabID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, labs.GetOperationOptions{})
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dev_test_lab", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	parameters := labs.Lab{
		Location: utils.String(location),
		Tags:     expandTags(d.Get("tags").(map[string]interface{})),
	}

	err := client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDevTestLabRead(d, meta)
}

func resourceDevTestLabRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := labs.ParseLabID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, labs.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		if props := model.Properties; props != nil {
			// Computed fields
			d.Set("artifacts_storage_account_id", props.ArtifactsStorageAccount)
			d.Set("default_storage_account_id", props.DefaultStorageAccount)
			d.Set("default_premium_storage_account_id", props.DefaultPremiumStorageAccount)

			kvId := ""
			if props.VaultName != nil {
				id, err := commonids.ParseKeyVaultIDInsensitively(*props.VaultName)
				if err != nil {
					return fmt.Errorf("parsing %q: %+v", *props.VaultName, err)
				}
				kvId = id.ID()
			}
			d.Set("key_vault_id", kvId)
			d.Set("premium_data_disk_storage_account_id", props.PremiumDataDiskStorageAccount)
			d.Set("unique_identifier", props.UniqueIdentifier)
		}
		if err = tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceDevTestLabDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := labs.ParseLabID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id, labs.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}
