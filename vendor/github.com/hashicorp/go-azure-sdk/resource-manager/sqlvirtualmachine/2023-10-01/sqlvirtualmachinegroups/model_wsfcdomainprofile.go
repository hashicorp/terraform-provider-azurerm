package sqlvirtualmachinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WsfcDomainProfile struct {
	ClusterBootstrapAccount  *string            `json:"clusterBootstrapAccount,omitempty"`
	ClusterOperatorAccount   *string            `json:"clusterOperatorAccount,omitempty"`
	ClusterSubnetType        *ClusterSubnetType `json:"clusterSubnetType,omitempty"`
	DomainFqdn               *string            `json:"domainFqdn,omitempty"`
	FileShareWitnessPath     *string            `json:"fileShareWitnessPath,omitempty"`
	IsSqlServiceAccountGmsa  *bool              `json:"isSqlServiceAccountGmsa,omitempty"`
	OuPath                   *string            `json:"ouPath,omitempty"`
	SqlServiceAccount        *string            `json:"sqlServiceAccount,omitempty"`
	StorageAccountPrimaryKey *string            `json:"storageAccountPrimaryKey,omitempty"`
	StorageAccountURL        *string            `json:"storageAccountUrl,omitempty"`
}
