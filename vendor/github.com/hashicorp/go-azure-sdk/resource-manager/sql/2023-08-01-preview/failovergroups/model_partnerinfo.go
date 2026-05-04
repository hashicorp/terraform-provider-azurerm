package failovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerInfo struct {
	Id              string                        `json:"id"`
	Location        *string                       `json:"location,omitempty"`
	ReplicationRole *FailoverGroupReplicationRole `json:"replicationRole,omitempty"`
}
