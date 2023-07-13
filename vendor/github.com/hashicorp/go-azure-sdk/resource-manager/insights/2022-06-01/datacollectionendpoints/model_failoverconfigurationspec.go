package datacollectionendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverConfigurationSpec struct {
	ActiveLocation *string         `json:"activeLocation,omitempty"`
	Locations      *[]LocationSpec `json:"locations,omitempty"`
}
