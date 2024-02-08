package emailtemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailTemplateContractProperties struct {
	Body        string                                       `json:"body"`
	Description *string                                      `json:"description,omitempty"`
	IsDefault   *bool                                        `json:"isDefault,omitempty"`
	Parameters  *[]EmailTemplateParametersContractProperties `json:"parameters,omitempty"`
	Subject     string                                       `json:"subject"`
	Title       *string                                      `json:"title,omitempty"`
}
