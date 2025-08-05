package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestContract struct {
	Description     *string                   `json:"description,omitempty"`
	Headers         *[]ParameterContract      `json:"headers,omitempty"`
	QueryParameters *[]ParameterContract      `json:"queryParameters,omitempty"`
	Representations *[]RepresentationContract `json:"representations,omitempty"`
}
