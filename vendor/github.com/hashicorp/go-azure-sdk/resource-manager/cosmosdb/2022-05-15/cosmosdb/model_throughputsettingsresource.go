package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThroughputSettingsResource struct {
	AutoScaleSettings   *AutoscaleSettingsResource `json:"autoscaleSettings,omitempty"`
	MinimumThroughput   *string                    `json:"minimumThroughput,omitempty"`
	OfferReplacePending *string                    `json:"offerReplacePending,omitempty"`
	Throughput          *int64                     `json:"throughput,omitempty"`
}
