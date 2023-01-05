package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WinRMListener struct {
	CertificateUrl *string        `json:"certificateUrl,omitempty"`
	Protocol       *ProtocolTypes `json:"protocol,omitempty"`
}
