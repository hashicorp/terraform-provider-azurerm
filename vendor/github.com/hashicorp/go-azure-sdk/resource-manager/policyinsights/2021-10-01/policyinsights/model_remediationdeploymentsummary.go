package policyinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RemediationDeploymentSummary struct {
	FailedDeployments     *int64 `json:"failedDeployments,omitempty"`
	SuccessfulDeployments *int64 `json:"successfulDeployments,omitempty"`
	TotalDeployments      *int64 `json:"totalDeployments,omitempty"`
}
