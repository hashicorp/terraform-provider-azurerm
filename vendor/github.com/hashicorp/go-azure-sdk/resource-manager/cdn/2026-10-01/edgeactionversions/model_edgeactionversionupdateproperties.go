package edgeactionversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionVersionUpdateProperties struct {
	DeploymentType   *EdgeActionVersionDeploymentType `json:"deploymentType,omitempty"`
	IsDefaultVersion *EdgeActionIsDefaultVersion      `json:"isDefaultVersion,omitempty"`
}
