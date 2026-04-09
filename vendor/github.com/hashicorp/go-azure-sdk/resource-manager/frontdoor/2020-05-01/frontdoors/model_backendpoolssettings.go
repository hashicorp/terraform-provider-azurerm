package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendPoolsSettings struct {
	EnforceCertificateNameCheck *EnforceCertificateNameCheckEnabledState `json:"enforceCertificateNameCheck,omitempty"`
	SendRecvTimeoutSeconds      *int64                                   `json:"sendRecvTimeoutSeconds,omitempty"`
}
