package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogMonitorResourceUpdateParameters struct {
	Properties *MonitorUpdateProperties `json:"properties,omitempty"`
	Sku        *ResourceSku             `json:"sku,omitempty"`
	Tags       *map[string]string       `json:"tags,omitempty"`
}
