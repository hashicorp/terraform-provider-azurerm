package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectStatusDefinition struct {
	AppliedBy             *ObjectReferenceDefinition         `json:"appliedBy,omitempty"`
	ComplianceState       *FluxComplianceState               `json:"complianceState,omitempty"`
	HelmReleaseProperties *HelmReleasePropertiesDefinition   `json:"helmReleaseProperties,omitempty"`
	Kind                  *string                            `json:"kind,omitempty"`
	Name                  *string                            `json:"name,omitempty"`
	Namespace             *string                            `json:"namespace,omitempty"`
	StatusConditions      *[]ObjectStatusConditionDefinition `json:"statusConditions,omitempty"`
}
