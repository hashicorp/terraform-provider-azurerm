package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPodIdentityException struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	PodLabels map[string]string `json:"podLabels"`
}
