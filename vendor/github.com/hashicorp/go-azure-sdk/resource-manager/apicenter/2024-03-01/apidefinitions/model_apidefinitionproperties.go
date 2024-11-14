package apidefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiDefinitionProperties struct {
	Description   *string                               `json:"description,omitempty"`
	Specification *ApiDefinitionPropertiesSpecification `json:"specification,omitempty"`
	Title         string                                `json:"title"`
}
