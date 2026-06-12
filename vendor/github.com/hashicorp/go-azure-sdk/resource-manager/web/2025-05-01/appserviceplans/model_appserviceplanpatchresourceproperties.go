package appserviceplans

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServicePlanPatchResourceProperties struct {
	ElasticScaleEnabled       *bool                      `json:"elasticScaleEnabled,omitempty"`
	FreeOfferExpirationTime   *string                    `json:"freeOfferExpirationTime,omitempty"`
	GeoRegion                 *string                    `json:"geoRegion,omitempty"`
	HostingEnvironmentProfile *HostingEnvironmentProfile `json:"hostingEnvironmentProfile,omitempty"`
	HyperV                    *bool                      `json:"hyperV,omitempty"`
	IsSpot                    *bool                      `json:"isSpot,omitempty"`
	IsXenon                   *bool                      `json:"isXenon,omitempty"`
	KubeEnvironmentProfile    *KubeEnvironmentProfile    `json:"kubeEnvironmentProfile,omitempty"`
	MaximumElasticWorkerCount *int64                     `json:"maximumElasticWorkerCount,omitempty"`
	MaximumNumberOfWorkers    *int64                     `json:"maximumNumberOfWorkers,omitempty"`
	NumberOfSites             *int64                     `json:"numberOfSites,omitempty"`
	NumberOfWorkers           *int64                     `json:"numberOfWorkers,omitempty"`
	PerSiteScaling            *bool                      `json:"perSiteScaling,omitempty"`
	ProvisioningState         *ProvisioningState         `json:"provisioningState,omitempty"`
	Reserved                  *bool                      `json:"reserved,omitempty"`
	ResourceGroup             *string                    `json:"resourceGroup,omitempty"`
	SpotExpirationTime        *string                    `json:"spotExpirationTime,omitempty"`
	Status                    *StatusOptions             `json:"status,omitempty"`
	Subscription              *string                    `json:"subscription,omitempty"`
	TargetWorkerCount         *int64                     `json:"targetWorkerCount,omitempty"`
	TargetWorkerSizeId        *int64                     `json:"targetWorkerSizeId,omitempty"`
	WorkerTierName            *string                    `json:"workerTierName,omitempty"`
	ZoneRedundant             *bool                      `json:"zoneRedundant,omitempty"`
}

func (o *AppServicePlanPatchResourceProperties) GetFreeOfferExpirationTimeAsTime() (*time.Time, error) {
	if o.FreeOfferExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FreeOfferExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AppServicePlanPatchResourceProperties) SetFreeOfferExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FreeOfferExpirationTime = &formatted
}

func (o *AppServicePlanPatchResourceProperties) GetSpotExpirationTimeAsTime() (*time.Time, error) {
	if o.SpotExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SpotExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AppServicePlanPatchResourceProperties) SetSpotExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SpotExpirationTime = &formatted
}
