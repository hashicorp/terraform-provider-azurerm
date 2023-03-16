package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PerfCounterDataSource struct {
	CounterSpecifiers          *[]string                            `json:"counterSpecifiers,omitempty"`
	Name                       *string                              `json:"name,omitempty"`
	SamplingFrequencyInSeconds *int64                               `json:"samplingFrequencyInSeconds,omitempty"`
	Streams                    *[]KnownPerfCounterDataSourceStreams `json:"streams,omitempty"`
}
