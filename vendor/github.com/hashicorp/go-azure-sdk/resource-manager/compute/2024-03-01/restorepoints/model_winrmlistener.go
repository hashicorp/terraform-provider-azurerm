package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WinRMListener struct {
	CertificateURL *string        `json:"certificateUrl,omitempty"`
	Protocol       *ProtocolTypes `json:"protocol,omitempty"`
}
