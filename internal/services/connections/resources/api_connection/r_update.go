// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_connection

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiConnectionUpdate(d *schema.ResourceData, meta any) error {
	client := meta.(*clients.Client).Connections.ConnectionsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connections.ParseConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}
	props := existing.Model.Properties

	if d.HasChange("display_name") {
		props.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	// The GET operation returns `NonSecretParameterValues` but we're making updates through `ParameterValues`
	// so we remove `NonSecretParameterValues` from the request to avoid conflicting parameters.
	// this is fixed in later (preview) versions of the API but these don't have an API spec available.
	props.NonSecretParameterValues = nil
	if d.HasChange("parameter_values") {
		props.ParameterValues = pointer.To(d.Get("parameter_values").(map[string]any))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceApiConnectionRead(d, meta)
}
