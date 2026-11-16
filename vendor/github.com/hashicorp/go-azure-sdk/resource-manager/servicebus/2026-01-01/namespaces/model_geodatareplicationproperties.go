package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoDataReplicationProperties struct {
	Locations                          *[]NamespaceReplicaLocation `json:"locations,omitempty"`
	MaxReplicationLagDurationInSeconds *int64                      `json:"maxReplicationLagDurationInSeconds,omitempty"`
}
