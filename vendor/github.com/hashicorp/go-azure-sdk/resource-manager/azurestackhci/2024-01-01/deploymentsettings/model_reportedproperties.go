package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportedProperties struct {
	DeploymentStatus *DeploymentStatus `json:"deploymentStatus,omitempty"`
	ValidationStatus *ValidationStatus `json:"validationStatus,omitempty"`
}
