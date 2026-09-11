package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPodIdentity struct {
	BindingSelector   *string                                     `json:"bindingSelector,omitempty"`
	Identity          UserAssignedIdentity                        `json:"identity"`
	Name              string                                      `json:"name"`
	Namespace         string                                      `json:"namespace"`
	ProvisioningInfo  *ManagedClusterPodIdentityProvisioningInfo  `json:"provisioningInfo,omitempty"`
	ProvisioningState *ManagedClusterPodIdentityProvisioningState `json:"provisioningState,omitempty"`
}
