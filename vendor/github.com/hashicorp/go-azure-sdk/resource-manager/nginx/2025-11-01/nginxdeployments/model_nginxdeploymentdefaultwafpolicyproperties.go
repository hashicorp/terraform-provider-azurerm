package nginxdeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentDefaultWafPolicyProperties struct {
	Content  *string `json:"content,omitempty"`
	Filepath *string `json:"filepath,omitempty"`
}
