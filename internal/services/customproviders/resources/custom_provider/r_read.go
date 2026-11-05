// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCustomProviderRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).CustomProviders.CustomProviderClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customresourceprovider.ParseResourceProviderID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ResourceProviderName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("resource_type", flattenCustomProviderResourceType(props.ResourceTypes)); err != nil {
				return fmt.Errorf("setting `resource_type`: %+v", err)
			}

			if err := d.Set("action", flattenCustomProviderAction(props.Actions)); err != nil {
				return fmt.Errorf("setting `action`: %+v", err)
			}

			if err := d.Set("validation", flattenCustomProviderValidation(props.Validations)); err != nil {
				return fmt.Errorf("setting `validation`: %+v", err)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
