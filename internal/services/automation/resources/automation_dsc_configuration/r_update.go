// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_dsc_configuration

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/dscconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationDscConfigurationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.DscConfiguration
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dscconfiguration.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", *id)
	}

	content, err := client.GetContent(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving content for %s: %+v", *id, err)
	}

	if content.Model == nil {
		return fmt.Errorf("retrieving content for %s: response was nil", *id)
	}

	// NOTE: We use CreateOrUpdate (PUT) rather than Update (PATCH) here because the
	// PATCH API re-parses the DSC configuration source content and resets the description
	// field, causing updates to description to be silently lost.
	parameters := dscconfiguration.DscConfigurationCreateOrUpdateParameters{
		Properties: dscconfiguration.DscConfigurationCreateOrUpdateProperties{
			LogVerbose:  existing.Model.Properties.LogVerbose,
			Description: existing.Model.Properties.Description,
			// The GET response does not contain source, so we must set it based on the retrieved content to prevent 400s
			Source: dscconfiguration.ContentSource{
				Type:  pointer.To(dscconfiguration.ContentSourceTypeEmbeddedContent),
				Value: pointer.To(string(*content.Model)),
			},
		},
		Location: pointer.To(existing.Model.Location),
		Tags:     existing.Model.Tags,
	}

	if d.HasChange("content_embedded") {
		parameters.Properties.Source.Value = pointer.To(d.Get("content_embedded").(string))
	}

	if d.HasChange("log_verbose") {
		parameters.Properties.LogVerbose = pointer.To(d.Get("log_verbose").(bool))
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("tags") {
		parameters.Tags = pointer.To(helpers.ExpandStringInterfaceMap(d.Get("tags").(map[string]any)))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationDscConfigurationRead(d, meta)
}
