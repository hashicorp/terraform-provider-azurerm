// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package digital_twins_instance

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/digitaltwinsinstance"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDigitalTwinsInstanceCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).DigitalTwins.InstanceClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := digitaltwinsinstance.NewDigitalTwinsInstanceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.DigitalTwinsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_digital_twins_instance", id.ID())
		}
	}

	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	payload := digitaltwinsinstance.DigitalTwinsDescription{
		Location: location.Normalize(d.Get("location").(string)),
		Identity: expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.DigitalTwinsCreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDigitalTwinsInstanceRead(d, meta)
}
