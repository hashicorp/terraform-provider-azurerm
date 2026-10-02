package nginxdeploymentwafpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyError struct {
	Code    *string `json:"code,omitempty"`
	Field   *string `json:"field,omitempty"`
	Message *string `json:"message,omitempty"`
}
