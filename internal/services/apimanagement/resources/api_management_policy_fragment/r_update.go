// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_policy_fragment

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/policyfragment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementPolicyFragmentUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.PolicyFragmentClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := policyfragment.ParsePolicyFragmentID(d.Id())
	if err != nil {
		return err
	}

	format := policyfragment.PolicyFragmentContentFormat(d.Get("format").(string))

	opts := policyfragment.DefaultGetOperationOptions()
	opts.Format = pointer.To(format)
	existing, err := client.Get(ctx, *id, opts)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("description") {
		payload.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("value") {
		payload.Properties.Value = d.Get("value").(string)
	}

	if d.HasChange("format") {
		payload.Properties.Format = pointer.To(format)
		// on format change we send the value also because it might be interpreted differently because of the changed format
		payload.Properties.Value = d.Get("value").(string)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload, policyfragment.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementPolicyFragmentRead(d, meta)
}
