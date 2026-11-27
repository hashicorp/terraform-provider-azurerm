// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_workspace

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	healthcareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisWorkspaceUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := healthcareClient.From(meta).HealthcareWorkspaceClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := workspaces.ParseWorkspaceID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]any)

	parameters := workspaces.ResourceTags{
		Tags: tags.Expand(t),
	}

	if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceHealthcareApisWorkspaceRead(d, meta)
}
