package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentRouting struct {
	Mode   *RoutingMode       `json:"mode,omitempty"`
	Models *[]DeploymentModel `json:"models,omitempty"`
}
