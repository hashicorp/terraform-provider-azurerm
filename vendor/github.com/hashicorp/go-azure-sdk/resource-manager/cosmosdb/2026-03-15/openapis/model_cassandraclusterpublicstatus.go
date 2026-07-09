package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraClusterPublicStatus struct {
	ConnectionErrors *[]ConnectionError                             `json:"connectionErrors,omitempty"`
	DataCenters      *[]CassandraClusterPublicStatusDataCentersItem `json:"dataCenters,omitempty"`
	ETag             *string                                        `json:"eTag,omitempty"`
	Errors           *[]CassandraError                              `json:"errors,omitempty"`
	ReaperStatus     *ManagedCassandraReaperStatus                  `json:"reaperStatus,omitempty"`
}
