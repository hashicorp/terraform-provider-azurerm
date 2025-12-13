package broker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertManagerCertOptions struct {
	Duration    string                `json:"duration"`
	PrivateKey  CertManagerPrivateKey `json:"privateKey"`
	RenewBefore string                `json:"renewBefore"`
}
