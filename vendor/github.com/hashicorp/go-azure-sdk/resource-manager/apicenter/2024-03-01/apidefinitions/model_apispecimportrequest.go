package apidefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiSpecImportRequest struct {
	Format        *ApiSpecImportSourceFormat         `json:"format,omitempty"`
	Specification *ApiSpecImportRequestSpecification `json:"specification,omitempty"`
	Value         *string                            `json:"value,omitempty"`
}
