package staticsites

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResponseMessageEnvelopeRemotePrivateEndpointConnection struct {
	Error      *ErrorEntity                       `json:"error,omitempty"`
	Id         *string                            `json:"id,omitempty"`
	Identity   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Location   *string                            `json:"location,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Plan       *ArmPlan                           `json:"plan,omitempty"`
	Properties *RemotePrivateEndpointConnection   `json:"properties,omitempty"`
	Sku        *SkuDescription                    `json:"sku,omitempty"`
	Status     *string                            `json:"status,omitempty"`
	Tags       *map[string]string                 `json:"tags,omitempty"`
	Type       *string                            `json:"type,omitempty"`
	Zones      *zones.Schema                      `json:"zones,omitempty"`
}
