package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabricksNotebookActivityTypeProperties struct {
	BaseParameters *map[string]interface{}   `json:"baseParameters,omitempty"`
	Libraries      *[]map[string]interface{} `json:"libraries,omitempty"`
	NotebookPath   interface{}               `json:"notebookPath"`
}
