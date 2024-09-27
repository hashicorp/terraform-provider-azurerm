package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentSettingsProperties struct {
	ArcNodeResourceIds      []string                `json:"arcNodeResourceIds"`
	DeploymentConfiguration DeploymentConfiguration `json:"deploymentConfiguration"`
	DeploymentMode          DeploymentMode          `json:"deploymentMode"`
	ProvisioningState       *ProvisioningState      `json:"provisioningState,omitempty"`
	ReportedProperties      *ReportedProperties     `json:"reportedProperties,omitempty"`
}
