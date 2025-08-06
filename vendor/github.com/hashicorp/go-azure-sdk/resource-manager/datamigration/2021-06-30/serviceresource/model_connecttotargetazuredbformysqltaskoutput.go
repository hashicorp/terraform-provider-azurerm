package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectToTargetAzureDbForMySqlTaskOutput struct {
	Databases                *[]string              `json:"databases,omitempty"`
	Id                       *string                `json:"id,omitempty"`
	ServerVersion            *string                `json:"serverVersion,omitempty"`
	TargetServerBrandVersion *string                `json:"targetServerBrandVersion,omitempty"`
	ValidationErrors         *[]ReportableException `json:"validationErrors,omitempty"`
}
