// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_credential

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/credential"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationCredentialUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.Credential
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := credential.ParseCredentialID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	parameters := credential.CredentialCreateOrUpdateParameters{
		Name: id.CredentialName,
		Properties: credential.CredentialCreateOrUpdateProperties{
			// password is sensitive and not returned by the API, so it must always
			// be read from state to ensure the value is preserved on update
			Password: d.Get("password").(string),
		},
	}

	if existing.Model != nil && existing.Model.Properties != nil {
		parameters.Properties.UserName = pointer.From(existing.Model.Properties.UserName)
		parameters.Properties.Description = existing.Model.Properties.Description
	}

	if d.HasChange("username") {
		parameters.Properties.UserName = d.Get("username").(string)
	}

	if d.HasChange("description") {
		parameters.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationCredentialRead(d, meta)
}
