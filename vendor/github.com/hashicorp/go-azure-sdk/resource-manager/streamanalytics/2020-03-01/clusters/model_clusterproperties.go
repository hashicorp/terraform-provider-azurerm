package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	CapacityAllocated *int64                    `json:"capacityAllocated,omitempty"`
	CapacityAssigned  *int64                    `json:"capacityAssigned,omitempty"`
	ClusterId         *string                   `json:"clusterId,omitempty"`
	CreatedDate       *string                   `json:"createdDate,omitempty"`
	ProvisioningState *ClusterProvisioningState `json:"provisioningState,omitempty"`
}

func (o *ClusterProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}
