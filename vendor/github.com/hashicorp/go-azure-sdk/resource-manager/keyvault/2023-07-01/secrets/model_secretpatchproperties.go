package secrets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretPatchProperties struct {
	Attributes  *Attributes `json:"attributes,omitempty"`
	ContentType *string     `json:"contentType,omitempty"`
	Value       *string     `json:"value,omitempty"`
}
