// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSecurityCenterAssessmentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SecurityCenter.AssessmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	// The stored attribute needs no state migration as Read rewrites it with the canonical casing.
	metadataID, err := assessmentsmetadata.ParseProviderAssessmentMetadataIDInsensitively(d.Get("assessment_policy_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewAssessmentID(d.Get("target_resource_id").(string), metadataID.AssessmentMetadataName)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.TargetResourceID, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for present of existing Security Center Assessments %q : %+v", id.ID(), err)
				}
			}

			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_security_center_assessment", id.ID())
			}
		}
	}

	assessment := security.Assessment{
		AssessmentProperties: &security.AssessmentProperties{
			AdditionalData: pluginsdk.ExpandMapStringPtrString(d.Get("additional_data").(map[string]any)),
			ResourceDetails: &security.AzureResourceDetails{
				Source: security.SourceAzure,
			},
			Status: expandSecurityCenterAssessmentStatus(d.Get("status").([]any)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.TargetResourceID, id.Name, assessment); err != nil {
		return fmt.Errorf("creating/updating Security Center Assessment %q (target resource id %q) : %+v", id.Name, id.TargetResourceID, err)
	}

	d.SetId(id.ID())

	return resourceSecurityCenterAssessmentRead(d, meta)
}

func expandSecurityCenterAssessmentStatus(input []any) *security.AssessmentStatus {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	return &security.AssessmentStatus{
		Code:        security.AssessmentStatusCode(v["code"].(string)),
		Cause:       pointer.To(v["cause"].(string)),
		Description: pointer.To(v["description"].(string)),
	}
}
