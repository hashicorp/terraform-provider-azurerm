package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksSparkJarActivityTypeProperties struct {
	Libraries     *[]map[string]interface{} `json:"libraries,omitempty"`
	MainClassName interface{}               `json:"mainClassName"`
	Parameters    *[]interface{}            `json:"parameters,omitempty"`
}
