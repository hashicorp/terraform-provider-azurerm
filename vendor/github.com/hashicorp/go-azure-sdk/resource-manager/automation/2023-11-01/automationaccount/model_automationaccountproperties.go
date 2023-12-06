package automationaccount

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationAccountProperties struct {
	AutomationHybridServiceUrl *string                      `json:"automationHybridServiceUrl,omitempty"`
	CreationTime               *string                      `json:"creationTime,omitempty"`
	Description                *string                      `json:"description,omitempty"`
	DisableLocalAuth           *bool                        `json:"disableLocalAuth,omitempty"`
	Encryption                 *EncryptionProperties        `json:"encryption,omitempty"`
	LastModifiedBy             *string                      `json:"lastModifiedBy,omitempty"`
	LastModifiedTime           *string                      `json:"lastModifiedTime,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	PublicNetworkAccess        *bool                        `json:"publicNetworkAccess,omitempty"`
	Sku                        *Sku                         `json:"sku,omitempty"`
	State                      *AutomationAccountState      `json:"state,omitempty"`
}

func (o *AutomationAccountProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AutomationAccountProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *AutomationAccountProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AutomationAccountProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
