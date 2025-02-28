package contenttype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentTypeContractProperties struct {
	Description *string      `json:"description,omitempty"`
	Id          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Schema      *interface{} `json:"schema,omitempty"`
	Version     *string      `json:"version,omitempty"`
}
