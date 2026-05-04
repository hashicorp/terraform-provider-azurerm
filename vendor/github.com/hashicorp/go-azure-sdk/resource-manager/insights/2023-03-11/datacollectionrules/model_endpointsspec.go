package datacollectionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EndpointsSpec struct {
	LogsIngestion    *string `json:"logsIngestion,omitempty"`
	MetricsIngestion *string `json:"metricsIngestion,omitempty"`
}
