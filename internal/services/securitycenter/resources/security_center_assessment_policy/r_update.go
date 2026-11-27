// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSecurityCenterAssessmentPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.GetInSubscription(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	if d.HasChange("description") {
		existing.Model.Properties.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("display_name") {
		existing.Model.Properties.DisplayName = d.Get("display_name").(string)
	}

	if d.HasChange("severity") {
		existing.Model.Properties.Severity = assessmentsmetadata.Severity(d.Get("severity").(string))
	}

	if d.HasChange("categories") {
		categories := make([]assessmentsmetadata.Categories, 0)
		for _, item := range d.Get("categories").(*pluginsdk.Set).List() {
			categories = append(categories, assessmentsmetadata.Categories(item.(string)))
		}
		existing.Model.Properties.Categories = &categories
	}

	if d.HasChange("threats") {
		threats := make([]assessmentsmetadata.Threats, 0)
		for _, item := range d.Get("threats").(*pluginsdk.Set).List() {
			threats = append(threats, assessmentsmetadata.Threats(item.(string)))
		}
		existing.Model.Properties.Threats = &threats
	}

	if d.HasChange("implementation_effort") {
		existing.Model.Properties.ImplementationEffort = pointer.ToEnum[assessmentsmetadata.ImplementationEffort](d.Get("implementation_effort").(string))
	}

	if d.HasChange("remediation_description") {
		existing.Model.Properties.RemediationDescription = pointer.To(d.Get("remediation_description").(string))
	}

	if d.HasChange("user_impact") {
		existing.Model.Properties.UserImpact = pointer.ToEnum[assessmentsmetadata.UserImpact](d.Get("user_impact").(string))
	}

	if _, err := client.CreateInSubscription(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmSecurityCenterAssessmentPolicyRead(d, meta)
}
