// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package powerbi_embedded

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	powerbiClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePowerBIEmbeddedCreate(d *pluginsdk.ResourceData, meta any) error {
	client := powerbiClient.From(meta).CapacityClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := capacities.NewCapacityID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetDetails(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_powerbi_embedded", id.ID())
		}
	}

	administrators := d.Get("administrators").(*pluginsdk.Set).List()
	mode := capacities.Mode(d.Get("mode").(string))

	parameters := capacities.DedicatedCapacity{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &capacities.DedicatedCapacityProperties{
			Administration: &capacities.DedicatedCapacityAdministrators{
				Members: pluginsdk.ExpandStringSlice(administrators),
			},
			Mode: &mode,
		},
		Sku: capacities.CapacitySku{
			Name: d.Get("sku_name").(string),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourcePowerBIEmbeddedRead(d, meta)
}
