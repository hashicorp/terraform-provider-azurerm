package apischema

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchemaContractProperties struct {
	ContentType string                   `json:"contentType"`
	Document    SchemaDocumentProperties `json:"document"`
}
