package deploymentsafeguards

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSafeguardsProperties struct {
	ExcludedNamespaces        *[]string                  `json:"excludedNamespaces,omitempty"`
	Level                     DeploymentSafeguardsLevel  `json:"level"`
	PodSecurityStandardsLevel *PodSecurityStandardsLevel `json:"podSecurityStandardsLevel,omitempty"`
	ProvisioningState         *ProvisioningState         `json:"provisioningState,omitempty"`
	SystemExcludedNamespaces  []string                   `json:"systemExcludedNamespaces"`
}
