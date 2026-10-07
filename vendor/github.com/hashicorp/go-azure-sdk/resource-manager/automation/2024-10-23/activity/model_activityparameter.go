package activity

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivityParameter struct {
	Description                     *string                           `json:"description,omitempty"`
	IsDynamic                       *bool                             `json:"isDynamic,omitempty"`
	IsMandatory                     *bool                             `json:"isMandatory,omitempty"`
	Name                            *string                           `json:"name,omitempty"`
	Position                        *int64                            `json:"position,omitempty"`
	Type                            *string                           `json:"type,omitempty"`
	ValidationSet                   *[]ActivityParameterValidationSet `json:"validationSet,omitempty"`
	ValueFromPipeline               *bool                             `json:"valueFromPipeline,omitempty"`
	ValueFromPipelineByPropertyName *bool                             `json:"valueFromPipelineByPropertyName,omitempty"`
	ValueFromRemainingArguments     *bool                             `json:"valueFromRemainingArguments,omitempty"`
}
