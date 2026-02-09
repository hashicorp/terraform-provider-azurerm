package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationContractProperties struct {
	Description        *string              `json:"description,omitempty"`
	DisplayName        string               `json:"displayName"`
	Method             string               `json:"method"`
	Policies           *string              `json:"policies,omitempty"`
	Request            *RequestContract     `json:"request,omitempty"`
	Responses          *[]ResponseContract  `json:"responses,omitempty"`
	TemplateParameters *[]ParameterContract `json:"templateParameters,omitempty"`
	UrlTemplate        string               `json:"urlTemplate"`
}
