package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2AgreementContent struct {
	ReceiveAgreement AS2OneWayAgreement `json:"receiveAgreement"`
	SendAgreement    AS2OneWayAgreement `json:"sendAgreement"`
}
