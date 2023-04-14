package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Criteria struct {
	Dimensions *[]Dimension `json:"dimensions,omitempty"`
	MetricName string       `json:"metricName"`
}
