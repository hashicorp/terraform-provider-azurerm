package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThroughputSettingsGetPropertiesResource struct {
	AutoScaleSettings   *AutoscaleSettingsResource `json:"autoscaleSettings,omitempty"`
	Etag                *string                    `json:"_etag,omitempty"`
	MinimumThroughput   *string                    `json:"minimumThroughput,omitempty"`
	OfferReplacePending *string                    `json:"offerReplacePending,omitempty"`
	Rid                 *string                    `json:"_rid,omitempty"`
	Throughput          *int64                     `json:"throughput,omitempty"`
	Ts                  *float64                   `json:"_ts,omitempty"`
}
