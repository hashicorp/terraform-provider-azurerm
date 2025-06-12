package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPodIdentityProvisioningError struct {
	Error *ManagedClusterPodIdentityProvisioningErrorBody `json:"error,omitempty"`
}
