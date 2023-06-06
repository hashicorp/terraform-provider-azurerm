package connectedregistries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedRegistryProperties struct {
	Activation        *ActivationProperties     `json:"activation,omitempty"`
	ClientTokenIds    *[]string                 `json:"clientTokenIds,omitempty"`
	ConnectionState   *ConnectionState          `json:"connectionState,omitempty"`
	LastActivityTime  *string                   `json:"lastActivityTime,omitempty"`
	Logging           *LoggingProperties        `json:"logging,omitempty"`
	LoginServer       *LoginServerProperties    `json:"loginServer,omitempty"`
	Mode              ConnectedRegistryMode     `json:"mode"`
	NotificationsList *[]string                 `json:"notificationsList,omitempty"`
	Parent            ParentProperties          `json:"parent"`
	ProvisioningState *ProvisioningState        `json:"provisioningState,omitempty"`
	StatusDetails     *[]StatusDetailProperties `json:"statusDetails,omitempty"`
	Version           *string                   `json:"version,omitempty"`
}

func (o *ConnectedRegistryProperties) GetLastActivityTimeAsTime() (*time.Time, error) {
	if o.LastActivityTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastActivityTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectedRegistryProperties) SetLastActivityTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastActivityTime = &formatted
}
