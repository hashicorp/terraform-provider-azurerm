package brokerauthentication

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerAuthenticatorMethodCustom struct {
	Auth            *BrokerAuthenticatorCustomAuth `json:"auth,omitempty"`
	CaCertConfigMap *string                        `json:"caCertConfigMap,omitempty"`
	Endpoint        string                         `json:"endpoint"`
	Headers         *map[string]string             `json:"headers,omitempty"`
}
