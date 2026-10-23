package nginxdeploymentwafpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyProperties struct {
	ApplyingState     *NginxDeploymentWafPolicyApplyingStatus  `json:"applyingState,omitempty"`
	CompilingState    *NginxDeploymentWafPolicyCompilingStatus `json:"compilingState,omitempty"`
	Content           *string                                  `json:"content,omitempty"`
	Filepath          *string                                  `json:"filepath,omitempty"`
	ProvisioningState *ProvisioningState                       `json:"provisioningState,omitempty"`
}
