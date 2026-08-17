package dataconnectors

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TiTaxiiDataConnectorProperties struct {
	CollectionId        *string                       `json:"collectionId,omitempty"`
	DataTypes           TiTaxiiDataConnectorDataTypes `json:"dataTypes"`
	FriendlyName        *string                       `json:"friendlyName,omitempty"`
	Password            *string                       `json:"password,omitempty"`
	PollingFrequency    PollingFrequency              `json:"pollingFrequency"`
	TaxiiLookbackPeriod *string                       `json:"taxiiLookbackPeriod,omitempty"`
	TaxiiServer         *string                       `json:"taxiiServer,omitempty"`
	TenantId            string                        `json:"tenantId"`
	UserName            *string                       `json:"userName,omitempty"`
	WorkspaceId         *string                       `json:"workspaceId,omitempty"`
}

func (o *TiTaxiiDataConnectorProperties) GetTaxiiLookbackPeriodAsTime() (*time.Time, error) {
	if o.TaxiiLookbackPeriod == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TaxiiLookbackPeriod, "2006-01-02T15:04:05Z07:00")
}

func (o *TiTaxiiDataConnectorProperties) SetTaxiiLookbackPeriodAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TaxiiLookbackPeriod = &formatted
}
