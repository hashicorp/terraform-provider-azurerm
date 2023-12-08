package apioperation

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResponseContract struct {
	Description     *string                   `json:"description,omitempty"`
	Headers         *[]ParameterContract      `json:"headers,omitempty"`
	Representations *[]RepresentationContract `json:"representations,omitempty"`
	StatusCode      int64                     `json:"statusCode"`
}
