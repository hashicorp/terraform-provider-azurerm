package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplianceListCredentialResults struct {
	HybridConnectionConfig *HybridConnectionConfig          `json:"hybridConnectionConfig,omitempty"`
	Kubeconfigs            *[]ApplianceCredentialKubeconfig `json:"kubeconfigs,omitempty"`
}
