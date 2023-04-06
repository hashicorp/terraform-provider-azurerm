package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPodIdentityProvisioningErrorBody struct {
	Code    *string                                           `json:"code,omitempty"`
	Details *[]ManagedClusterPodIdentityProvisioningErrorBody `json:"details,omitempty"`
	Message *string                                           `json:"message,omitempty"`
	Target  *string                                           `json:"target,omitempty"`
}
