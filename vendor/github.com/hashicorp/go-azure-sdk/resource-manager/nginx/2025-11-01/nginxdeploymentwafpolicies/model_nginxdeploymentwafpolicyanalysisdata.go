package nginxdeploymentwafpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxDeploymentWafPolicyAnalysisData struct {
	Errors *[]NginxDeploymentWafPolicyError `json:"errors,omitempty"`
}
