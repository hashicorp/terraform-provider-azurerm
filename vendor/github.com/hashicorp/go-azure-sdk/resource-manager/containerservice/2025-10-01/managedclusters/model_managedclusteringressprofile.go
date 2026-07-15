package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterIngressProfile struct {
	WebAppRouting *ManagedClusterIngressProfileWebAppRouting `json:"webAppRouting,omitempty"`
}
