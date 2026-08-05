package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretItem struct {
	Attributes  *SecretAttributes  `json:"attributes,omitempty"`
	ContentType *string            `json:"contentType,omitempty"`
	Id          *string            `json:"id,omitempty"`
	Managed     *bool              `json:"managed,omitempty"`
	Tags        *map[string]string `json:"tags,omitempty"`
}
