package connectedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CredentialResults struct {
	HybridConnectionConfig *HybridConnectionConfig `json:"hybridConnectionConfig,omitempty"`
	Kubeconfigs            *[]CredentialResult     `json:"kubeconfigs,omitempty"`
}
