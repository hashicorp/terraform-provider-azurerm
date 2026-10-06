package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatePatchResource struct {
	Id         *string                             `json:"id,omitempty"`
	Kind       *string                             `json:"kind,omitempty"`
	Name       *string                             `json:"name,omitempty"`
	Properties *CertificatePatchResourceProperties `json:"properties,omitempty"`
	Tags       *map[string]string                  `json:"tags,omitempty"`
	Type       *string                             `json:"type,omitempty"`
}
