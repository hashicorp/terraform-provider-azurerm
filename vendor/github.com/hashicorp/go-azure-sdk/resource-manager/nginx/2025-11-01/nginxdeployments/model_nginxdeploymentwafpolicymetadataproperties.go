package nginxdeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyMetadataProperties struct {
	ApplyingState     *NginxDeploymentWafPolicyApplyingStatus  `json:"applyingState,omitempty"`
	CompilingState    *NginxDeploymentWafPolicyCompilingStatus `json:"compilingState,omitempty"`
	Filepath          *string                                  `json:"filepath,omitempty"`
	ProvisioningState *ProvisioningState                       `json:"provisioningState,omitempty"`
}
