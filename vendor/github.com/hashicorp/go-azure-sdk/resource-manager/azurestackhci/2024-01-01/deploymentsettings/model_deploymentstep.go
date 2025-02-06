package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentStep struct {
	Description   *string           `json:"description,omitempty"`
	EndTimeUtc    *string           `json:"endTimeUtc,omitempty"`
	Exception     *[]string         `json:"exception,omitempty"`
	FullStepIndex *string           `json:"fullStepIndex,omitempty"`
	Name          *string           `json:"name,omitempty"`
	StartTimeUtc  *string           `json:"startTimeUtc,omitempty"`
	Status        *string           `json:"status,omitempty"`
	Steps         *[]DeploymentStep `json:"steps,omitempty"`
}
