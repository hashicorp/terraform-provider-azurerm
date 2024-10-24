package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HopLinkProperties struct {
	RoundTripTimeAvg *int64 `json:"roundTripTimeAvg,omitempty"`
	RoundTripTimeMax *int64 `json:"roundTripTimeMax,omitempty"`
	RoundTripTimeMin *int64 `json:"roundTripTimeMin,omitempty"`
}
