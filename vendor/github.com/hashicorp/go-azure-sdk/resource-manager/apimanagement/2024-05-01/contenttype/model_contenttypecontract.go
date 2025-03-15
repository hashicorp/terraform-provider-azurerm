package contenttype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentTypeContract struct {
	Id         *string                        `json:"id,omitempty"`
	Name       *string                        `json:"name,omitempty"`
	Properties *ContentTypeContractProperties `json:"properties,omitempty"`
	Type       *string                        `json:"type,omitempty"`
}
