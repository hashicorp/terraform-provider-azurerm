package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksSparkJarActivityTypeProperties struct {
	Libraries     *[]map[string]interface{} `json:"libraries,omitempty"`
	MainClassName string                    `json:"mainClassName"`
	Parameters    *[]string                 `json:"parameters,omitempty"`
}
