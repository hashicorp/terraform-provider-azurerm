// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package security_center_assessment_policy

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func securityCenterAssessmentPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"severity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(security.SeverityMedium),
			ValidateFunc: validation.StringInEnumSlice(security.PossibleSeverityValues(), false),
		},

		"categories": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			// Note: O+C the API would return `Unknown` when `categories` isn't set.
			// After synced with service team, they confirmed will add `Unknown` as possible value to this property and it will be published as a new version of this API.
			// https://github.com/Azure/azure-rest-api-specs/issues/14918
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"Unknown",
					string(security.Compute),
					string(security.Data),
					string(security.IdentityAndAccess),
					string(security.IoT),
					string(security.Networking),
				}, false),
			},
		},

		"implementation_effort": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInEnumSlice(security.PossibleImplementationEffortValues(), false),
		},

		"remediation_description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"threats": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"AccountBreach",
					"DataExfiltration",
					"DataSpillage",
					"MaliciousInsider",
					"ElevationOfPrivilege",
					"ThreatResistance",
					"MissingCoverage",
					"DenialOfService",
				}, false),
			},
		},

		"user_impact": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInEnumSlice(security.PossibleUserImpactValues(), false),
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
