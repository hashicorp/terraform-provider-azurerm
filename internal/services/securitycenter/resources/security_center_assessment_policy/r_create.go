// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"fmt"

	//nolint:staticcheck
	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmSecurityCenterAssessmentPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).SecurityCenter.AssessmentsMetadataClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := uuid.New().String()

	id := assessmentsmetadata.NewProviderAssessmentMetadataID(subscriptionId, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.GetInSubscription(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_security_center_assessment_policy", id.ID())
		}
	}

	params := assessmentsmetadata.SecurityAssessmentMetadataResponse{
		Properties: &assessmentsmetadata.SecurityAssessmentMetadataPropertiesResponse{
			AssessmentType: assessmentsmetadata.AssessmentTypeCustomerManaged,
			Description:    pointer.To(d.Get("description").(string)),
			DisplayName:    d.Get("display_name").(string),
			Severity:       assessmentsmetadata.Severity(d.Get("severity").(string)),
		},
	}

	if v, ok := d.GetOk("categories"); ok {
		categories := make([]assessmentsmetadata.Categories, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			categories = append(categories, assessmentsmetadata.Categories(item.(string)))
		}
		params.Properties.Categories = &categories
	}

	if v, ok := d.GetOk("threats"); ok {
		threats := make([]assessmentsmetadata.Threats, 0)
		for _, item := range v.(*pluginsdk.Set).List() {
			threats = append(threats, assessmentsmetadata.Threats(item.(string)))
		}
		params.Properties.Threats = &threats
	}

	if v, ok := d.GetOk("implementation_effort"); ok {
		params.Properties.ImplementationEffort = pointer.ToEnum[assessmentsmetadata.ImplementationEffort](v.(string))
	}

	if v, ok := d.GetOk("remediation_description"); ok {
		params.Properties.RemediationDescription = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("user_impact"); ok {
		params.Properties.UserImpact = pointer.ToEnum[assessmentsmetadata.UserImpact](v.(string))
	}

	if _, err := client.CreateInSubscription(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmSecurityCenterAssessmentPolicyRead(d, meta)
}
