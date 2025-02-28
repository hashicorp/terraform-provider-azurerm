package emailtemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailTemplateUpdateParameterProperties struct {
	Body        *string                                      `json:"body,omitempty"`
	Description *string                                      `json:"description,omitempty"`
	Parameters  *[]EmailTemplateParametersContractProperties `json:"parameters,omitempty"`
	Subject     *string                                      `json:"subject,omitempty"`
	Title       *string                                      `json:"title,omitempty"`
}
