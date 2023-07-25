package contact

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableContacts struct {
	GroundStationName *string                    `json:"groundStationName,omitempty"`
	Properties        *ContactInstanceProperties `json:"properties,omitempty"`
	Spacecraft        *ResourceReference         `json:"spacecraft,omitempty"`
}
