package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainProperties struct {
	AppName           *string                                `json:"appName,omitempty"`
	CertName          *string                                `json:"certName,omitempty"`
	ProvisioningState *CustomDomainResourceProvisioningState `json:"provisioningState,omitempty"`
	Thumbprint        *string                                `json:"thumbprint,omitempty"`
}
