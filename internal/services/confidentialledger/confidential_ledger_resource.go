package resource

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/sdk/2021-05-13-preview/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceConfidentialLedger() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceConfidentialLedgerCreateUpdate,
		// Read:   resourceResourceGroupRead,
		// Update: resourceResourceGroupCreateUpdate,
		// Delete: resourceResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ResourceGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceConfidentialLedgerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing resource group: %+v", err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_resource_group", *existing.ID)
		}
	}

	parameters := confidentialledger.ConfidentialLedger{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, name, parameters); err != nil {
		return fmt.Errorf("creating Resource Group %q: %+v", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving Resource Group %q: %+v", name, err)
	}

	// @tombuildsstuff: intentionally leaving this for now, since this'll need
	// details in the upgrade notes given how the Resource Group ID is cased incorrectly
	// but needs to be fixed (resourcegroups -> resourceGroups)
	d.SetId(*resp.ID)

	return resourceResourceGroupRead(d, meta)
}
