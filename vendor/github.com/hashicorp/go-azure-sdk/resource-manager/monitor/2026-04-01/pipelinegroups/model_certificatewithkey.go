package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateWithKey struct {
	Certificate CertificateSource `json:"certificate"`
	PrivateKey  PrivateKeySource  `json:"privateKey"`
}
