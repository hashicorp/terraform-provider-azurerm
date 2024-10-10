package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToTargetSqlMITaskOutput struct {
	AgentJobs                *[]string              `json:"agentJobs,omitempty"`
	Id                       *string                `json:"id,omitempty"`
	Logins                   *[]string              `json:"logins,omitempty"`
	TargetServerBrandVersion *string                `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion      *string                `json:"targetServerVersion,omitempty"`
	ValidationErrors         *[]ReportableException `json:"validationErrors,omitempty"`
}
