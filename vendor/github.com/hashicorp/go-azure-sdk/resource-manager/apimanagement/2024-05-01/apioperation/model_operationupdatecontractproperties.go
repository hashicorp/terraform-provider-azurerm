package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationUpdateContractProperties struct {
	Description        *string              `json:"description,omitempty"`
	DisplayName        *string              `json:"displayName,omitempty"`
	Method             *string              `json:"method,omitempty"`
	Policies           *string              `json:"policies,omitempty"`
	Request            *RequestContract     `json:"request,omitempty"`
	Responses          *[]ResponseContract  `json:"responses,omitempty"`
	TemplateParameters *[]ParameterContract `json:"templateParameters,omitempty"`
	UrlTemplate        *string              `json:"urlTemplate,omitempty"`
}
