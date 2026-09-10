package connectedregistries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TlsCertificateProperties struct {
	Location *string          `json:"location,omitempty"`
	Type     *CertificateType `json:"type,omitempty"`
}
