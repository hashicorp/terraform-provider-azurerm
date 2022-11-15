package namespaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RelayNamespaceProperties struct {
	CreatedAt          *string                `json:"createdAt,omitempty"`
	MetricId           *string                `json:"metricId,omitempty"`
	ProvisioningState  *ProvisioningStateEnum `json:"provisioningState,omitempty"`
	ServiceBusEndpoint *string                `json:"serviceBusEndpoint,omitempty"`
	UpdatedAt          *string                `json:"updatedAt,omitempty"`
}

func (o *RelayNamespaceProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *RelayNamespaceProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *RelayNamespaceProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *RelayNamespaceProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
