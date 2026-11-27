// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSecurityCenterAssessmentPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetInSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AssessmentMetadataName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("severity", string(props.Severity))
			d.Set("implementation_effort", pointer.FromEnum(props.ImplementationEffort))
			d.Set("remediation_description", pointer.From(props.RemediationDescription))
			d.Set("user_impact", pointer.FromEnum(props.UserImpact))

			categories := make([]string, 0)
			if props.Categories != nil {
				for _, item := range *props.Categories {
					categories = append(categories, string(item))
				}
			}
			d.Set("categories", pluginsdk.FlattenSlice(&categories))

			threats := make([]string, 0)
			if props.Threats != nil {
				for _, item := range *props.Threats {
					threats = append(threats, string(item))
				}
			}
			d.Set("threats", pluginsdk.FlattenSlice(&threats))
		}
	}

	return nil
}
