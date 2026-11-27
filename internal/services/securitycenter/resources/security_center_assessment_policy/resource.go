// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/security/2021-06-01/assessmentsmetadata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages the Security Center Assessment Metadata for Azure Security Center.",
		Create:      resourceArmSecurityCenterAssessmentPolicyCreate,
		Read:        resourceArmSecurityCenterAssessmentPolicyRead,
		Update:      resourceArmSecurityCenterAssessmentPolicyUpdate,
		Delete:      resourceArmSecurityCenterAssessmentPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assessmentsmetadata.ParseProviderAssessmentMetadataID(id)
			return err
		}),

		Schema: securityCenterAssessmentPolicySchema(),
	}
}
