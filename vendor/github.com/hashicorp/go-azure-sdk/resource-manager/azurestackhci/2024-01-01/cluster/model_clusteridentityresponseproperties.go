package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterIdentityResponseProperties struct {
	AadApplicationObjectId      *string `json:"aadApplicationObjectId,omitempty"`
	AadClientId                 *string `json:"aadClientId,omitempty"`
	AadServicePrincipalObjectId *string `json:"aadServicePrincipalObjectId,omitempty"`
	AadTenantId                 *string `json:"aadTenantId,omitempty"`
}
