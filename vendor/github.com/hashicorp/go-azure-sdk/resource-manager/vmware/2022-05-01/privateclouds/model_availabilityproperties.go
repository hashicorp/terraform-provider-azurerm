package privateclouds

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailabilityProperties struct {
	SecondaryZone *int64                `json:"secondaryZone,omitempty"`
	Strategy      *AvailabilityStrategy `json:"strategy,omitempty"`
	Zone          *int64                `json:"zone,omitempty"`
}
