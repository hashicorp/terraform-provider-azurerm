package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointAuthenticationSasl struct {
	SaslType  DataflowEndpointAuthenticationSaslType `json:"saslType"`
	SecretRef string                                 `json:"secretRef"`
}
