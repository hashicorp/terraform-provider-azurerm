package firewallpolicykubeselectorgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyKubeSelectorGroupProperties struct {
	NamespaceSelector *KubeLabelSelector `json:"namespaceSelector,omitempty"`
	PodSelector       *KubeLabelSelector `json:"podSelector,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
