package datadogmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorUpdateProperties struct {
	Cspm               *bool             `json:"cspm,omitempty"`
	MonitoringStatus   *MonitoringStatus `json:"monitoringStatus,omitempty"`
	ResourceCollection *bool             `json:"resourceCollection,omitempty"`
}
