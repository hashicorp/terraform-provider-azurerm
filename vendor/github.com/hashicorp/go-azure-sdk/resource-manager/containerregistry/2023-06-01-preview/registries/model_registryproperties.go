package registries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryProperties struct {
	AdminUserEnabled           *bool                        `json:"adminUserEnabled,omitempty"`
	AnonymousPullEnabled       *bool                        `json:"anonymousPullEnabled,omitempty"`
	CreationDate               *string                      `json:"creationDate,omitempty"`
	DataEndpointEnabled        *bool                        `json:"dataEndpointEnabled,omitempty"`
	DataEndpointHostNames      *[]string                    `json:"dataEndpointHostNames,omitempty"`
	Encryption                 *EncryptionProperty          `json:"encryption,omitempty"`
	LoginServer                *string                      `json:"loginServer,omitempty"`
	NetworkRuleBypassOptions   *NetworkRuleBypassOptions    `json:"networkRuleBypassOptions,omitempty"`
	NetworkRuleSet             *NetworkRuleSet              `json:"networkRuleSet,omitempty"`
	Policies                   *Policies                    `json:"policies,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState          *ProvisioningState           `json:"provisioningState,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccess         `json:"publicNetworkAccess,omitempty"`
	Status                     *Status                      `json:"status,omitempty"`
	ZoneRedundancy             *ZoneRedundancy              `json:"zoneRedundancy,omitempty"`
}

func (o *RegistryProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *RegistryProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
