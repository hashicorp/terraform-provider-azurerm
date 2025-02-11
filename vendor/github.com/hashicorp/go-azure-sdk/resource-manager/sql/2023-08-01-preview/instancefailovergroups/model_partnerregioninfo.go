package instancefailovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartnerRegionInfo struct {
	Location        *string                               `json:"location,omitempty"`
	ReplicationRole *InstanceFailoverGroupReplicationRole `json:"replicationRole,omitempty"`
}
