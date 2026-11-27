// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterAssessmentRead(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).AssessmentsClient
	subscriptionID := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.AssessmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] security Center Assessment %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	d.Set("assessment_policy_id", assessmentsmetadata.NewProviderAssessmentMetadataID(subscriptionID, id.Name).ID())
	d.Set("target_resource_id", id.TargetResourceID)
	if props := resp.AssessmentProperties; props != nil {
		d.Set("additional_data", pluginsdk.FlattenMapStringPtrString(props.AdditionalData))
		if err := d.Set("status", flattenSecurityCenterAssessmentStatus(props.Status)); err != nil {
			return fmt.Errorf("setting `status`: %s", err)
		}
	}

	return nil
}

func flattenSecurityCenterAssessmentStatus(input *security.AssessmentStatus) []any {
	if input == nil {
		return []any{}
	}

	var cause, description string
	if input.Cause != nil {
		cause = *input.Cause
	}
	if input.Description != nil {
		description = *input.Description
	}

	return []any{
		map[string]any{
			"code":        string(input.Code),
			"cause":       cause,
			"description": description,
		},
	}
}
