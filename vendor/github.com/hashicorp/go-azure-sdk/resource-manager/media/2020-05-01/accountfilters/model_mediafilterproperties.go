package accountfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MediaFilterProperties struct {
	FirstQuality          *FirstQuality           `json:"firstQuality"`
	PresentationTimeRange *PresentationTimeRange  `json:"presentationTimeRange"`
	Tracks                *[]FilterTrackSelection `json:"tracks,omitempty"`
}
