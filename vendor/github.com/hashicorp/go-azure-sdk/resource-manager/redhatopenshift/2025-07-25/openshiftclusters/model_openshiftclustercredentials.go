package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenShiftClusterCredentials struct {
	KubeadminPassword *string `json:"kubeadminPassword,omitempty"`
	KubeadminUsername *string `json:"kubeadminUsername,omitempty"`
}
