package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterReportedProperties struct {
	ClusterId      *string        `json:"clusterId,omitempty"`
	ClusterName    *string        `json:"clusterName,omitempty"`
	ClusterVersion *string        `json:"clusterVersion,omitempty"`
	LastUpdated    *string        `json:"lastUpdated,omitempty"`
	Nodes          *[]ClusterNode `json:"nodes,omitempty"`
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
