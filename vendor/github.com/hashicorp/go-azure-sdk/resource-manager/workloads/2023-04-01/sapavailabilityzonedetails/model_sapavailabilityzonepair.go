package sapavailabilityzonedetails

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPAvailabilityZonePair struct {
	ZoneA *int64 `json:"zoneA,omitempty"`
	ZoneB *int64 `json:"zoneB,omitempty"`
}
