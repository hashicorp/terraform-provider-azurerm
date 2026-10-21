// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_application_definition

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedapplications/2021-07-01/applicationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedApplicationDefinitionUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ManagedApplication.ApplicationDefinitionClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationdefinitions.ParseApplicationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("display_name") {
		payload.Properties.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	if d.HasChange("package_enabled") {
		payload.Properties.IsEnabled = pointer.To(d.Get("package_enabled").(bool))
	}

	if d.HasChange("authorization") {
		payload.Properties.Authorizations = expandManagedApplicationDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List())
	}

	if d.HasChange("create_ui_definition") {
		// handle API error: The 'MainTemplate, CreateUiDefinition' properties should be empty if package zip file uri is provided.
		if v, ok := d.GetOk("create_ui_definition"); ok {
			payload.Properties.CreateUiDefinition = pointer.To(v)
		} else {
			payload.Properties.CreateUiDefinition = nil
		}
	}

	if d.HasChange("main_template") {
		// handle API error: The 'MainTemplate, CreateUiDefinition' properties should be empty if package zip file uri is provided.
		if v, ok := d.GetOk("main_template"); ok {
			payload.Properties.MainTemplate = pointer.To(v)
		} else {
			payload.Properties.MainTemplate = nil
		}
	}

	if d.HasChange("package_file_uri") {
		payload.Properties.PackageFileUri = pointer.To(d.Get("package_file_uri").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// update payload only supports tags, so we'll continue using CreateOrUpdate method here
	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceManagedApplicationDefinitionRead(d, meta)
}
