package blobservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BlobServiceProperties struct {
	Id         *string                          `json:"id,omitempty"`
	Name       *string                          `json:"name,omitempty"`
	Properties *BlobServicePropertiesProperties `json:"properties,omitempty"`
	Sku        *Sku                             `json:"sku,omitempty"`
	Type       *string                          `json:"type,omitempty"`
}
