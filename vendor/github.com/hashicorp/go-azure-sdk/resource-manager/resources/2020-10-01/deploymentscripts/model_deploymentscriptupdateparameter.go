package deploymentscripts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScriptUpdateParameter struct {
	Id   *string            `json:"id,omitempty"`
	Name *string            `json:"name,omitempty"`
	Tags *map[string]string `json:"tags,omitempty"`
	Type *string            `json:"type,omitempty"`
}
