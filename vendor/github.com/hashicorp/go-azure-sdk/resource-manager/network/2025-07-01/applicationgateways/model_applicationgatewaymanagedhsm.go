package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayManagedHsm struct {
	KeyId          *string `json:"keyId,omitempty"`
	PublicCertData *string `json:"publicCertData,omitempty"`
}
