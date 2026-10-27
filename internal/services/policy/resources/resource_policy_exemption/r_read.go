// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_policy_exemption

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourcePolicyExemptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourcePolicyExemptionID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Exemption: %+v", err)
	}

	resp, err := client.Get(ctx, id.ResourceId, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Error reading Policy Exemption %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_id", id.ResourceId)
	if props := resp.ExemptionProperties; props != nil {
		d.Set("policy_assignment_id", props.PolicyAssignmentID)
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)
		d.Set("exemption_category", string(props.ExemptionCategory))

		if err := d.Set("policy_definition_reference_ids", pluginsdk.FlattenSlice(props.PolicyDefinitionReferenceIds)); err != nil {
			return fmt.Errorf("setting `policy_definition_reference_ids: %+v", err)
		}

		expiresOn := ""
		if expiresTime := props.ExpiresOn; expiresTime != nil {
			expiresOn = expiresTime.String()
		}
		d.Set("expires_on", expiresOn)

		if metadataStr := helpers.FlattenJSON(props.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}
	}

	return nil
}
