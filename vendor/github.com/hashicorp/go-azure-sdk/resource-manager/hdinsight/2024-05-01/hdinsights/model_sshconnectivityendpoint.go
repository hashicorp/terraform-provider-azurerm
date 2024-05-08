package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SshConnectivityEndpoint struct {
	Endpoint           string  `json:"endpoint"`
	PrivateSshEndpoint *string `json:"privateSshEndpoint,omitempty"`
}
