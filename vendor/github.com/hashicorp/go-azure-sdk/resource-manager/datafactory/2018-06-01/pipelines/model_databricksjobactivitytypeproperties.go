package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksJobActivityTypeProperties struct {
	JobId         interface{}             `json:"jobId"`
	JobParameters *map[string]interface{} `json:"jobParameters,omitempty"`
}
