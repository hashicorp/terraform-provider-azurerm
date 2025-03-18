package routetables

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceConnection struct {
	Etag       *string                                 `json:"etag,omitempty"`
	Id         *string                                 `json:"id,omitempty"`
	Name       *string                                 `json:"name,omitempty"`
	Properties *PrivateLinkServiceConnectionProperties `json:"properties,omitempty"`
	Type       *string                                 `json:"type,omitempty"`
}
