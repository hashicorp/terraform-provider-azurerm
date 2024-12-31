package assets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Dataset struct {
	DataPoints           *[]DataPoint `json:"dataPoints,omitempty"`
	DatasetConfiguration *string      `json:"datasetConfiguration,omitempty"`
	Name                 string       `json:"name"`
	Topic                *Topic       `json:"topic,omitempty"`
}
