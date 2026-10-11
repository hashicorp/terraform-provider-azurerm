package azuremonitorworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Metrics struct {
	InternalId              *string `json:"internalId,omitempty"`
	PrometheusQueryEndpoint *string `json:"prometheusQueryEndpoint,omitempty"`
}
