package zones

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ZoneProperties struct {
	MaxNumberOfRecordSets          *int64         `json:"maxNumberOfRecordSets,omitempty"`
	MaxNumberOfRecordsPerRecordSet *int64         `json:"maxNumberOfRecordsPerRecordSet,omitempty"`
	NameServers                    *[]string      `json:"nameServers,omitempty"`
	NumberOfRecordSets             *int64         `json:"numberOfRecordSets,omitempty"`
	RegistrationVirtualNetworks    *[]SubResource `json:"registrationVirtualNetworks,omitempty"`
	ResolutionVirtualNetworks      *[]SubResource `json:"resolutionVirtualNetworks,omitempty"`
	ZoneType                       *ZoneType      `json:"zoneType,omitempty"`
}
