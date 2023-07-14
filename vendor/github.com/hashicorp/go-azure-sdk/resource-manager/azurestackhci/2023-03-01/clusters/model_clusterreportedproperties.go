package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterReportedProperties struct {
	ClusterId             *string          `json:"clusterId,omitempty"`
	ClusterName           *string          `json:"clusterName,omitempty"`
	ClusterType           *ClusterNodeType `json:"clusterType,omitempty"`
	ClusterVersion        *string          `json:"clusterVersion,omitempty"`
	DiagnosticLevel       *DiagnosticLevel `json:"diagnosticLevel,omitempty"`
	ImdsAttestation       *ImdsAttestation `json:"imdsAttestation,omitempty"`
	LastUpdated           *string          `json:"lastUpdated,omitempty"`
	Manufacturer          *string          `json:"manufacturer,omitempty"`
	Nodes                 *[]ClusterNode   `json:"nodes,omitempty"`
	SupportedCapabilities *[]string        `json:"supportedCapabilities,omitempty"`
}

func (o *ClusterReportedProperties) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterReportedProperties) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
