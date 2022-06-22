package connections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiReference struct {
	BrandColor  *string      `json:"brandColor,omitempty"`
	Description *string      `json:"description,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	IconUri     *string      `json:"iconUri,omitempty"`
	Id          *string      `json:"id,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Swagger     *interface{} `json:"swagger,omitempty"`
	Type        *string      `json:"type,omitempty"`
}
