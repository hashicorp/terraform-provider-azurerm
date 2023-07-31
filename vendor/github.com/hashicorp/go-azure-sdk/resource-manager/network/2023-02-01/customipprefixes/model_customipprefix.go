package customipprefixes

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomIPPrefix struct {
	Etag             *string                         `json:"etag,omitempty"`
	ExtendedLocation *edgezones.Model                `json:"extendedLocation,omitempty"`
	Id               *string                         `json:"id,omitempty"`
	Location         *string                         `json:"location,omitempty"`
	Name             *string                         `json:"name,omitempty"`
	Properties       *CustomIPPrefixPropertiesFormat `json:"properties,omitempty"`
	Tags             *map[string]string              `json:"tags,omitempty"`
	Type             *string                         `json:"type,omitempty"`
	Zones            *zones.Schema                   `json:"zones,omitempty"`
}
