package networkmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerDeploymentStatusParameter struct {
	DeploymentTypes *[]ConfigurationType `json:"deploymentTypes,omitempty"`
	Regions         *[]string            `json:"regions,omitempty"`
	SkipToken       *string              `json:"skipToken,omitempty"`
}
