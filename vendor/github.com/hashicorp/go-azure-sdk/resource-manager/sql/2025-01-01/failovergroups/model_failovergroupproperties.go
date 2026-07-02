package failovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverGroupProperties struct {
	Databases         *[]string                            `json:"databases,omitempty"`
	PartnerServers    []PartnerInfo                        `json:"partnerServers"`
	ReadOnlyEndpoint  *FailoverGroupReadOnlyEndpoint       `json:"readOnlyEndpoint,omitempty"`
	ReadWriteEndpoint FailoverGroupReadWriteEndpoint       `json:"readWriteEndpoint"`
	ReplicationRole   *FailoverGroupReplicationRole        `json:"replicationRole,omitempty"`
	ReplicationState  *string                              `json:"replicationState,omitempty"`
	SecondaryType     *FailoverGroupDatabasesSecondaryType `json:"secondaryType,omitempty"`
}
