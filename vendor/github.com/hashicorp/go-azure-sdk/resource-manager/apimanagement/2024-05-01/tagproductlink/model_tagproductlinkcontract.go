package tagproductlink

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagProductLinkContract struct {
	Id         *string                           `json:"id,omitempty"`
	Name       *string                           `json:"name,omitempty"`
	Properties *TagProductLinkContractProperties `json:"properties,omitempty"`
	Type       *string                           `json:"type,omitempty"`
}
