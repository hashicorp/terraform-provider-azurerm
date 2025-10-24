package brokerauthentication

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerAuthenticatorMethods struct {
	CustomSettings              *BrokerAuthenticatorMethodCustom `json:"customSettings,omitempty"`
	Method                      BrokerAuthenticationMethod       `json:"method"`
	ServiceAccountTokenSettings *BrokerAuthenticatorMethodSat    `json:"serviceAccountTokenSettings,omitempty"`
	X509Settings                *BrokerAuthenticatorMethodX509   `json:"x509Settings,omitempty"`
}
