package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2AcknowledgementConnectionSettings struct {
	IgnoreCertificateNameMismatch bool `json:"ignoreCertificateNameMismatch"`
	KeepHTTPConnectionAlive       bool `json:"keepHttpConnectionAlive"`
	SupportHTTPStatusCodeContinue bool `json:"supportHttpStatusCodeContinue"`
	UnfoldHTTPHeaders             bool `json:"unfoldHttpHeaders"`
}
