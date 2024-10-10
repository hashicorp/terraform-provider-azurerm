package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToTargetAzureDbForPostgreSqlSyncTaskOutput struct {
	Databases                *[]string              `json:"databases,omitempty"`
	Id                       *string                `json:"id,omitempty"`
	TargetServerBrandVersion *string                `json:"targetServerBrandVersion,omitempty"`
	TargetServerVersion      *string                `json:"targetServerVersion,omitempty"`
	ValidationErrors         *[]ReportableException `json:"validationErrors,omitempty"`
}
