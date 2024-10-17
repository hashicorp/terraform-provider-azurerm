package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OptionsResource struct {
	AutoScaleSettings *AutoScaleSettings `json:"autoscaleSettings,omitempty"`
	Throughput        *int64             `json:"throughput,omitempty"`
}
