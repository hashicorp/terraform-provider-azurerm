package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationStatus struct {
	Status *string           `json:"status,omitempty"`
	Steps  *[]DeploymentStep `json:"steps,omitempty"`
}
