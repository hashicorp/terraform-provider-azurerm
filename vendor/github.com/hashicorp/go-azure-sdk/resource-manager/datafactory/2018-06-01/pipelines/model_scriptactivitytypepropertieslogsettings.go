package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActivityTypePropertiesLogSettings struct {
	LogDestination      ScriptActivityLogDestination `json:"logDestination"`
	LogLocationSettings *LogLocationSettings         `json:"logLocationSettings,omitempty"`
}
