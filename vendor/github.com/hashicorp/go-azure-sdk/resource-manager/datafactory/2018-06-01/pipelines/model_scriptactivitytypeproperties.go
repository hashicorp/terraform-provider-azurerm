package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActivityTypeProperties struct {
	LogSettings                 *ScriptActivityTypePropertiesLogSettings `json:"logSettings,omitempty"`
	ReturnMultistatementResult  *bool                                    `json:"returnMultistatementResult,omitempty"`
	ScriptBlockExecutionTimeout *string                                  `json:"scriptBlockExecutionTimeout,omitempty"`
	Scripts                     *[]ScriptActivityScriptBlock             `json:"scripts,omitempty"`
}
