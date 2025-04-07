package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksSparkPythonActivityTypeProperties struct {
	Libraries  *[]map[string]interface{} `json:"libraries,omitempty"`
	Parameters *[]interface{}            `json:"parameters,omitempty"`
	PythonFile interface{}               `json:"pythonFile"`
}
