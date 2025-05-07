package deployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScopedDeploymentWhatIf struct {
	Location   string                     `json:"location"`
	Properties DeploymentWhatIfProperties `json:"properties"`
}
