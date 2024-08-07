package instancefailovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceFailoverGroupProperties struct {
	ManagedInstancePairs []ManagedInstancePairInfo              `json:"managedInstancePairs"`
	PartnerRegions       []PartnerRegionInfo                    `json:"partnerRegions"`
	ReadOnlyEndpoint     *InstanceFailoverGroupReadOnlyEndpoint `json:"readOnlyEndpoint,omitempty"`
	ReadWriteEndpoint    InstanceFailoverGroupReadWriteEndpoint `json:"readWriteEndpoint"`
	ReplicationRole      *InstanceFailoverGroupReplicationRole  `json:"replicationRole,omitempty"`
	ReplicationState     *string                                `json:"replicationState,omitempty"`
	SecondaryType        *SecondaryInstanceType                 `json:"secondaryType,omitempty"`
}
