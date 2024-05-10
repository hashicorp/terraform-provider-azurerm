package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretProperties struct {
	Attributes           *Attributes `json:"attributes,omitempty"`
	ContentType          *string     `json:"contentType,omitempty"`
	SecretUri            *string     `json:"secretUri,omitempty"`
	SecretUriWithVersion *string     `json:"secretUriWithVersion,omitempty"`
	Value                *string     `json:"value,omitempty"`
}
