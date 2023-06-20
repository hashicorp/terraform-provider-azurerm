package protectionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Settings struct {
	IsCompression    *bool   `json:"isCompression,omitempty"`
	Issqlcompression *bool   `json:"issqlcompression,omitempty"`
	TimeZone         *string `json:"timeZone,omitempty"`
}
