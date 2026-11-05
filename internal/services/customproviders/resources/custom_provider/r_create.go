// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCustomProviderCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := location.Normalize(d.Get("location").(string))
	id := customresourceprovider.NewResourceProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_custom_resource_provider", id.ID())
			}
		}
	}

	provider := customresourceprovider.CustomRPManifest{
		Properties: &customresourceprovider.CustomRPManifestProperties{
			ResourceTypes: expandCustomProviderResourceType(d.Get("resource_type").(*pluginsdk.Set).List()),
			Actions:       expandCustomProviderAction(d.Get("action").(*pluginsdk.Set).List()),
			Validations:   expandCustomProviderValidation(d.Get("validation").(*pluginsdk.Set).List()),
		},
		Location: location,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, provider, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, provider); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceCustomProviderRead(d, meta)
}
