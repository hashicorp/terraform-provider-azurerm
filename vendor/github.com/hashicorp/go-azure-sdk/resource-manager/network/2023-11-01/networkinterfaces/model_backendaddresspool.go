package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendAddressPool struct {
	Etag       *string                             `json:"etag,omitempty"`
	Id         *string                             `json:"id,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *BackendAddressPoolPropertiesFormat `json:"properties,omitempty"`
	Type       *string                             `json:"type,omitempty"`
}
