package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigDiagnosticsValidatorResultIssue struct {
	DescriptionParams *[]string `json:"descriptionParams,omitempty"`
	Id                *string   `json:"id,omitempty"`
}
