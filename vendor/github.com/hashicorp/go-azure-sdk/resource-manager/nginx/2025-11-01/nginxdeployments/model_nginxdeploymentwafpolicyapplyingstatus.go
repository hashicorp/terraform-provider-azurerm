package nginxdeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyApplyingStatus struct {
	Code          *NginxDeploymentWafPolicyApplyingStatusCode `json:"code,omitempty"`
	DisplayStatus *string                                     `json:"displayStatus,omitempty"`
	Time          *string                                     `json:"time,omitempty"`
}
