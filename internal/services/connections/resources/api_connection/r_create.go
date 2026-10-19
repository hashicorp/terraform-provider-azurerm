// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/managedapis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiConnectionCreate(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := connections.NewConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_connection", id.ID())
		}
	}

	managedAppId, err := managedapis.ParseManagedApiID(d.Get("managed_api_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `managed_app_id`: %+v", err)
	}
	location := location.Normalize(managedAppId.LocationName)
	model := connections.ApiConnectionDefinition{
		Location: pointer.To(location),
		Properties: &connections.ApiConnectionDefinitionProperties{
			Api: &connections.ApiReference{
				Id: pointer.To(managedAppId.ID()),
			},
			DisplayName:     pointer.To(d.Get("display_name").(string)),
			ParameterValues: pointer.To(d.Get("parameter_values").(map[string]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}
	if v := d.Get("display_name").(string); v != "" {
		model.Properties.DisplayName = pointer.To(v)
	}

	if _, err := client.CreateOrUpdate(ctx, id, model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApiConnectionRead(d, meta)
}
