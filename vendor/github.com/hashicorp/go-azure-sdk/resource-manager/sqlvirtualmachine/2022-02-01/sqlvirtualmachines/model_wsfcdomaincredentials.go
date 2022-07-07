package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WsfcDomainCredentials struct {
	ClusterBootstrapAccountPassword *string `json:"clusterBootstrapAccountPassword,omitempty"`
	ClusterOperatorAccountPassword  *string `json:"clusterOperatorAccountPassword,omitempty"`
	SqlServiceAccountPassword       *string `json:"sqlServiceAccountPassword,omitempty"`
}
