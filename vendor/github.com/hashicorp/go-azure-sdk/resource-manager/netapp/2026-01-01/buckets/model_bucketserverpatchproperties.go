package buckets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BucketServerPatchProperties struct {
	CertificateObject           *string                      `json:"certificateObject,omitempty"`
	Fqdn                        *string                      `json:"fqdn,omitempty"`
	OnCertificateConflictAction *OnCertificateConflictAction `json:"onCertificateConflictAction,omitempty"`
}
