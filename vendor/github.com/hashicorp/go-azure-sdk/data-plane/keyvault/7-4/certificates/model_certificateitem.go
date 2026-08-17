package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateItem struct {
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	X5t        *string                `json:"x5t,omitempty"`
}
