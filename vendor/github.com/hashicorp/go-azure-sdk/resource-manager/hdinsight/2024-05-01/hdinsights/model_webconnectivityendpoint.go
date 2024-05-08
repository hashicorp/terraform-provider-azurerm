package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebConnectivityEndpoint struct {
	Fqdn        string  `json:"fqdn"`
	PrivateFqdn *string `json:"privateFqdn,omitempty"`
}
