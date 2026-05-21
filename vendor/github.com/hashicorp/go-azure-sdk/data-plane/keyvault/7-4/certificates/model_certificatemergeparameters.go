package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateMergeParameters struct {
	Attributes *CertificateAttributes `json:"attributes,omitempty"`
	Tags       *map[string]string     `json:"tags,omitempty"`
	X5c        []string               `json:"x5c"`
}
