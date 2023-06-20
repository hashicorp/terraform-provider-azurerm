package namespaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceProperties struct {
	CreatedAt          *string        `json:"createdAt,omitempty"`
	Critical           *bool          `json:"critical,omitempty"`
	DataCenter         *string        `json:"dataCenter,omitempty"`
	Enabled            *bool          `json:"enabled,omitempty"`
	MetricId           *string        `json:"metricId,omitempty"`
	Name               *string        `json:"name,omitempty"`
	NamespaceType      *NamespaceType `json:"namespaceType,omitempty"`
	ProvisioningState  *string        `json:"provisioningState,omitempty"`
	Region             *string        `json:"region,omitempty"`
	ScaleUnit          *string        `json:"scaleUnit,omitempty"`
	ServiceBusEndpoint *string        `json:"serviceBusEndpoint,omitempty"`
	Status             *string        `json:"status,omitempty"`
	SubscriptionId     *string        `json:"subscriptionId,omitempty"`
	UpdatedAt          *string        `json:"updatedAt,omitempty"`
}

func (o *NamespaceProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *NamespaceProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *NamespaceProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *NamespaceProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
