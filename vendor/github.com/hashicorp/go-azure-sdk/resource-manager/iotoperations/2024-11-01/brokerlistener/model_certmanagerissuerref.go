package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertManagerIssuerRef struct {
	Group string                `json:"group"`
	Kind  CertManagerIssuerKind `json:"kind"`
	Name  string                `json:"name"`
}
