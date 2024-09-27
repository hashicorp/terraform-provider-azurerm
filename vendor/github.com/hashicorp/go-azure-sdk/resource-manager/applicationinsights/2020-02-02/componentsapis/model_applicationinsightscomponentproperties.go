package componentsapis

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentProperties struct {
	AppId                           *string                      `json:"AppId,omitempty"`
	ApplicationId                   *string                      `json:"ApplicationId,omitempty"`
	ApplicationType                 ApplicationType              `json:"Application_Type"`
	ConnectionString                *string                      `json:"ConnectionString,omitempty"`
	CreationDate                    *string                      `json:"CreationDate,omitempty"`
	DisableIPMasking                *bool                        `json:"DisableIpMasking,omitempty"`
	DisableLocalAuth                *bool                        `json:"DisableLocalAuth,omitempty"`
	FlowType                        *FlowType                    `json:"Flow_Type,omitempty"`
	ForceCustomerStorageForProfiler *bool                        `json:"ForceCustomerStorageForProfiler,omitempty"`
	HockeyAppId                     *string                      `json:"HockeyAppId,omitempty"`
	HockeyAppToken                  *string                      `json:"HockeyAppToken,omitempty"`
	ImmediatePurgeDataOn30Days      *bool                        `json:"ImmediatePurgeDataOn30Days,omitempty"`
	IngestionMode                   *IngestionMode               `json:"IngestionMode,omitempty"`
	InstrumentationKey              *string                      `json:"InstrumentationKey,omitempty"`
	LaMigrationDate                 *string                      `json:"LaMigrationDate,omitempty"`
	Name                            *string                      `json:"Name,omitempty"`
	PrivateLinkScopedResources      *[]PrivateLinkScopedResource `json:"PrivateLinkScopedResources,omitempty"`
	ProvisioningState               *string                      `json:"provisioningState,omitempty"`
	PublicNetworkAccessForIngestion *PublicNetworkAccessType     `json:"publicNetworkAccessForIngestion,omitempty"`
	PublicNetworkAccessForQuery     *PublicNetworkAccessType     `json:"publicNetworkAccessForQuery,omitempty"`
	RequestSource                   *RequestSource               `json:"Request_Source,omitempty"`
	RetentionInDays                 *int64                       `json:"RetentionInDays,omitempty"`
	SamplingPercentage              *float64                     `json:"SamplingPercentage,omitempty"`
	TenantId                        *string                      `json:"TenantId,omitempty"`
	WorkspaceResourceId             *string                      `json:"WorkspaceResourceId,omitempty"`
}

func (o *ApplicationInsightsComponentProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ApplicationInsightsComponentProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}

func (o *ApplicationInsightsComponentProperties) GetLaMigrationDateAsTime() (*time.Time, error) {
	if o.LaMigrationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LaMigrationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ApplicationInsightsComponentProperties) SetLaMigrationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LaMigrationDate = &formatted
}
