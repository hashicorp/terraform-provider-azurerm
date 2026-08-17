package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProfileProperties struct {
	ExtendedProperties           *map[string]string        `json:"extendedProperties,omitempty"`
	FrontDoorId                  *string                   `json:"frontDoorId,omitempty"`
	LogScrubbing                 *ProfileLogScrubbing      `json:"logScrubbing,omitempty"`
	OriginResponseTimeoutSeconds *int64                    `json:"originResponseTimeoutSeconds,omitempty"`
	ProvisioningState            *ProfileProvisioningState `json:"provisioningState,omitempty"`
	ResourceState                *ProfileResourceState     `json:"resourceState,omitempty"`
}
