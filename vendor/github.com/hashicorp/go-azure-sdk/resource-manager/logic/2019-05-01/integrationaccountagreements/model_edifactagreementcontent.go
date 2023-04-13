package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdifactAgreementContent struct {
	ReceiveAgreement EdifactOneWayAgreement `json:"receiveAgreement"`
	SendAgreement    EdifactOneWayAgreement `json:"sendAgreement"`
}
