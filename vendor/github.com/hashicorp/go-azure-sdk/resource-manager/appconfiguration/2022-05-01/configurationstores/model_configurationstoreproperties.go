package configurationstores

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationStoreProperties struct {
	CreateMode                 *CreateMode                           `json:"createMode,omitempty"`
	CreationDate               *string                               `json:"creationDate,omitempty"`
	DisableLocalAuth           *bool                                 `json:"disableLocalAuth,omitempty"`
	EnablePurgeProtection      *bool                                 `json:"enablePurgeProtection,omitempty"`
	Encryption                 *EncryptionProperties                 `json:"encryption,omitempty"`
	Endpoint                   *string                               `json:"endpoint,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnectionReference `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState                    `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess                  `json:"publicNetworkAccess,omitempty"`
	SoftDeleteRetentionInDays  *int64                                `json:"softDeleteRetentionInDays,omitempty"`
}

func (o *ConfigurationStoreProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ConfigurationStoreProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
