package tenantconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeployConfigurationParameterProperties struct {
	Branch string `json:"branch"`
	Force  *bool  `json:"force,omitempty"`
}
